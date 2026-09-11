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
	object, err := getObject(ctx, minioClient, asset.Original)

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

	var view *vips.Image
	if original.Width()*original.Height() > MAX_VIEW_PIXEL {
		view, err = original.Copy(nil)
		if err != nil {
			return fmt.Errorf("unable to copy original image: %w", err)
		}
		defer view.Close()
	} else {
		view = original
	}

	err = populateView(ctx, minioClient, asset, view, view != original)
	if err != nil {
		return fmt.Errorf("unable to populate view image: %e", err)
	}

	if view == nil {
		view = original
	} else {
		defer view.Close()
	}

	err = populatePreview(ctx, minioClient, asset, view)
	if err != nil {
		return fmt.Errorf("unable to populate preview image: %e", err)
	}

	err = populateThumbnail(ctx, minioClient, asset, view)
	if err != nil {
		return fmt.Errorf("unable to populate thumbnail: %e", err)
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
	view *vips.Image,
	requireProcessing bool,
) error {
	slog.Info("populating view media for asset", slog.String("id", asset.ID.String()))

	err := ctx.Err()
	if err != nil {
		return fmt.Errorf("context cancelled: %w", err)
	}

	if requireProcessing {
		err = view.ThumbnailImage(1_000_000, &vips.ThumbnailImageOptions{
			Height: VIEW_HEIGHT,
			Size:   vips.SizeDown,
		})
		if err != nil {
			return fmt.Errorf("unable to process view image: %w", err)
		}
	}

	asset.ViewWidth = int32(view.Width())
	asset.ViewHeight = int32(view.Height())
	asset.ImageFrames = int32(view.Pages())

	if view.Pages() > 1 {
		asset.ViewHeight = int32(view.PageHeight())
	}

	if !requireProcessing {
		asset.View = asset.Original
		return nil
	}

	buf, err := view.WebpsaveBuffer(nil)
	if err != nil {
		return fmt.Errorf("unable to save to webp image: %w", err)
	}

	if asset.View == "" || asset.View == asset.Original {
		asset.View = createAssetKey("webp")
	}

	err = putObject(ctx, err, minioClient, asset.View, buf)
	if err != nil {
		return fmt.Errorf("unable to put object to S3: %w", err)
	}

	return nil
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

	asset.ImageFrames = int32(view.Pages())

	if asset.ImageFrames == 1 {
		asset.Preview = asset.View

		return nil
	}

	preview, err := view.Copy(nil)
	if err != nil {
		return err
	}

	defer preview.Close()

	err = preview.ThumbnailImage(1_000_000, &vips.ThumbnailImageOptions{
		Height: THUMBNAIL_HEIGHT,
	})
	if err != nil {
		return fmt.Errorf("unable to create preview image: %w", err)
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

	err = putObject(ctx, err, minioClient, asset.Preview, buf)

	if err != nil {
		return fmt.Errorf("unable to put preview object to S3: %w", err)
	}

	return nil
}

func getObject(ctx context.Context, minioClient *minio.Client, key string) (obj *minio.Object, err error) {
	obj, err = minioClient.GetObject(
		ctx, os.Getenv("S3_BUCKET"),
		key,
		minio.GetObjectOptions{},
	)

	if err != nil {
		err = fmt.Errorf("unable to get object from s3: %w", err)
	}

	return
}

func putObject(ctx context.Context, err error, minioClient *minio.Client, key string, buf []byte) error {
	_, err = minioClient.PutObject(
		ctx, os.Getenv("S3_BUCKET"),
		key,
		bytes.NewReader(buf),
		int64(len(buf)),
		minio.PutObjectOptions{
			ContentType: "image/webp",
		},
	)
	return err
}

func populateThumbnail(
	ctx context.Context,
	minioClient *minio.Client,
	asset *db.Asset,
	view *vips.Image,
) error {
	slog.Info("populating thumbnail media for asset", slog.String("id", asset.ID.String()))

	err := ctx.Err()
	if err != nil {
		return fmt.Errorf("context cancelled: %w", err)
	}

	thumbnail, err := view.Copy(nil)
	if err != nil {
		return fmt.Errorf("unable to create thumbnail copy: %w", err)
	}
	defer thumbnail.Close()
	thumbnail.SetPages(1)

	err = thumbnail.ThumbnailImage(1_000_000, &vips.ThumbnailImageOptions{
		Height: THUMBNAIL_HEIGHT,
	})
	if err != nil {
		return fmt.Errorf("unable to create thumbnail image: %w", err)
	}

	asset.ThumbnailWidth = int32(thumbnail.Width())
	asset.ThumbnailHeight = int32(thumbnail.Height())

	if thumbnail.Pages() == 1 {
		asset.ThumbnailHeight = int32(thumbnail.PageHeight())

		return nil
	}

	params := vips.DefaultWebpsaveBufferOptions()
	params.Q = THUMBNAIL_QUALITY

	buf, err := thumbnail.WebpsaveBuffer(params)
	if err != nil {
		return fmt.Errorf("unable to write preview image: %w", err)
	}

	if asset.Thumbnail == "" || asset.Thumbnail == asset.Original {
		asset.Thumbnail = createAssetKey("webp")
	}

	err = putObject(ctx, err, minioClient, asset.Thumbnail, buf)
	if err != nil {
		return fmt.Errorf("unable to put object to S3: %w", err)
	}

	return nil
}
