package api

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"log/slog"
	"net/http"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/wutipong/albums/albumscli/server/types"
)

type PostAssetResposnse struct {
	Asset   types.Asset `json:"asset"`
	Success bool        `json:"success"`
}

var (
	ErrDuplicateAsset = errors.New("duplicate asset")
)

func PostAsset(
	ctx context.Context,
	server ServerConfig,
	albumID string,
	containerPath string,
	path string,
	reader io.Reader,
	size int64,
) (result PostAssetResposnse, err error) {
	start := time.Now()
	defer func() {
		slog.Info("PostAsset completed",
			slog.Duration("duration", time.Since(start)),
			slog.String("path", path),
		)
	}()

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
		if err.Error() == "duplicate asset" {
			err = ErrDuplicateAsset
		}
		err = fmt.Errorf("request to upload failed: %w", err)
		return
	}

	success := true
	err = doPutObject(ctx, postAssetRequest.URL, data, size)
	if err != nil {
		slog.Error("put object fails", slog.String("error", err.Error()))
		success = false
	}

	postAssetCommit, err := Post[PostAssetCommitResponse](
		ctx, server, "/api/asset/upload/commit",
		PostAssetCommitRequest{
			ID:      postAssetRequest.ID,
			Success: success,
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

func doPutObject(ctx context.Context, url string, data []byte, size int64) error {
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(data))
	if err != nil {
		err = fmt.Errorf("failed to create request for put object: %w", err)
		return err
	}
	req.ContentLength = size
	for i := range 10 {
		_, err = http.DefaultClient.Do(req)
		if err == nil {
			break
		}

		slog.Warn("Error occurred when uploading to S3, retrying.",
			slog.Int("retry", i),
			slog.String("error", err.Error()),
		)

		select {
		case <-ctx.Done():
			err = fmt.Errorf("context error during put object: %w", ctx.Err())
			return err
		case <-time.After(time.Duration(i+1) * 10 * time.Second):
		}
	}

	if err != nil {
		err = fmt.Errorf("failed to put object: %w", err)
		return err
	}

	return nil
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
	ID      string `json:"id"`
	Success bool   `json:"success"`
}

type PostAssetCommitResponse struct {
	Asset   types.Asset `json:"asset"`
	Success bool        `json:"success"`
}
