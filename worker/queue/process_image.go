package queue

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log/slog"
	"math"
	"os"
	"path/filepath"

	vips "github.com/cshum/vipsgen/vips816"
	"github.com/minio/minio-go/v7"
	"github.com/pgvector/pgvector-go"
	"github.com/wutipong/albums/worker/clip"
	"github.com/wutipong/albums/worker/db"
)

const THUMBNAIL_HEIGHT = 200

func processImageAsset(ctx context.Context, minioClient *minio.Client, asset *db.Asset) error {
	slog.Info("processing image asset", slog.String("id", asset.ID.String()))

	err := ctx.Err()
	if err != nil {
		slog.Info("context.", slog.String("error", err.Error()))
		return fmt.Errorf("context cancelled: %w", err)
	}

	slog.Info("getting object from S3.", slog.String("id", asset.Original))
	object, err := minioClient.GetObject(
		ctx, os.Getenv("S3_BUCKET"),
		asset.Original,
		minio.GetObjectOptions{},
	)

	if err != nil {
		return fmt.Errorf("unable to get object from s3: %w", err)
	}
	defer object.Close()

	slog.Info("reading data", slog.String("id", asset.Original))
	buff, err := io.ReadAll(object)
	if err != nil {
		return fmt.Errorf("unable to read object from s3: %w", err)
	}

	slog.Info("read original image file.")

	params := vips.DefaultLoadOptions()
	if hasAnimationExt(filepath.Ext(asset.Filename)) {
		params.N = -1
	}

	original, err := vips.NewImageFromBuffer(buff, params)
	if err != nil {
		return fmt.Errorf("unable to read original image: %w", err)
	}
	defer original.Close()

	err = populateView(ctx, minioClient, asset, original)
	if err != nil {
		return fmt.Errorf("unable to populate view image: %e", err)
	}

	err = populatePreview(ctx, minioClient, asset, original)
	if err != nil {
		return fmt.Errorf("unable to populate preview image: %e", err)
	}

	err = populateThumbnail(ctx, minioClient, asset, original)
	if err != nil {
		return fmt.Errorf("unable to populate thumbnail: %e", err)
	}

	err = PopulateImageEmbedding(ctx, asset, original)
	if err != nil {
		return fmt.Errorf("unable to populate image embedding: %w", err)
	}
	return nil
}

func populateView(
	ctx context.Context,
	minioClient *minio.Client,
	asset *db.Asset,
	original *vips.Image,
) error {
	slog.Info("populating view media for asset", slog.String("id", asset.ID.String()))

	err := ctx.Err()
	if err != nil {
		return fmt.Errorf("context cancelled: %w", err)
	}

	asset.ViewWidth = int32(original.Width())
	asset.ViewHeight = int32(original.Height())
	if original.Pages() > 1 {
		asset.ViewHeight = int32(original.PageHeight())
	}

	slog.Info("filename", slog.String("name", asset.Filename))

	if filepath.Ext(asset.Filename) != ".gif" {
		asset.View = asset.Original

		return nil
	}

	view, err := original.Copy(nil)
	if err != nil {
		return fmt.Errorf("unable to copy original image: %w", err)
	}

	buf, err := view.WebpsaveBuffer(nil)
	if err != nil {
		return fmt.Errorf("unable to save to webp image.")
	}

	asset.View = createAssetKey()

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
		return fmt.Errorf("unable to put object to S3: %w", err)
	}

	return nil
}

func populatePreview(
	ctx context.Context,
	minioClient *minio.Client,
	asset *db.Asset,
	original *vips.Image,
) error {
	slog.Info(
		"populating preview media for asset",
		slog.String("id", asset.ID.String()),
	)

	err := ctx.Err()
	if err != nil {
		return fmt.Errorf("context cancelled: %w", err)
	}

	asset.ImageFrames = int32(original.Pages())

	if asset.ImageFrames == 1 {
		asset.Preview = asset.Original

		return nil
	}

	preview, err := createPreviewForAnimationImage(original)
	if err != nil {
		return err
	}

	params := vips.DefaultWebpsaveBufferOptions()
	params.Q = THUMBNAIL_QUALITY

	buf, err := preview.WebpsaveBuffer(params)
	if err != nil {
		return fmt.Errorf("unable to write preview image: %w", err)
	}

	asset.Preview = createAssetKey()

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

func createPreviewForAnimationImage(original *vips.Image) (*vips.Image, error) {
	slog.Debug("original image",
		slog.Int("width", original.Width()),
		slog.Int("height", original.Height()),
		slog.Int("page_height", original.PageHeight()),
		slog.Int("loop", original.Loop()),
		slog.Int("pages", original.Pages()),
	)

	preview, err := original.Copy(nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create a preview copy from original image: %w", err)
	}

	factor := float64(THUMBNAIL_HEIGHT) / float64(original.PageHeight())

	preview.Resize(factor, &vips.ResizeOptions{
		Kernel: vips.KernelLanczos3,
		Gap:    2,
	})

	preview.SetPageHeight(preview.Height() / preview.Pages())

	slog.Debug("preview image",
		slog.Int("width", preview.Width()),
		slog.Int("height", preview.Height()),
		slog.Int("page_height", preview.PageHeight()),
		slog.Int("loop", preview.Loop()),
		slog.Int("pages", preview.Pages()),
	)
	return preview, nil
}

func populateThumbnail(
	ctx context.Context,
	minioClient *minio.Client,
	asset *db.Asset,
	original *vips.Image,
) error {
	slog.Info("populating thumbnail media for asset", slog.String("id", asset.ID.String()))

	err := ctx.Err()
	if err != nil {
		return fmt.Errorf("context cancelled: %w", err)
	}

	asset.ThumbnailWidth = int32((original.Width() * THUMBNAIL_HEIGHT) / original.Height())
	asset.ThumbnailHeight = THUMBNAIL_HEIGHT

	if original.Pages() == 1 {
		asset.Thumbnail = asset.Original

		return nil
	}

	asset.ThumbnailWidth = int32((original.Width() * THUMBNAIL_HEIGHT) / original.PageHeight())

	thumbnail, err := createThumbnailForAnimationImage(original, err)
	if err != nil {
		return err
	}

	params := vips.DefaultWebpsaveBufferOptions()
	params.Q = THUMBNAIL_QUALITY

	buf, err := thumbnail.WebpsaveBuffer(params)
	if err != nil {
		return fmt.Errorf("unable to write preview image: %w", err)
	}

	asset.Thumbnail = createAssetKey()

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

func createThumbnailForAnimationImage(original *vips.Image, err error) (*vips.Image, error) {
	slog.Debug("original image",
		slog.Int("width", original.Width()),
		slog.Int("height", original.Height()),
		slog.Int("page_height", original.PageHeight()),
		slog.Int("loop", original.Loop()),
		slog.Int("pages", original.Pages()),
	)

	copyOptions := vips.DefaultCopyOptions()

	thumbnail, _ := original.Copy(copyOptions)

	defer thumbnail.Close()

	err = thumbnail.Autorot(nil)
	if err != nil {
		return nil, fmt.Errorf("unable to perform auto rotating: %w", err)
	}

	factor := float64(THUMBNAIL_HEIGHT) / float64(original.PageHeight())
	thumbnail.Resize(factor, &vips.ResizeOptions{
		Kernel: vips.KernelLanczos3,
		Gap:    2,
	})

	err = thumbnail.ExtractArea(0, 0, thumbnail.Width(), THUMBNAIL_HEIGHT)
	if err != nil {
		return nil, fmt.Errorf("unable to extract area: %w", err)
	}
	thumbnail.SetPages(1)
	thumbnail.SetPageHeight(THUMBNAIL_HEIGHT)

	slog.Debug("thumbnail image",
		slog.Int("width", thumbnail.Width()),
		slog.Int("height", thumbnail.Height()),
		slog.Int("page_height", thumbnail.PageHeight()),
		slog.Int("loop", thumbnail.Loop()),
		slog.Int("pages", thumbnail.Pages()),
	)
	return thumbnail, nil
}

func PopulateImageEmbedding(
	ctx context.Context,
	asset *db.Asset,
	original *vips.Image,
) error {
	slog.Info("populating image embedding for asset", slog.String("id", asset.ID.String()))
	spec, err := clip.GetImageSpec(ctx)
	if err != nil {
		return fmt.Errorf("unable to get image spec: %w", err)
	}
	copyOptions := vips.DefaultCopyOptions()
	img, _ := original.Copy(copyOptions)

	defer img.Close()

	err = img.Autorot(nil)
	if err != nil {
		return fmt.Errorf("unable to perform auto rotating: %w", err)
	}

	width := img.Width()
	pageHeight := img.PageHeight()

	img.ExtractArea(0, 0, width, pageHeight)
	img.SetPages(1)

	options := vips.DefaultThumbnailImageOptions()
	options.Height = int(spec.Height)
	options.Crop = vips.InterestingAttention
	options.Size = vips.SizeBoth

	err = img.ThumbnailImage(int(spec.Width), options)
	if err != nil {
		return fmt.Errorf("unable to resize image: %w", err)
	}

	buff, err := img.WebpsaveBuffer(vips.DefaultWebpsaveBufferOptions())
	if err != nil {
		return fmt.Errorf("unable to save image: %w", err)
	}

	resp, err := clip.EncodeImage(ctx, buff)
	if err != nil {
		return fmt.Errorf("unable to get image embedding: %w", err)
	}

	embedding, err := ParseNumpyBytes(resp.Embedding)
	if err != nil {
		return fmt.Errorf("unable to decode embedding: %w", err)
	}
	asset.ImageEmbedding = &embedding

	return nil
}

func ParseNumpyBytes(b []byte) (pgvector.Vector, error) {
	// 4 bytes per float32
	length := len(b) / 4
	vec := make([]float32, length)

	for i := range length {
		bits := binary.LittleEndian.Uint32(b[i*4 : (i+1)*4])
		vec[i] = math.Float32frombits(bits)
	}

	return pgvector.NewVector(vec), nil
}
