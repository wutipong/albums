package api

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"log/slog"
	"net/http"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/wutipong/albums/albumscli/server/types"
)

type PostAssetResposnse struct {
	Asset   types.Asset `json:"asset"`
	Success bool        `json:"success"`
}

func PostAsset(
	ctx context.Context,
	server ServerConfig,
	albumID string,
	containerPath string,
	path string,
	reader io.Reader,
	size int64,
) (result PostAssetResposnse, err error) {
	if ctx.Err() != nil {
		err = fmt.Errorf("context error: %w", ctx.Err())
		return
	}
	if server.DryRun {
		slog.Debug(
			"Dry run: skipping asset upload",
			slog.String("path", path),
		)
		result = PostAssetResposnse{
			Asset: types.Asset{
				ID: uuid.NewString(),
			},
			Success: true,
		}

		return
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		err = fmt.Errorf("unable to read from reader: %w", err)
		return
	}

	checksum := crc32.ChecksumIEEE(data)
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, checksum)
	encoded := base64.StdEncoding.EncodeToString(buf)

	assetFileName := filepath.Join(containerPath, path)

	postAssetRequest, err := Post[PostAssetRequestResponse](
		ctx, server, "/api/asset/upload/request",
		PostAssetRequestRequest{
			AlbumID:  albumID,
			Filename: assetFileName,
			Checksum: encoded,
			Network:  server.Network,
		})
	if err != nil {
		err = fmt.Errorf("request to upload failed: %w", err)
		return
	}

	req, err := http.NewRequest(http.MethodPut, postAssetRequest.URL, bytes.NewBuffer(data))
	if err != nil {
		err = fmt.Errorf("failed to create request for put object: %w", err)
		return
	}
	req.ContentLength = size

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		err = fmt.Errorf("failed to put object: %w", err)
		return
	}
	defer resp.Body.Close()

	postAssetCommit, err := Post[PostAssetCommitResponse](
		ctx, server, "/api/asset/upload/commit",
		PostAssetCommitRequest{
			ID: postAssetRequest.ID,
		})

	if err != nil {
		err = fmt.Errorf("unable to commit asset upload %s: %w", postAssetRequest.ID, err)
	}

	result = PostAssetResposnse{
		Asset:   postAssetCommit.Asset,
		Success: postAssetCommit.Success,
	}

	return
}

type PostAssetRequestRequest struct {
	AlbumID  string `json:"album_id"`
	Filename string `json:"filename"`
	Checksum string `json:"checksum"`
	Network  string `json:"network"`
}

type PostAssetRequestResponse struct {
	ID      string `json:"id"`
	URL     string `json:"url"`
	Success bool   `json:"success"`
}

type PostAssetCommitRequest struct {
	ID string `json:"id"`
}

type PostAssetCommitResponse struct {
	Asset   types.Asset `json:"asset"`
	Success bool        `json:"success"`
}
