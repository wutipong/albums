package queue

import (
	"context"
	"encoding/binary"
	"fmt"
	"log/slog"
	"math"
	"os"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/minio/minio-go/v7"
	"github.com/pgvector/pgvector-go"
	"github.com/wutipong/albums/worker/clip"
	"github.com/wutipong/albums/worker/db"
	vips "github.com/wutipong/albums/worker/vips"
)

func PopulateImageEmbedding(
	ctx context.Context, minioClient *minio.Client, id string,
) error {
	slog.Info("populating image embedding for asset", slog.String("id", id))

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

	if asset.Type != db.AssetTypeTImage {
		slog.Info("Embedding only supports image asset type")
		return nil
	}

	err = DoPopulateEmbedding(ctx, minioClient, &asset)
	if err != nil {
		return err
	}

	_, err = queries.UpdateAsset(ctx, db.UpdateAssetParams{
		ID:              uuid,
		Filename:        asset.Filename,
		Type:            asset.Type,
		Original:        asset.Original,
		Preview:         asset.Preview,
		Thumbnail:       asset.Thumbnail,
		View:            asset.View,
		ProcessStatus:   asset.ProcessStatus,
		ThumbnailWidth:  asset.ThumbnailWidth,
		ThumbnailHeight: asset.ThumbnailHeight,
		ViewWidth:       asset.ViewWidth,
		ViewHeight:      asset.ViewHeight,
		ImageFrames:     asset.ImageFrames,
		VideoDuration:   asset.VideoDuration,
		ImageEmbedding:  asset.ImageEmbedding,
	})

	if err != nil {
		slog.Error("update asset fails.", slog.String("error", err.Error()))
		return fmt.Errorf("unable to save asset metadata: %w", err)
	}

	return nil
}

func DoPopulateEmbedding(ctx context.Context, minioClient *minio.Client, asset *db.Asset) error {
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

	source := vips.NewSource(object)
	defer source.Close()

	slog.Info("read original image file.")

	params := vips.DefaultLoadOptions()

	original, err := vips.NewImageFromSource(source, params)
	if err != nil {
		return fmt.Errorf("unable to read original image: %w", err)
	}
	defer original.Close()

	embedding, err := GetImageEmbedding(ctx, original)
	if err != nil {
		return err
	}
	asset.ImageEmbedding = &embedding
	return nil
}

func GetImageEmbedding(ctx context.Context, original *vips.Image) (pgvector.Vector, error) {
	spec, err := clip.GetImageSpec(ctx)
	if err != nil {
		return pgvector.Vector{}, fmt.Errorf("unable to get image spec: %w", err)
	}
	copyOptions := vips.DefaultCopyOptions()
	img, _ := original.Copy(copyOptions)

	defer img.Close()

	err = img.Autorot(nil)
	if err != nil {
		return pgvector.Vector{}, fmt.Errorf("unable to perform auto rotating: %w", err)
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
		return pgvector.Vector{}, fmt.Errorf("unable to resize image: %w", err)
	}

	buff, err := img.WebpsaveBuffer(vips.DefaultWebpsaveBufferOptions())
	if err != nil {
		return pgvector.Vector{}, fmt.Errorf("unable to save image: %w", err)
	}

	resp, err := clip.EncodeImage(ctx, buff)
	if err != nil {
		return pgvector.Vector{}, fmt.Errorf("unable to get image embedding: %w", err)
	}

	embedding, err := ParseNumpyBytes(resp.Embedding)
	if err != nil {
		return pgvector.Vector{}, fmt.Errorf("unable to decode embedding: %w", err)
	}
	return embedding, nil
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
