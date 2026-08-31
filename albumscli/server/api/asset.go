package api

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/imroc/req/v3"
	"github.com/wutipong/albums/albumscli/server/types"
)

type PostAssetResposnse struct {
	Asset types.Asset `json:"asset"`
}

var (
	ErrDuplicateAsset = errors.New("duplicate asset")
	ErrUploadFailed   = errors.New("asset upload fails")
)

func PostAsset(
	ctx context.Context,
	server ServerConfig,
	albumID string,
	containerPath string,
	path string,
	reader io.Reader,
	size int64,
) (asset types.Asset, err error) {
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
		asset = types.Asset{
			ID: uuid.NewString(),
		}

		return
	}

	assetFileName := filepath.Join(containerPath, path)

	data, err := io.ReadAll(reader)
	if err != nil {
		err = fmt.Errorf("unable to read from reader: %w", err)
		return
	}

	c := NewClient(server)

	postAssetRequest, err := doPostAssetRequest(ctx, albumID, assetFileName, c)
	if err != nil {
		return
	}

	slog.Debug("upload request results",
		slog.String("ID", postAssetRequest.ID),
		slog.String("url", postAssetRequest.URL),
		slog.Bool("success", postAssetRequest.Success),
	)

	slog.Debug("put object", slog.Int("size", len(data)), slog.String("url", postAssetRequest.URL))

	err = doPutObject(ctx, c, postAssetRequest.URL, data)
	success := true

	if err != nil {
		slog.Error("put object fails", slog.String("error", err.Error()))
		success = false
		err = nil
	}

	slog.Debug("upload commit", slog.String("id", postAssetRequest.ID), slog.String("url", postAssetRequest.URL))

	asset, err = doCommitAsset(ctx, c, success, postAssetRequest.ID)

	return
}

func doCommitAsset(ctx context.Context, c *req.Client, success bool, id string) (asset types.Asset, err error) {
	var errorResponse ErrorResponse
	var postAssetCommit PostAssetCommitResponse

	r := c.Post("/api/asset/upload/commit").
		SetSuccessResult(&postAssetCommit).
		SetErrorResult(&errorResponse).
		SetBodyJsonMarshal(PostAssetCommitRequest{
			ID:      id,
			Success: success,
		}).Do(ctx)

	if r.IsErrorState() {
		err = fmt.Errorf("commit asset fails: %w", errorResponse)
		return
	}

	if r.Err != nil {
		err = fmt.Errorf("commit asset fails: %w", r.Err)
		return
	}

	switch postAssetCommit.Status {
	case "asset status is updated, but the upload failed.":
		err = ErrUploadFailed
		return
	case "asset is commited, but it is not queued to processing.":
		slog.Info("asset is commited, but it is not queued to processing.")
	}

	asset = postAssetCommit.Asset

	return

}

func doPutObject(ctx context.Context, c *req.Client, presigned string, data []byte) error {
	r := c.Put(presigned).
		SetBodyBytes(data).
		SetRetryCount(10).
		Do(ctx)

	slog.Debug("put object response", slog.String("status", r.Status))

	if r.Err != nil {
		return fmt.Errorf("unable to put object to s3: %w", r.Err)
	}

	return nil
}

func doPostAssetRequest(ctx context.Context, albumID string, assetFileName string, c *req.Client) (request PostAssetRequestResponse, err error) {
	slog.Debug("upload request", slog.String("album_id", albumID), slog.String("filename", assetFileName))

	var errorResponse ErrorResponse
	r := c.Post("/api/asset/upload/request").
		SetBodyJsonMarshal(PostAssetRequestRequest{
			AlbumID:  albumID,
			Filename: assetFileName,
		}).
		SetSuccessResult(&request).
		SetErrorResult(&errorResponse).
		Do(ctx)

	if r.IsErrorState() {
		slog.Debug("upload request error", slog.String("error", errorResponse.Status))
		if errorResponse.Status == "duplicate asset" {
			err = fmt.Errorf("request to upload asset failed: %w", ErrDuplicateAsset)

			return
		} else {
			err = fmt.Errorf("request to upload asset failed: %w", errorResponse)
		}
	}

	if r.Err != nil {
		err = fmt.Errorf("request to upload asset failed: %w", r.Err)
		return
	}
	return
}

type PostAssetRequestRequest struct {
	AlbumID  string `json:"album_id"`
	Filename string `json:"filename"`
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
	Asset  types.Asset `json:"asset"`
	Status string      `json:"status"`
}
