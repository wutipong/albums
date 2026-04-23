package queue

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/wutipong/albums/worker/db"
)

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
		return fmt.Errorf("unable to parse album id: %w", err)
	}

	album, err := queries.GetAlbum(ctx, albumIdUUID)

	var asset db.Asset
	if assetId != "" {
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
	album.Cover = asset.Thumbnail

	_, err := queries.UpdateAlbumThumbnail(ctx, db.UpdateAlbumThumbnailParams{
		ID:    album.ID,
		Cover: album.Cover,
	})

	if err != nil {
		return fmt.Errorf("unable update data: %w", err)
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
