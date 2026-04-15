package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/jackc/pgx/v5/pgtype"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"github.com/wutipong/albums/worker/db"
)

const VIDEO_VIEW_FILE = "view.mp4"
const VIDEO_WIDTH = 1280
const VIDEO_HEIGHT = 720

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

	probe, err := ffmpeg.Probe(originalPath)
	if err != nil {
		return fmt.Errorf("unable to probe original video: %w", err)
	}

	var info Probe
	json.Unmarshal([]byte(probe), &info)

	err = processVideoThumbnail(ctx, asset, info)
	if err != nil {
		return fmt.Errorf("unable to process video asset thumbnail: %w", err)
	}

	err = processVideoPreview(ctx, asset, info)
	if err != nil {
		return fmt.Errorf("unable to process video asset preview: %w", err)
	}

	err = processVideoView(ctx, asset, info)
	if err != nil {
		return fmt.Errorf("unable to process video asset view: %w", err)
	}

	asset.Type = db.AssetTypeTVideo

	return nil
}

func processVideoView(ctx context.Context, asset *db.Asset, info Probe) error {
	slog.Info("process video asset view media", slog.Any("id", asset.ID))
	err := ctx.Err()
	if err != nil {
		return fmt.Errorf("context cancelled: %w", err)
	}

	videoStream, err := info.Video()
	if err != nil {
		return fmt.Errorf("video stream not found: %w", err)
	}

	width := videoStream.Width
	height := videoStream.Height

	if width <= VIDEO_WIDTH &&
		height <= VIDEO_HEIGHT &&
		isVideoBrowserSafe(info) {

		asset.View = asset.Original
		asset.ViewWidth = int32(width)
		asset.ViewHeight = int32(height)

		return nil
	}

	asset.View = VIDEO_VIEW_FILE

	originalPath := createCacheAssetPath(asset.ID.String(), asset.Original)
	viewPath := createCacheAssetPath(asset.ID.String(), asset.View)

	err = ffmpeg.Input(originalPath).
		Output(viewPath, ffmpeg.KwArgs{
			"vf": fmt.Sprintf(
				"scale=%d:%d:force_original_aspect_ratio=decrease,scale=trunc(iw/2)*2:trunc(ih/2)*2",
				VIDEO_WIDTH, VIDEO_HEIGHT,
			),
			"c:v":      "libx264",
			"preset":   "superfast",
			"crf":      "30",
			"pix_fmt":  "yuv420p",    // Fixes browser incompatibility [6]
			"c:a":      "aac",        // Standard audio
			"movflags": "+faststart", // Enables progressive loading [5]
		}).OverWriteOutput().ErrorToStdOut().Run()

	if err != nil {
		return fmt.Errorf("unable to create view asset for video asset: %w", err)
	}

	probe, err := ffmpeg.Probe(viewPath)
	if err != nil {
		return fmt.Errorf("unable to probe original video: %w", err)
	}

	var viewInfo Probe
	json.Unmarshal([]byte(probe), &viewInfo)

	viewVideoStream, err := viewInfo.Video()
	if err != nil {
		return fmt.Errorf("unable to get video stream from video asset: %w", err)
	}
	asset.ViewWidth = int32(viewVideoStream.Width)
	asset.ViewHeight = int32(viewVideoStream.Height)

	return nil
}

func processVideoThumbnail(ctx context.Context, asset *db.Asset, info Probe) error {
	slog.Info("process video asset thumbnail", slog.Any("id", asset.ID))
	err := ctx.Err()
	if err != nil {
		return fmt.Errorf("context cancelled: %w", err)
	}

	asset.Thumbnail = THUMBNAIL_FILE

	originalPath := createCacheAssetPath(asset.ID.String(), asset.Original)
	thumbnailPath := createCacheAssetPath(asset.ID.String(), asset.Thumbnail)

	duration, err := strconv.ParseFloat(info.Format.Duration, 10)
	if err != nil {
		return fmt.Errorf("unable to parse duration: %w", err)
	}

	// save thumbnail at 1/3 duration
	err = ffmpeg.
		Input(originalPath, ffmpeg.KwArgs{
			"ss": fmt.Sprintf("%f", duration/3),
		}).
		Output(thumbnailPath, ffmpeg.KwArgs{
			"c:v":     "libwebp",
			"vframes": "1",
			"quality": fmt.Sprintf("%d", THUMBNAIL_QUALITY),
			"vf": fmt.Sprintf(
				"scale=%d:%d:force_original_aspect_ratio=decrease",
				THUMBNAIL_HEIGHT, THUMBNAIL_HEIGHT,
			),
		}).OverWriteOutput().ErrorToStdOut().Run()

	if err != nil {
		return fmt.Errorf("unable to create thumbnail asset for video asset: %w", err)
	}

	thumbnail, err := vips.NewImageFromFile(thumbnailPath)
	if err != nil {
		return fmt.Errorf("unable to read thumbnail image: %w", err)
	}

	asset.ThumbnailWidth = int32(thumbnail.Width())
	asset.ThumbnailHeight = int32(thumbnail.Height())

	videoDuration := time.Duration(duration) * time.Second
	asset.VideoDuration = pgtype.Interval{
		Microseconds: videoDuration.Microseconds(),
		Valid:        true,
	}
	return nil
}

func processVideoPreview(ctx context.Context, asset *db.Asset, info Probe) error {
	slog.Info("process video preview", slog.Any("id", asset.ID))
	err := ctx.Err()
	if err != nil {
		return fmt.Errorf("context cancelled: %w", err)
	}

	asset.Preview = PREVIEW_FILE

	originalPath := createCacheAssetPath(asset.ID.String(), asset.Original)
	previewPath := createCacheAssetPath(asset.ID.String(), asset.Preview)

	duration, err := strconv.ParseFloat(info.Format.Duration, 10)
	if err != nil {
		return fmt.Errorf("unable to parse duration: %w", err)
	}

	// save preview at 1/3 duration, 5 seconds-long in 5 fps.
	err = ffmpeg.
		Input(originalPath, ffmpeg.KwArgs{
			"ss": fmt.Sprintf("%f", duration/3),
		}).
		Output(previewPath, ffmpeg.KwArgs{
			"c:v":     "libwebp",
			"t":       "5",
			"loop":    "0",
			"quality": fmt.Sprintf("%d", THUMBNAIL_QUALITY),
			"vf": fmt.Sprintf(
				"fps=5,scale=%d:%d:force_original_aspect_ratio=decrease",
				THUMBNAIL_HEIGHT, THUMBNAIL_HEIGHT,
			),
		}).OverWriteOutput().ErrorToStdOut().Run()

	if err != nil {
		return fmt.Errorf("unable to create thumbnail asset for video asset: %w", err)
	}

	return nil
}
