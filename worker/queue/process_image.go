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

		err = view.ThumbnailImage(1_000_000, &vips.ThumbnailImageOptions{
			Height: VIEW_HEIGHT,
			Size:   vips.SizeDown,
		})
		if err != nil {
			return fmt.Errorf("unable to process view image: %w", err)
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

	} else {
		view = original
		asset.View = asset.Original
	}

	err = populateView(ctx, asset, view)
	if err != nil {
		return fmt.Errorf("unable to populate view image: %w", err)
	}

	var preview *vips.Image
	if (original.Pages() > 1 && original.PageHeight() > THUMBNAIL_HEIGHT) ||
		(original.Pages() == 1 && original.Height() > THUMBNAIL_HEIGHT) {
		preview, err = original.Copy(nil)
		if err != nil {
			return fmt.Errorf("unable to copy original image: %w", err)
		}
		defer preview.Close()

		err = preview.ThumbnailImage(1_000_000, &vips.ThumbnailImageOptions{
			Height: THUMBNAIL_HEIGHT,
		})
		if err != nil {
			return fmt.Errorf("unable to process preview image: %w", err)
		}

		buf, err := preview.WebpsaveBuffer(nil)
		if err != nil {
			return fmt.Errorf("unable to save to webp image: %w", err)
		}

		if asset.Preview == "" || asset.Preview == asset.Original {
			asset.Preview = createAssetKey("webp")
		}

		err = putObject(ctx, err, minioClient, asset.Preview, buf)
		if err != nil {
			return fmt.Errorf("unable to put object to S3: %w", err)
		}
	} else {
		preview = original
		asset.Preview = asset.Original
	}

	err = populatePreview(ctx, asset, view)
	if err != nil {
		return fmt.Errorf("unable to populate preview image: %w", err)
	}

	var thumbnail *vips.Image
	if original.Pages() > 1 || (original.Pages() == 1 && original.Height() > THUMBNAIL_HEIGHT) {
		thumbnail, err = original.Copy(nil)
		if err != nil {
			return fmt.Errorf("unable to copy original image: %w", err)
		}
		defer thumbnail.Close()

		if thumbnail.Pages() > 1 {
			thumbnail.SetPages(1)
			thumbnail.ExtractArea(0, 0, thumbnail.Width(), thumbnail.PageHeight())
		}

		err = thumbnail.ThumbnailImage(1_000_000, &vips.ThumbnailImageOptions{
			Height: THUMBNAIL_HEIGHT,
		})
		if err != nil {
			return fmt.Errorf("unable to process preview image: %w", err)
		}

		buf, err := thumbnail.WebpsaveBuffer(nil)
		if err != nil {
			return fmt.Errorf("unable to save to webp image: %w", err)
		}

		if asset.Thumbnail == "" || asset.Thumbnail == asset.Original {
			asset.Thumbnail = createAssetKey("webp")
		}

		err = putObject(ctx, err, minioClient, asset.Thumbnail, buf)
		if err != nil {
			return fmt.Errorf("unable to put object to S3: %w", err)
		}
	} else {
		thumbnail = original
		asset.Thumbnail = asset.Original
	}

	err = populateThumbnail(ctx, asset, thumbnail)
	if err != nil {
		return fmt.Errorf("unable to populate thumbnail: %w", err)
	}

	embedding, err := GetImageEmbedding(ctx, original)
	if err == nil {
		asset.ImageEmbedding = &embedding
	} else {
		slog.Warn("unable to populate embedding. Skip.")
	}
	return nil
}

func populateView(
	ctx context.Context,
	asset *db.Asset,
	view *vips.Image,
) error {
	slog.Info("populating view media for asset", slog.String("id", asset.ID.String()))

	err := ctx.Err()
	if err != nil {
		return fmt.Errorf("context cancelled: %w", err)
	}

	asset.ViewWidth = int32(view.Width())
	asset.ViewHeight = int32(view.Height())
	asset.ImageFrames = int32(view.Pages())

	if view.Pages() > 1 {
		asset.ViewHeight = int32(view.PageHeight())
	}

	return nil
}

func populatePreview(
	ctx context.Context,
	asset *db.Asset,
	preview *vips.Image,
) error {
	slog.Info(
		"populating preview media for asset",
		slog.String("id", asset.ID.String()),
	)

	err := ctx.Err()
	if err != nil {
		return fmt.Errorf("context cancelled: %w", err)
	}

	asset.ImageFrames = int32(preview.Pages())

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
	asset *db.Asset,
	thumbnail *vips.Image,
) error {
	slog.Info("populating thumbnail media for asset", slog.String("id", asset.ID.String()))

	err := ctx.Err()
	if err != nil {
		return fmt.Errorf("context cancelled: %w", err)
	}

	asset.ThumbnailWidth = int32(thumbnail.Width())
	asset.ThumbnailHeight = int32(thumbnail.Height())

	if thumbnail.Pages() == 1 {
		asset.ThumbnailHeight = int32(thumbnail.PageHeight())

		return nil
	}

	return nil
}
