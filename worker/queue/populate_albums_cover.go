package queue

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/minio/minio-go/v7"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"github.com/wutipong/albums/worker/db"
	vips "github.com/wutipong/albums/worker/vips"
)

const COVER_WIDTH = 300
const COVER_HEIGHT_PORTRAIT = 450
const COVER_HEIGHT_LANDSCAPE = 200

func PopulateAlbumCover(
	ctx context.Context,
	minioClient *minio.Client,
	albumId string,
	assetId string,
) error {
	queries, _ := db.Get()

	var albumIdUUID pgtype.UUID
	err := albumIdUUID.Scan(albumId)
	if err != nil {
		return fmt.Errorf("unable to parse album id: %w", err)
	}

	album, err := queries.GetAlbum(ctx, albumIdUUID)

	var assetIdUUID pgtype.UUID
	err = assetIdUUID.Scan(assetId)
	if err != nil {
		return fmt.Errorf("unable to parse asset id: %w", err)
	}
	asset, err := queries.GetAsset(ctx, assetIdUUID)

	if err != nil {
		return fmt.Errorf("unable to find asset :%w", err)
	}

	err = SetAlbumCoverFromAsset(ctx, minioClient, queries, asset, album)

	if err != nil {
		slog.Error(
			"unable to update album cover",
			slog.Any("id", album.ID),
			slog.String("error", err.Error()),
		)
	}

	return nil
}

func SetAlbumCoverFromAsset(
	ctx context.Context,
	minioClient *minio.Client,
	queries *db.Queries,
	asset db.Asset,
	album db.Album,
) error {
	slog.Info("getting object from S3.", slog.String("id", asset.Original))
	object, err := getObjectFromS3(ctx, minioClient, asset.Original)
	if err != nil {
		return fmt.Errorf("unable to get object from s3: %w", err)
	}
	defer object.Close()

	var original *vips.Image = nil
	if asset.Type == db.AssetTypeTImage {
		source := vips.NewSource(object)
		defer source.Close()

		slog.Info("read original image file.")

		params := vips.DefaultLoadOptions()
		if hasAnimationExt(filepath.Ext(asset.Filename)) {
			params.N = -1
		}

		original, err = vips.NewImageFromSource(source, params)
		if err != nil {
			return fmt.Errorf("unable to read original image: %w", err)
		}
	} else {
		originalFile, err := os.CreateTemp("",
			fmt.Sprintf("*.%s", filepath.Base(asset.Filename)),
		)

		if err != nil {
			return fmt.Errorf("unable to create temp file for original asset: %w", err)
		}
		defer os.Remove(originalFile.Name())

		io.Copy(originalFile, object)

		probe, err := ffmpeg.Probe(originalFile.Name())
		if err != nil {
			return fmt.Errorf("unable to probe original video: %w", err)
		}

		var info Probe
		json.Unmarshal([]byte(probe), &info)

		original, err = extractVideoThumbnail(info, originalFile, &asset, "vframes", "1")
		if err != nil {
			return fmt.Errorf("unable to read original image: %w", err)
		}
	}
	defer original.Close()

	slog.Debug("original image",
		slog.Int("width", original.Width()),
		slog.Int("height", original.Height()),
		slog.Int("page_height", original.PageHeight()),
		slog.Int("loop", original.Loop()),
		slog.Int("pages", original.Pages()),
	)

	copyOptions := vips.DefaultCopyOptions()
	cover, err := original.Copy(copyOptions)
	if err != nil {
		return fmt.Errorf("unable to copy from original image: %w", err)
	}

	defer cover.Close()

	err = cover.Autorot(nil)
	if err != nil {
		return fmt.Errorf("unable to perform auto rotating: %w", err)
	}

	targetHeight := COVER_HEIGHT_LANDSCAPE
	if os.Getenv("COVER_ASPECT") == "portrait" {
		targetHeight = COVER_HEIGHT_LANDSCAPE
	}

	factor := float64(targetHeight) / float64(original.PageHeight())
	cover.Resize(factor, &vips.ResizeOptions{
		Kernel: vips.KernelLanczos3,
		Gap:    2,
	})

	err = cover.ExtractArea(0, 0, COVER_WIDTH, targetHeight)
	if err != nil {
		return fmt.Errorf("unable to extract area: %w", err)
	}

	if original.Pages() > 1 {
		cover.SetPages(1)
		cover.SetPageHeight(targetHeight)
	}
	slog.Debug("thumbnail image",
		slog.Int("width", cover.Width()),
		slog.Int("height", cover.Height()),
		slog.Int("page_height", cover.PageHeight()),
		slog.Int("loop", cover.Loop()),
		slog.Int("pages", cover.Pages()),
	)

	saveParams := vips.DefaultWebpsaveBufferOptions()
	saveParams.Q = THUMBNAIL_QUALITY

	buf, err := cover.WebpsaveBuffer(saveParams)
	if err != nil {
		return fmt.Errorf("unable to write preview image: %w", err)
	}

	if album.Cover == "" {
		album.Cover = createAssetKey("webp")
	}

	err = putObjectToS3(ctx, minioClient, album.Cover, bytes.NewReader(buf), "image/webp")

	if err != nil {
		return fmt.Errorf("unable to put object to S3: %w", err)
	}
	_, err = queries.UpdateAlbumThumbnail(ctx, db.UpdateAlbumThumbnailParams{
		ID:    album.ID,
		Cover: album.Cover,
	})

	if err != nil {
		return fmt.Errorf("unable update data: %w", err)
	}

	return nil
}

func createAlbumCover(original *vips.Image) (thumbnail *vips.Image, err error) {
	slog.Debug("original image",
		slog.Int("width", original.Width()),
		slog.Int("height", original.Height()),
		slog.Int("page_height", original.PageHeight()),
		slog.Int("loop", original.Loop()),
		slog.Int("pages", original.Pages()),
	)

	copyOptions := vips.DefaultCopyOptions()
	thumbnail, err = original.Copy(copyOptions)
	if err != nil {
		err = fmt.Errorf("unable to copy from original image: %w", err)
		return
	}

	err = thumbnail.Autorot(nil)
	if err != nil {
		err = fmt.Errorf("unable to perform auto rotating: %w", err)
		return
	}

	factor := float64(THUMBNAIL_HEIGHT) / float64(original.PageHeight())
	thumbnail.Resize(factor, &vips.ResizeOptions{
		Kernel: vips.KernelLanczos3,
		Gap:    2,
	})

	err = thumbnail.ExtractArea(0, 0, thumbnail.Width(), THUMBNAIL_HEIGHT)
	if err != nil {
		err = fmt.Errorf("unable to extract area: %w", err)
		return
	}

	if original.Pages() > 1 {
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
	return
}
