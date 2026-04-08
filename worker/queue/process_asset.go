package queue

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

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

func ProcessAsset(ctx context.Context, id string) error {
	var uuid pgtype.UUID

	uuid.Scan(id)
	queries := db.New(db.Connection())
	asset, err := queries.GetAsset(ctx, uuid)
	if err != nil {
		return fmt.Errorf("unable to read asset data: %w", err)
	}

	originalPath := createCacheAssetPath(id, asset.Original)
	original, err := vips.NewImageFromFile(originalPath)
	if err != nil {
		return fmt.Errorf("unable to read original image: %w", err)
	}
	originalMeta := original.Metadata()
	err = populateThumbnail(ctx, &asset, original, originalMeta)
	if err != nil {
		return fmt.Errorf("unable to populate thumbnail: %e", err)
	}

	err = populatePreview(ctx, &asset, original, originalMeta)
	if err != nil {
		return fmt.Errorf("unable to populate preview image: %e", err)
	}

	err = populateView(ctx, &asset, original, originalMeta)
	if err != nil {
		return fmt.Errorf("unable to populate view image: %e", err)
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
		ProcessStatus: asset.ProcessStatus,
	})

	if err != nil {
		return fmt.Errorf("unable to save image metadata: %w", err)
	}

	return nil
}

func populateView(
	ctx context.Context,
	asset *db.Asset,
	original *vips.ImageRef,
	originalMeta *vips.ImageMetadata,
) error {
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
		THUMBNAIL_SIZE, THUMBNAIL_SIZE, vips.InterestingNone, vips.SizeUp,
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
		THUMBNAIL_SIZE, THUMBNAIL_SIZE, vips.InterestingNone, vips.SizeUp,
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
		THUMBNAIL_SIZE, THUMBNAIL_SIZE, vips.InterestingNone, vips.SizeUp,
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

	combined := []string{topLevelDir, secondLevelDir}
	combined = append(combined, args...)

	return filepath.Join(combined...)
}
