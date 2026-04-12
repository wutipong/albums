package queue

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/jackc/pgx/v5/pgtype"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"github.com/wutipong/albums/worker/db"
)

const ALBUM_COVER_WIDTH = 300
const ALBUM_COVER_HEIGHT = 200
const ALBUM_COVER_FILE = "cover.webp"

func PopulateAlbumsCover(ctx context.Context) error {
	queries, _ := db.Get()

	albums, err := queries.GetAlbumsWithoutCover(ctx)
	if err != nil {
		return fmt.Errorf("unable to populate albums without cover image: %w", err)
	}

	for _, album := range albums {
		err = PopulateAlbumCover(ctx, album.ID.String(), "")

		if err != nil {
			slog.Error(
				"unable to update album cover",
				slog.Any("id", album.ID),
				slog.String("error", err.Error()),
			)
		}

	}

	return nil
}

func PopulateAlbumCover(ctx context.Context, albumId string, assetId string) error {
	queries, _ := db.Get()

	var albumIdUUID pgtype.UUID
	err := albumIdUUID.Scan(albumId)
	if err != nil {
		return fmt.Errorf("unable to parse album id: %w")
	}

	album, err := queries.GetAlbum(ctx, albumIdUUID)

	var asset db.Asset
	if albumId == "" {
		var assetIdUUID pgtype.UUID
		err = assetIdUUID.Scan(assetId)
		if err != nil {
			return fmt.Errorf("unable to parse asset id: %w", err)
		}
		asset, err = queries.GetAsset(ctx, assetIdUUID)
	} else {
		asset, err = queries.GetRandomAlbumAsset(ctx, album.ID)
		if err != nil {
			return fmt.Errorf("unable to retrive random asset: %w", err)

		}
	}

	err = SetAlbumCoverFromAsset(ctx, queries, asset, album)

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
	queries *db.Queries,
	asset db.Asset,
	album db.Album,
) error {
	var err error
	switch asset.Type {
	case "image":
		err = populateAlbumCoverFromImageAsset(ctx, &album, &asset)
	case "video":
		err = populateAlbumCoverFromVideoAsset(ctx, &album, &asset)
	}

	if err != nil {
		return fmt.Errorf("unable to populate cover from asset: %w", err)
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

func populateAlbumCoverFromImageAsset(
	ctx context.Context,
	album *db.Album,
	asset *db.Asset,
) error {
	err := ctx.Err()
	if err != nil {
		slog.Info("context.", slog.String("error", err.Error()))
		return fmt.Errorf("context cancelled: %w", err)
	}

	assetId := asset.ID.String()
	originalPath := createCacheAssetPath(assetId, asset.Original)
	slog.Info("original asset path", slog.String("path", originalPath))

	params := vips.NewImportParams()
	cover, err := vips.LoadImageFromFile(originalPath, params)
	if err != nil {
		return fmt.Errorf("unable to read image: %w", err)
	}

	return writeAlbumCover(err, cover, album)
}

func populateAlbumCoverFromVideoAsset(ctx context.Context, album *db.Album, asset *db.Asset) error {
	err := ctx.Err()
	if err != nil {
		slog.Info("context.", slog.String("error", err.Error()))
		return fmt.Errorf("context cancelled: %w", err)
	}

	assetId := asset.ID.String()
	originalPath := createCacheAssetPath(assetId, asset.Original)
	slog.Info("original asset path", slog.String("path", originalPath))

	probe, err := ffmpeg.Probe(originalPath)
	if err != nil {
		return fmt.Errorf("unable to probe original video: %w", err)
	}

	var info Probe
	json.Unmarshal([]byte(probe), &info)

	buffer := new(bytes.Buffer)
	err = ffmpeg.Input(originalPath).
		WithOutput(buffer).
		Output("pipe:", ffmpeg.KwArgs{
			"c:v":     "libwebp",
			"f":       "webp",
			"vframes": "1",
			"quality": fmt.Sprintf("%d", THUMBNAIL_QUALITY),
		}).OverWriteOutput().ErrorToStdOut().Run()

	image, err := vips.NewImageFromReader(buffer)
	if err != nil {
		return fmt.Errorf("unable to read image from ffmpeg output: %w", err)
	}

	slog.Debug("image", slog.Any("image", image))

	return writeAlbumCover(err, image, album)
}

func writeAlbumCover(err error, cover *vips.ImageRef, album *db.Album) error {
	err = cover.AutoRotate()
	if err != nil {
		return fmt.Errorf("unable to perform auto rotating: %w", err)
	}

	err = cover.ThumbnailWithSize(
		ALBUM_COVER_WIDTH,
		ALBUM_COVER_HEIGHT,
		vips.InterestingAttention,
		vips.SizeBoth,
	)
	if err != nil {
		return fmt.Errorf("unable to resize preview image: %w", err)
	}

	album.Cover = ALBUM_COVER_FILE
	coverPath := createCacheAlbumPath(album.ID.String(), album.Cover)
	err = os.MkdirAll(filepath.Dir(coverPath), 0755)
	if err != nil {
		return fmt.Errorf("unable to reate directory: %w", err)
	}

	saveParams := vips.NewWebpExportParams()
	saveParams.Quality = THUMBNAIL_QUALITY

	buf, _, err := cover.ExportWebp(saveParams)
	if err != nil {
		return fmt.Errorf("unable to write preview image: %w", err)
	}

	err = os.WriteFile(coverPath, buf, 0644)
	if err != nil {
		return fmt.Errorf("unable to save file: %w", err)
	}
	return nil
}

func createCacheAlbumPath(id string, args ...string) string {
	topLevelDir := id[0:2]
	secondLevelDir := id[2:4]

	combined := []string{
		os.Getenv("CACHE_DIR"),
		"albums",
		topLevelDir,
		secondLevelDir,
		id,
	}
	combined = append(combined, args...)

	return filepath.Join(combined...)
}
