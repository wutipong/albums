package queue

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/minio/minio-go/v7"
	"github.com/wutipong/albums/worker/db"
	vips "github.com/wutipong/albums/worker/vips"
)

const THUMBNAIL_HEIGHT = 200
const MAX_VIEW_PIXEL = 50_000_000
const VIEW_HEIGHT = 2000

func processImageAsset(ctx context.Context, minioClient *minio.Client, asset *db.Asset) error {
	slog.Info("processing image asset", slog.String("id", asset.ID.String()))

	err := ctx.Err()
	if err != nil {
		slog.Info("context.", slog.String("error", err.Error()))
		return fmt.Errorf("context cancelled: %w", err)
	}

	slog.Info("getting object from S3.", slog.String("id", asset.Original))
	object, err := getObjectFromS3(ctx, minioClient, asset.Original)

	if err != nil {
		return fmt.Errorf("unable to get object from s3: %w", err)
	}
	defer object.Close()

	source := vips.NewSource(object)
	defer source.Close()

	slog.Info("read original image file.")

	params := vips.DefaultLoadOptions()
	if hasAnimationExt(filepath.Ext(asset.Filename)) {
		params.N = -1
	}

	original, err := vips.NewImageFromSource(source, params)
	if err != nil {
		return fmt.Errorf("unable to read original image: %w", err)
	}
	defer original.Close()

	slog.Info("original", slog.Any("image", original), slog.Bool("available", original != nil), slog.Int("width", original.Width()))

	view, err := populateView(ctx, minioClient, asset, original)
	if err != nil {
		return fmt.Errorf("unable to populate view image: %e", err)
	}

	if view == nil {
		view = original
	} else {
		defer view.Close()
	}

	err = populateThumbnail(ctx, minioClient, asset, view)
	if err != nil {
		return fmt.Errorf("unable to populate thumbnail: %e", err)
	}

	err = populatePreview(ctx, minioClient, asset, view)
	if err != nil {
		return fmt.Errorf("unable to populate preview image: %e", err)
	}

	embedding, err := GetImageEmbedding(ctx, original)
	if err == nil {
		asset.ImageEmbedding = &embedding
	} else {
		slog.Warn("Unable to populate embedding. Skip.")
	}
	return nil
}

func populateView(
	ctx context.Context,
	minioClient *minio.Client,
	asset *db.Asset,
	original *vips.Image,
) (view *vips.Image, err error) {
	slog.Info("populating view media for asset", slog.String("id", asset.ID.String()))
	slog.Debug("original image",
		slog.Int("width", original.Width()),
		slog.Int("height", original.Height()),
		slog.Int("page_height", original.PageHeight()),
		slog.Int("loop", original.Loop()),
		slog.Int("pages", original.Pages()),
	)

	err = ctx.Err()
	if err != nil {
		err = fmt.Errorf("context cancelled: %w", err)
		return
	}

	if original == nil {
		err = fmt.Errorf("Invalid image")
		return
	}

	asset.ImageFrames = int32(original.Pages())

	if original.Pages() > 1 {
		asset.View = asset.Original
		asset.ViewWidth = int32(original.Width())
		asset.ViewHeight = int32(original.PageHeight())

		return
	}

	if (original.Width() * original.Height()) < MAX_VIEW_PIXEL {
		asset.View = asset.Original
		asset.ViewWidth = int32(original.Width())
		asset.ViewHeight = int32(original.Height())

		return
	}

	view, err = original.Copy(nil)
	if err != nil {
		err = fmt.Errorf("unable to copy original image: %w", err)
		return
	}

	factor := float64(view.Height()) / float64(VIEW_HEIGHT)

	err = view.Resize(factor, &vips.ResizeOptions{
		Kernel: vips.KernelLanczos3,
		Gap:    2,
	})
	if err != nil {
		err = fmt.Errorf("unable to resize view image: %w", err)
		return
	}

	asset.ViewWidth = int32(view.Width())
	asset.ViewHeight = int32(view.Height())

	buf, err := view.WebpsaveBuffer(nil)
	if err != nil {
		err = fmt.Errorf("unable to save to webp image: %w", err)
		return
	}

	if asset.View == "" || asset.View == asset.Original {
		asset.View = createAssetKey("webp")
	}

	_, err = minioClient.PutObject(
		ctx, os.Getenv("S3_BUCKET"),
		asset.View,
		bytes.NewReader(buf),
		int64(len(buf)),
		minio.PutObjectOptions{
			ContentType: "image/webp",
		},
	)

	if err != nil {
		err = fmt.Errorf("unable to put object to S3: %w", err)
		return
	}

	return
}

func populatePreview(
	ctx context.Context,
	minioClient *minio.Client,
	asset *db.Asset,
	view *vips.Image,
) error {
	slog.Info(
		"populating preview media for asset",
		slog.String("id", asset.ID.String()),
	)

	err := ctx.Err()
	if err != nil {
		return fmt.Errorf("context cancelled: %w", err)
	}

	if asset.ImageFrames == 1 {
		asset.Preview = asset.Thumbnail

		return nil
	}

	preview, err := view.Copy(nil)
	if err != nil {
		return fmt.Errorf("unable to create a preview copy from original image: %w", err)
	}
	defer preview.Close()

	err = processAnimatedPreview(preview)
	if err != nil {
		return err
	}

	params := vips.DefaultWebpsaveBufferOptions()
	params.Q = THUMBNAIL_QUALITY
	params.PageHeight = preview.PageHeight()

	buf, err := preview.WebpsaveBuffer(params)
	if err != nil {
		return fmt.Errorf("unable to write preview image: %w", err)
	}

	if asset.Preview == "" || asset.Preview == asset.View {
		asset.Preview = createAssetKey("webp")
	}

	_, err = minioClient.PutObject(
		ctx, os.Getenv("S3_BUCKET"),
		asset.Preview,
		bytes.NewReader(buf),
		int64(len(buf)),
		minio.PutObjectOptions{
			ContentType: "image/webp",
		},
	)

	if err != nil {
		return fmt.Errorf("unable to put preview object to S3: %w", err)
	}

	return nil
}

func processAnimatedPreview(preview *vips.Image) error {
	slog.Debug("original image",
		slog.Int("width", preview.Width()),
		slog.Int("height", preview.Height()),
		slog.Int("page_height", preview.PageHeight()),
		slog.Int("loop", preview.Loop()),
		slog.Int("pages", preview.Pages()),
	)

	factor := float64(THUMBNAIL_HEIGHT) / float64(preview.PageHeight())

	preview.Resize(factor, &vips.ResizeOptions{
		Kernel: vips.KernelLanczos3,
		Gap:    2,
	})

	preview.SetPageHeight(THUMBNAIL_HEIGHT)

	slog.Debug("preview image",
		slog.Int("width", preview.Width()),
		slog.Int("height", preview.Height()),
		slog.Int("page_height", preview.PageHeight()),
		slog.Int("loop", preview.Loop()),
		slog.Int("pages", preview.Pages()),
	)
	return nil
}

func populateThumbnail(
	ctx context.Context,
	minioClient *minio.Client,
	asset *db.Asset,
	view *vips.Image,
) error {
	slog.Info(
		"populating thumbnail media for asset",
		slog.String("id", asset.ID.String()),
	)

	err := ctx.Err()
	if err != nil {
		return fmt.Errorf("context cancelled: %w", err)
	}

	asset.ThumbnailWidth = int32((view.Width() * THUMBNAIL_HEIGHT) / view.Height())
	asset.ThumbnailHeight = THUMBNAIL_HEIGHT

	copyOptions := vips.DefaultCopyOptions()
	thumbnail, err := view.Copy(copyOptions)
	if err != nil {
		return fmt.Errorf("unable to copy from original image: %w", err)
	}

	defer thumbnail.Close()

	err = processThumbnail(thumbnail)

	if err != nil {
		return fmt.Errorf("Unable to process thumbnail image: %w", err)
	}

	asset.ThumbnailWidth = int32((view.Width() * THUMBNAIL_HEIGHT) / view.PageHeight())

	params := vips.DefaultWebpsaveBufferOptions()
	params.Q = THUMBNAIL_QUALITY

	buf, err := thumbnail.WebpsaveBuffer(params)
	if err != nil {
		return fmt.Errorf("unable to write preview image: %w", err)
	}

	if asset.Thumbnail == "" || asset.Thumbnail == asset.Original {
		asset.Thumbnail = createAssetKey("webp")
	}

	_, err = minioClient.PutObject(
		ctx, os.Getenv("S3_BUCKET"),
		asset.Thumbnail,
		bytes.NewReader(buf),
		int64(len(buf)),
		minio.PutObjectOptions{
			ContentType: "image/webp",
		},
	)

	if err != nil {
		return fmt.Errorf("unable to put object to S3: %w", err)
	}

	return nil
}

func processThumbnail(thumbnail *vips.Image) error {
	slog.Debug("original image",
		slog.Int("width", thumbnail.Width()),
		slog.Int("height", thumbnail.Height()),
		slog.Int("page_height", thumbnail.PageHeight()),
		slog.Int("loop", thumbnail.Loop()),
		slog.Int("pages", thumbnail.Pages()),
	)

	err := thumbnail.Autorot(nil)
	if err != nil {
		return fmt.Errorf("unable to perform auto rotating: %w", err)
	}

	factor := float64(THUMBNAIL_HEIGHT) / float64(thumbnail.PageHeight())
	thumbnail.Resize(factor, &vips.ResizeOptions{
		Kernel: vips.KernelLanczos3,
		Gap:    2,
	})

	err = thumbnail.ExtractArea(0, 0, thumbnail.Width(), THUMBNAIL_HEIGHT)
	if err != nil {
		return fmt.Errorf("unable to extract area: %w", err)
	}

	if thumbnail.Pages() > 1 {
		thumbnail.SetPages(1)
		thumbnail.SetPageHeight(THUMBNAIL_HEIGHT)
	}
	slog.Debug("thumbnail image",
		slog.Int("width", thumbnail.Width()),
		slog.Int("height", thumbnail.Height()),
		slog.Int("page_height", thumbnail.PageHeight()),
		slog.Int("loop", thumbnail.Loop()),
		slog.Int("pages", thumbnail.Pages()),
	)
	return nil
}
