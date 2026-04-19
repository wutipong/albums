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

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	vips "github.com/cshum/vipsgen/vips816"
	"github.com/pgvector/pgvector-go"
	"github.com/wutipong/albums/worker/clip"
	"github.com/wutipong/albums/worker/db"
)

const THUMBNAIL_HEIGHT = 200

func processImageAsset(ctx context.Context, s3Client *s3.Client, asset *db.Asset) error {
	slog.Info("processing image asset", slog.String("id", asset.ID.String()))

	err := ctx.Err()
	if err != nil {
		slog.Info("context.", slog.String("error", err.Error()))
		return fmt.Errorf("context cancelled: %w", err)
	}

	slog.Info("getting object from S3.", slog.String("id", asset.Original))
	s3Obj, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Key:    aws.String(asset.Original),
		Bucket: aws.String(os.Getenv("S3_BUCKET")),
	})

	if err != nil {
		return fmt.Errorf("unable to get object from s3: %w", err)
	}
	defer s3Obj.Body.Close()

	slog.Info("reading data", slog.String("id", asset.Original))
	buff, err := io.ReadAll(s3Obj.Body)
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

	err = populateView(ctx, s3Client, asset, original)
	if err != nil {
		return fmt.Errorf("unable to populate view image: %e", err)
	}

	err = populatePreview(ctx, s3Client, asset, original)
	if err != nil {
		return fmt.Errorf("unable to populate preview image: %e", err)
	}

	err = populateThumbnail(ctx, s3Client, asset, original)
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
	s3Client *s3.Client,
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

	_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("S3_BUCKET")),
		Body:   bytes.NewReader(buf),
		Key:    aws.String(asset.View),
	})

	if err != nil {
		return fmt.Errorf("unable to put object to S3: %w", err)
	}

	return nil
}

func populatePreview(
	ctx context.Context,
	_ *s3.Client,
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

	asset.Preview = asset.Original
	asset.ImageFrames = int32(original.Pages())

	return nil
}

func populateThumbnail(
	ctx context.Context,
	s3Client *s3.Client,
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

	copyOptions := vips.DefaultCopyOptions()

	thumbnail, _ := original.Copy(copyOptions)

	defer thumbnail.Close()

	err = thumbnail.Autorot(nil)
	if err != nil {
		return fmt.Errorf("unable to perform auto rotating: %w", err)
	}

	width := thumbnail.Width()
	pageHeight := thumbnail.PageHeight()

	thumbnail.ExtractArea(0, 0, width, pageHeight)
	thumbnail.SetPages(1)
	params := vips.DefaultWebpsaveBufferOptions()
	params.Q = THUMBNAIL_QUALITY

	buf, err := thumbnail.WebpsaveBuffer(params)
	if err != nil {
		return fmt.Errorf("unable to write preview image: %w", err)
	}

	asset.Thumbnail = createAssetKey()

	_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("S3_BUCKET")),
		Body:   bytes.NewReader(buf),
		Key:    aws.String(asset.Thumbnail),
	})

	if err != nil {
		return fmt.Errorf("unable to put object to S3: %w", err)
	}

	return nil
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
