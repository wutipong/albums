package queue

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/wutipong/albums/worker/db"
)

const ALBUM_COVER_FILE = "cover.webp"

func PopulateAlbumCover(ctx context.Context, albumId string, assetId string) error {
	queries, _ := db.Get()

	var albumIdUUID pgtype.UUID
	err := albumIdUUID.Scan(albumId)
	if err != nil {
		return fmt.Errorf("unable to parse album id: %w", err)
	}

	album, err := queries.GetAlbum(ctx, albumIdUUID)

	var asset db.Asset
	if assetId == "" {
		if album.Cover != "" {
			slog.Info("album already has a cover, skipping",
				slog.String("albumId", albumId),
				slog.String("cover", album.Cover),
			)
			return nil
		}

		if os.Getenv("COVER_ASPECT") == "portrait" {
			asset, err = queries.GetAlbumPortraitAssetForCover(ctx, albumIdUUID)
		} else {
			asset, err = queries.GetAlbumLandscapeAssetForCover(ctx, albumIdUUID)
		}
		if errors.Is(err, sql.ErrNoRows) {
			asset, err = queries.GetAlbumAssetForCover(ctx, albumIdUUID)
		}
		if err != nil {
			return fmt.Errorf("unable to find asset for cover: %w", err)
		}
	} else {
		var assetIdUUID pgtype.UUID
		err = assetIdUUID.Scan(assetId)
		if err != nil {
			return fmt.Errorf("unable to parse asset id: %w", err)
		}
		asset, err = queries.GetAsset(ctx, assetIdUUID)
	}
	if err != nil {
		return fmt.Errorf("unable to find asset :%w", err)
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
