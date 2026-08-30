package api

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"log/slog"
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

	c := NewClient(server)

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

	var postAssetRequest PostAssetRequestResponse
	var errorResponse ErrorResponse
	r := c.Post("/api/asset/upload/request").
		SetBodyJsonMarshal(PostAssetRequestRequest{
			AlbumID:  albumID,
			Filename: assetFileName,
			Checksum: encoded,
		}).
		SetSuccessResult(&postAssetRequest).
		SetErrorResult(&errorResponse).
		Do(ctx)

	if r.IsErrorState() {
		if errorResponse.Message == "duplicate asset" {
			err = fmt.Errorf("request to upload failed: %w", ErrDuplicateAsset)

			return
		}
	}

	if r.Err != nil {
		err = fmt.Errorf("request to upload failed: %w", r.Err)
		return
	}

	r = c.Put(postAssetRequest.URL).
		SetBodyBytes(data).
		SetRetryCount(10).
		Do(ctx)

	success := true
	if r.Err != nil {
		slog.Error("put object fails", slog.String("error", err.Error()))
		success = false
	}

	var postAssetCommit PostAssetCommitResponse

	r = c.Post("/api/asset/upload/commit").
		SetSuccessResult(&postAssetCommit).
		SetBodyJsonMarshal(PostAssetCommitRequest{
			ID:      postAssetRequest.ID,
			Success: success,
		}).Do(ctx)

	if r.Err != nil {
		err = fmt.Errorf("unable to commit asset upload %s: %w", postAssetRequest.ID, r.Err)
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
	ID      string `json:"id"`
	Success bool   `json:"success"`
}

type PostAssetCommitResponse struct {
	Asset   types.Asset `json:"asset"`
	Success bool        `json:"success"`
}
