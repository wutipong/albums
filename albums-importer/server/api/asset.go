package api

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/wutipong/albums/albums-importer/server/types"
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
	modDate time.Time,
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

	assetFileName := filepath.Join(containerPath, path)

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	_ = writer.WriteField("albumId", albumID)

	part, err := writer.CreateFormFile("file", assetFileName)
	if err != nil {
		err = fmt.Errorf("failed to create form file: %w", err)
		return
	}

	_, err = io.Copy(part, reader)
	if err != nil {
		err = fmt.Errorf("failed to write data to form file: %w", err)
		return
	}
	_ = writer.Close()

	return DoRequestWithReturnObject[PostAssetResposnse](
		ctx, "POST", server, "/api/asset", &body, writer.FormDataContentType(),
	)
}
