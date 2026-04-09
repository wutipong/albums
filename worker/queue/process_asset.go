package queue

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/wutipong/albums/worker/db"
)

const THUMBNAIL_FILE = "thumbnail.webp"
const THUMBNAIL_SIZE = 200
const THUMBNAIL_QUALITY = 60

const PREVIEW_FILE = "preview.webp"

const VIEW_FILE = "view.webp"
const VIEW_SIZE = 2000
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
	".mp4",
	".m4v",
	".webm",
}

func ProcessAsset(ctx context.Context, id string, queries *db.Queries) error {
	slog.Info("processing asset", slog.String("id", id))

	err := ctx.Err()
	if err != nil {
		slog.Info("context.", slog.String("error", err.Error()))
		return fmt.Errorf("context cancelled: %w", err)
	}

	var uuid pgtype.UUID

	slog.Info("parsing task id.")

	err = uuid.Scan(id)
	if err != nil {
		slog.Info("error parsing uuid.", slog.String("error", err.Error()))
		return fmt.Errorf("invalid asset id %s: %w", id, err)
	}

	slog.Info("get current asset data.")

	asset, err := queries.GetAsset(ctx, uuid)
	if err != nil {
		return fmt.Errorf("unable to read asset data: %w", err)
	}

	ext := filepath.Ext(asset.Original)
	if slices.Contains(imageExts, ext) {
		err = processImageAsset(ctx, &asset)
		if err != nil {
			return fmt.Errorf("unable to process image asset: %w", err)
		}
	} else if slices.Contains(videoExts, ext) {
		err = processVideoAsset(ctx, &asset)
		if err != nil {
			return fmt.Errorf("unable to process video asset: %w", err)
		}
	} else {
		asset.DeletedAt = pgtype.Timestamptz{Time: time.Now(), Valid: true}
	}

	_, err = queries.UpdateAsset(ctx, db.UpdateAssetParams{
		ID:            uuid,
		Filename:      asset.Filename,
		Checksum:      asset.Checksum,
		Type:          asset.Type,
		Original:      asset.Original,
		Preview:       asset.Preview,
		Thumbnail:     asset.Thumbnail,
		View:          asset.View,
		ProcessStatus: db.ProcessStatusTProcessed,
	})

	if err != nil {
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

	original, err := vips.NewImageFromFile(originalPath)
	if err != nil {
		return fmt.Errorf("unable to read original image: %w", err)
	}
	originalMeta := original.Metadata()
	err = populateThumbnail(ctx, asset, original, originalMeta)
	if err != nil {
		return fmt.Errorf("unable to populate thumbnail: %e", err)
	}

	err = populatePreview(ctx, asset, original, originalMeta)
	if err != nil {
		return fmt.Errorf("unable to populate preview image: %e", err)
	}

	err = populateView(ctx, asset, original, originalMeta)
	if err != nil {
		return fmt.Errorf("unable to populate view image: %e", err)
	}

	return nil
}

func populateView(
	ctx context.Context,
	asset *db.Asset,
	original *vips.ImageRef,
	originalMeta *vips.ImageMetadata,
) error {
	slog.Info("populating view media for asset", slog.String("id", asset.ID.String()))

	err := ctx.Err()
	if err != nil {
		return fmt.Errorf("context cancelled: %w", err)
	}

	if originalMeta.Width <= VIEW_SIZE &&
		originalMeta.Height <= VIEW_SIZE {
		asset.View = asset.Original

		return nil
	}

	view, err := original.Copy()
	if err != nil {
		return fmt.Errorf("unable to copy original: %w", err)
	}

	err = view.AutoRotate()
	if err != nil {
		return fmt.Errorf("unable to perform auto rotating: %w", err)
	}

	err = view.ThumbnailWithSize(
		VIEW_SIZE, VIEW_SIZE, vips.InterestingNone, vips.SizeDown,
	)
	if err != nil {
		return fmt.Errorf("unable to resize preview image")
	}

	params := vips.NewWebpExportParams()
	params.Quality = THUMBNAIL_QUALITY

	buf, _, err := view.ExportWebp(params)
	if err != nil {
		return fmt.Errorf("unable to write preview image: %w", err)
	}

	viewPath := createCacheAssetPath(asset.ID.String(), VIEW_FILE)
	err = os.WriteFile(viewPath, buf, 0644)
	if err != nil {
		return fmt.Errorf("unable to save file: %w", err)
	}

	asset.View = VIEW_FILE
	return nil
}

func populatePreview(
	ctx context.Context,
	asset *db.Asset,
	original *vips.ImageRef,
	originalMeta *vips.ImageMetadata,
) error {
	slog.Info("populating preview media for asset", slog.String("id", asset.ID.String()))

	err := ctx.Err()
	if err != nil {
		return fmt.Errorf("context cancelled: %w", err)
	}

	if originalMeta.Width <= THUMBNAIL_SIZE &&
		originalMeta.Height <= THUMBNAIL_SIZE {
		asset.Preview = asset.Original

		return nil
	}

	preview, err := original.Copy()
	if err != nil {
		return fmt.Errorf("unable to copy original: %w", err)
	}

	err = preview.AutoRotate()
	if err != nil {
		return fmt.Errorf("unable to perform auto rotating: %w", err)
	}

	err = preview.ThumbnailWithSize(
		THUMBNAIL_SIZE, THUMBNAIL_SIZE, vips.InterestingNone, vips.SizeDown,
	)
	if err != nil {
		return fmt.Errorf("unable to resize preview image")
	}

	params := vips.NewWebpExportParams()
	params.Quality = THUMBNAIL_QUALITY

	buf, _, err := preview.ExportWebp(params)
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
	original *vips.ImageRef,
	originalMeta *vips.ImageMetadata,
) error {
	slog.Info("populating thumbnail media for asset", slog.String("id", asset.ID.String()))

	err := ctx.Err()
	if err != nil {
		return fmt.Errorf("context cancelled: %w", err)
	}

	if originalMeta.Width <= THUMBNAIL_SIZE &&
		originalMeta.Height <= THUMBNAIL_SIZE &&
		originalMeta.Pages > 1 {
		asset.Thumbnail = asset.Original

		return nil
	}

	thumbnail, err := original.Copy()
	if err != nil {
		return fmt.Errorf("unable to copy original: %w", err)
	}

	err = thumbnail.AutoRotate()
	if err != nil {
		return fmt.Errorf("unable to perform auto rotating: %w", err)
	}

	err = thumbnail.SetPages(1)
	if err != nil {
		return fmt.Errorf("unable to set page count: %w", err)
	}

	err = thumbnail.ThumbnailWithSize(
		THUMBNAIL_SIZE, THUMBNAIL_SIZE, vips.InterestingNone, vips.SizeDown,
	)
	if err != nil {
		return fmt.Errorf("unable to resize preview image")
	}

	params := vips.NewWebpExportParams()
	params.Quality = THUMBNAIL_QUALITY

	buf, _, err := thumbnail.ExportWebp(params)
	if err != nil {
		return fmt.Errorf("unable to write preview image: %w", err)
	}

	previewPath := createCacheAssetPath(asset.ID.String(), THUMBNAIL_FILE)
	err = os.WriteFile(previewPath, buf, 0644)
	if err != nil {
		return fmt.Errorf("unable to save file: %w", err)
	}

	asset.Thumbnail = THUMBNAIL_FILE
	return nil
}

func createCacheAssetPath(id string, args ...string) string {
	topLevelDir := id[0:2]
	secondLevelDir := id[2:4]

	combined := []string{os.Getenv("CACHE_DIR"), topLevelDir, secondLevelDir, id}
	combined = append(combined, args...)

	return filepath.Join(combined...)
}
