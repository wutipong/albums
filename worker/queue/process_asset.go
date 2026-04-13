package queue

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	vips "github.com/cshum/vipsgen/vips816"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/wutipong/albums/worker/db"
)

const THUMBNAIL_FILE = "thumbnail.webp"
const THUMBNAIL_WIDTH = 100_000_000 // don't consider the width
const THUMBNAIL_HEIGHT = 200
const THUMBNAIL_QUALITY = 60

const PREVIEW_FILE = "preview.webp"

const VIEW_FILE = "view.webp"
const VIEW_HEIGHT = 1000
const VIEW_WIDTH = 100_000_000 // don't consider the width
const VIEW_QUALITY = 80

var imageExts = []string{
	".jpg",
	".jpeg",
	".png",
	".gif",
	".svg",
	".tiff",
	".webp",
}

var videoExts = []string{
	".3gp",
	".avi",
	".m4v",
	".mkv",
	".mov",
	".mp4",
	".webm",
}

var animationExts = []string{
	".gif",
	".webm",
}

func hasVideoExt(ext string) bool {
	return slices.Contains(videoExts, strings.ToLower(ext))
}

func hasImageExt(ext string) bool {
	return slices.Contains(imageExts, strings.ToLower(ext))
}
func hasAnimationExt(ext string) bool {
	return slices.Contains(animationExts, strings.ToLower(ext))
}

func ProcessAsset(ctx context.Context, id string) error {
	slog.Info("processing asset", slog.String("id", id))
	var uuid pgtype.UUID
	err := uuid.Scan(id)
	if err != nil {
		return fmt.Errorf("unable to parse id: %w", err)
	}

	queries, _ := db.Get()

	asset, err := queries.GetAsset(ctx, uuid)
	if err != nil {
		return fmt.Errorf("unable to read asset data: %w", err)
	}

	asset.ProcessStatus = db.ProcessStatusTProcessing

	asset, err = queries.UpdateAssetProcessStatus(ctx,
		db.UpdateAssetProcessStatusParams{
			ID:            uuid,
			ProcessStatus: db.ProcessStatusTProcessing,
		},
	)

	if err != nil {
		return fmt.Errorf("unable to save image metadata: %w", err)
	}

	ext := filepath.Ext(asset.Original)
	if hasImageExt(ext) {
		err = processImageAsset(ctx, &asset)
		if err != nil {
			slog.Info("error processing image asset.", slog.String("error", err.Error()))
			return fmt.Errorf("unable to process image asset: %w", err)
		}
	} else if hasVideoExt(ext) {
		err = processVideoAsset(ctx, &asset)
		if err != nil {
			slog.Info("error proessing video.", slog.String("error", err.Error()))
			return fmt.Errorf("unable to process video asset: %w", err)
		}
	} else {
		asset.DeletedAt = pgtype.Timestamptz{Time: time.Now(), Valid: true}
	}

	_, err = queries.UpdateAsset(ctx, db.UpdateAssetParams{
		ID:              uuid,
		Filename:        asset.Filename,
		Checksum:        asset.Checksum,
		Type:            asset.Type,
		Original:        asset.Original,
		Preview:         asset.Preview,
		Thumbnail:       asset.Thumbnail,
		View:            asset.View,
		ProcessStatus:   db.ProcessStatusTProcessed,
		ThumbnailWidth:  asset.ThumbnailWidth,
		ThumbnailHeight: asset.ThumbnailHeight,
		ViewWidth:       asset.ViewWidth,
		ViewHeight:      asset.ViewHeight,
		ImageFrames:     asset.ImageFrames,
		VideoDuration:   asset.VideoDuration,
	})

	if err != nil {
		slog.Info("error parsing uuid.", slog.String("error", err.Error()))
		return fmt.Errorf("unable to save image metadata: %w", err)
	}

	return nil
}

func processImageAsset(ctx context.Context, asset *db.Asset) error {
	slog.Info("processing image asset", slog.String("id", asset.ID.String()))

	err := ctx.Err()
	if err != nil {
		slog.Info("context.", slog.String("error", err.Error()))
		return fmt.Errorf("context cancelled: %w", err)
	}

	id := asset.ID.String()
	originalPath := createCacheAssetPath(id, asset.Original)
	slog.Info("original asset path", slog.String("path", originalPath))

	slog.Info("read original image file.")

	params := vips.DefaultLoadOptions()
	if hasAnimationExt(filepath.Ext(asset.Original)) {
		params.N = -1
	}

	original, err := vips.NewImageFromFile(originalPath, params)
	if err != nil {
		return fmt.Errorf("unable to read original image: %w", err)
	}
	defer original.Close()

	err = populateThumbnail(ctx, asset, original)
	if err != nil {
		return fmt.Errorf("unable to populate thumbnail: %e", err)
	}

	err = populatePreview(ctx, asset, original)
	if err != nil {
		return fmt.Errorf("unable to populate preview image: %e", err)
	}

	err = populateView(ctx, asset, original)
	if err != nil {
		return fmt.Errorf("unable to populate view image: %e", err)
	}

	return nil
}

func populateView(
	ctx context.Context,
	asset *db.Asset,
	original *vips.Image,
) error {
	slog.Info("populating view media for asset", slog.String("id", asset.ID.String()))

	err := ctx.Err()
	if err != nil {
		return fmt.Errorf("context cancelled: %w", err)
	}

	if original.Height() <= VIEW_HEIGHT {
		asset.View = asset.Original
		asset.ViewWidth = int32(original.Width())
		asset.ViewHeight = int32(original.Height())

		return nil
	}

	view, err := original.Copy(nil)
	if err != nil {
		return fmt.Errorf("unable to copy original: %w", err)
	}
	defer view.Close()

	err = view.Autorot(nil)
	if err != nil {
		return fmt.Errorf("unable to perform auto rotating: %w", err)
	}

	options := vips.DefaultThumbnailImageOptions()
	options.Height = VIEW_HEIGHT
	options.Crop = vips.InterestingNone
	options.Size = vips.SizeDown

	err = view.ThumbnailImage(VIEW_WIDTH, options)
	if err != nil {
		return fmt.Errorf("unable to resize preview image: %w", err)
	}

	params := vips.DefaultWebpsaveBufferOptions()
	params.Q = THUMBNAIL_QUALITY
	buf, err := view.WebpsaveBuffer(params)

	if err != nil {
		return fmt.Errorf("unable to write view image: %w", err)
	}

	viewPath := createCacheAssetPath(asset.ID.String(), VIEW_FILE)
	err = os.WriteFile(viewPath, buf, 0644)
	if err != nil {
		return fmt.Errorf("unable to save file: %w", err)
	}

	asset.View = VIEW_FILE
	asset.ViewWidth = int32(view.Width())
	asset.ViewHeight = int32(view.Height())

	return nil
}

func populatePreview(
	ctx context.Context,
	asset *db.Asset,
	original *vips.Image,
) error {
	slog.Info("populating preview media for asset", slog.String("id", asset.ID.String()))

	err := ctx.Err()
	if err != nil {
		return fmt.Errorf("context cancelled: %w", err)
	}

	if original.Height() <= THUMBNAIL_HEIGHT {
		asset.Preview = asset.Original
		asset.ImageFrames = 1

		return nil
	}

	if original.Pages() == 1 {
		asset.Preview = asset.Thumbnail
		asset.ImageFrames = 1

		return nil
	}

	asset.ImageFrames = int32(original.Pages())

	preview, err := original.Copy(nil)
	if err != nil {
		return fmt.Errorf("unable to copy original: %w", err)
	}
	defer preview.Close()

	err = preview.Autorot(nil)
	if err != nil {
		return fmt.Errorf("unable to perform auto rotating: %w", err)
	}

	options := vips.DefaultThumbnailImageOptions()
	options.Height = THUMBNAIL_HEIGHT
	options.Crop = vips.InterestingNone
	options.Size = vips.SizeDown

	err = preview.ThumbnailImage(
		THUMBNAIL_WIDTH, options,
	)

	if err != nil {
		return fmt.Errorf("unable to resize image: %w", err)
	}

	params := vips.DefaultWebpsaveBufferOptions()
	params.Q = THUMBNAIL_QUALITY

	buf, err := preview.WebpsaveBuffer(params)
	if err != nil {
		return fmt.Errorf("unable to write preview image: %w", err)
	}

	previewPath := createCacheAssetPath(asset.ID.String(), PREVIEW_FILE)
	err = os.WriteFile(previewPath, buf, 0644)
	if err != nil {
		return fmt.Errorf("unable to save file: %w", err)
	}

	asset.Preview = PREVIEW_FILE
	return nil
}

func populateThumbnail(
	ctx context.Context,
	asset *db.Asset,
	original *vips.Image,
) error {
	slog.Info("populating thumbnail media for asset", slog.String("id", asset.ID.String()))

	err := ctx.Err()
	if err != nil {
		return fmt.Errorf("context cancelled: %w", err)
	}

	if original.Height() <= THUMBNAIL_HEIGHT &&
		original.Pages() == 1 {
		asset.Thumbnail = asset.Original
		asset.ThumbnailWidth = int32(original.Width())
		asset.ThumbnailHeight = int32(original.Height())

		return nil
	}

	originalPath := createCacheAssetPath(asset.ID.String(), asset.Original)
	thumbnail, _ := vips.NewImageFromFile(originalPath, nil)
	defer thumbnail.Close()

	err = thumbnail.Autorot(nil)
	if err != nil {
		return fmt.Errorf("unable to perform auto rotating: %w", err)
	}

	options := vips.DefaultThumbnailImageOptions()
	options.Height = THUMBNAIL_HEIGHT
	options.Crop = vips.InterestingNone
	options.Size = vips.SizeDown

	err = thumbnail.ThumbnailImage(THUMBNAIL_WIDTH, options)
	if err != nil {
		return fmt.Errorf("unable to resize image: %w", err)
	}

	params := vips.DefaultWebpsaveBufferOptions()
	params.Q = THUMBNAIL_QUALITY

	buf, err := thumbnail.WebpsaveBuffer(params)
	if err != nil {
		return fmt.Errorf("unable to write preview image: %w", err)
	}

	previewPath := createCacheAssetPath(asset.ID.String(), THUMBNAIL_FILE)
	err = os.WriteFile(previewPath, buf, 0644)
	if err != nil {
		return fmt.Errorf("unable to save file: %w", err)
	}

	asset.Thumbnail = THUMBNAIL_FILE
	asset.ThumbnailWidth = int32(thumbnail.Width())
	asset.ThumbnailHeight = int32(thumbnail.Height())

	return nil
}

func createCacheAssetPath(id string, args ...string) string {
	topLevelDir := id[0:2]
	secondLevelDir := id[2:4]

	combined := []string{
		os.Getenv("CACHE_DIR"),
		"assets",
		topLevelDir,
		secondLevelDir,
		id,
	}
	combined = append(combined, args...)

	return filepath.Join(combined...)
}
