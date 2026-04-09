package queue

import (
	"context"
	"fmt"
	"log/slog"

	ffmpeg "github.com/u2takey/ffmpeg-go"
	"github.com/wutipong/albums/worker/db"
)

const VIDEO_VIEW_FILE = "view.webm"

func processVideoAsset(ctx context.Context, asset *db.Asset) error {
	slog.Info("process video asset", slog.Any("id", asset.ID))

	err := ctx.Err()
	if err != nil {
		slog.Info("context.", slog.String("error", err.Error()))
		return fmt.Errorf("context cancelled: %w", err)
	}
	id := asset.ID.String()
	originalPath := createCacheAssetPath(id, asset.Original)
	slog.Info("original asset path", slog.String("path", originalPath))

	err = processVideoPreview(ctx, asset)
	if err != nil {
		return fmt.Errorf("unable to process video asset preview: %w", err)
	}
	err = processVideoThumbnail(ctx, asset)
	if err != nil {
		return fmt.Errorf("unable to process video asset thumbnail: %w", err)
	}
	err = processVideoView(ctx, asset)
	if err != nil {
		return fmt.Errorf("unable to process video asset view: %w", err)
	}

	return nil
}

func processVideoView(ctx context.Context, asset *db.Asset) error {
	slog.Info("process video asset view media", slog.Any("id", asset.ID))
	err := ctx.Err()
	if err != nil {
		return fmt.Errorf("context cancelled: %w", err)
	}

	asset.View = VIDEO_VIEW_FILE

	originalPath := createCacheAssetPath(asset.ID.String(), asset.Original)
	viewPath := createCacheAssetPath(asset.ID.String(), asset.View)

	//transcode video to 720p WEBM format
	err = ffmpeg.Input(originalPath).
		Output(viewPath, ffmpeg.KwArgs{
			"vf":  "scale=1280:720:force_original_aspect_ratio=decrease",
			"c:v": "libvpx-vp9",
			"crf": "30",
			"b:v": "0",
			"c:a": "libopus",
		}).OverWriteOutput().ErrorToStdOut().Run()

	if err != nil {
		return fmt.Errorf("unable to create view asset for video asset: %w", err)
	}

	return nil
}

func processVideoThumbnail(ctx context.Context, asset *db.Asset) error {
	slog.Info("process video asset thumbnail", slog.Any("id", asset.ID))
	err := ctx.Err()
	if err != nil {
		return fmt.Errorf("context cancelled: %w", err)
	}

	asset.Thumbnail = THUMBNAIL_FILE

	originalPath := createCacheAssetPath(asset.ID.String(), asset.Original)
	thumbnailPath := createCacheAssetPath(asset.ID.String(), asset.Thumbnail)

	//transcode video to 720p WEBM format
	err = ffmpeg.Input(originalPath).
		Filter("select", ffmpeg.Args{"eq(n,trunc(30/100*n_frames))"}).
		Output(thumbnailPath, ffmpeg.KwArgs{
			"vf": fmt.Sprintf(
				"scale=%d:%d:force_original_aspect_ratio=decrease",
				THUMBNAIL_SIZE, THUMBNAIL_SIZE,
			),
			"quality": fmt.Sprintf("%d", THUMBNAIL_QUALITY),
		}).OverWriteOutput().ErrorToStdOut().Run()

	if err != nil {
		return fmt.Errorf("unable to create thumbnail asset for video asset: %w", err)
	}

	return nil
}

func processVideoPreview(ctx context.Context, asset *db.Asset) error {
	slog.Info("process video preview", slog.Any("id", asset.ID))
	err := ctx.Err()
	if err != nil {
		return fmt.Errorf("context cancelled: %w", err)
	}

	asset.Thumbnail = PREVIEW_FILE

	originalPath := createCacheAssetPath(asset.ID.String(), asset.Original)
	previewPath := createCacheAssetPath(asset.ID.String(), asset.Thumbnail)

	//transcode video to 720p WEBM format
	err = ffmpeg.Input(originalPath, ffmpeg.KwArgs{"ss": "5"}).
		Filter("select", ffmpeg.Args{"eq(n,trunc(30/100*n_frames))"}).
		Output(previewPath, ffmpeg.KwArgs{
			"vf": fmt.Sprintf(
				"scale=%d:%d:force_original_aspect_ratio=decrease",
				THUMBNAIL_SIZE, THUMBNAIL_SIZE,
			),
			"quality": fmt.Sprintf("%d", THUMBNAIL_QUALITY),
		}).OverWriteOutput().ErrorToStdOut().Run()

	if err != nil {
		return fmt.Errorf("unable to create thumbnail asset for video asset: %w", err)
	}

	return nil
}
