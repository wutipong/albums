package queue

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/minio/minio-go/v7"
	"github.com/wutipong/albums/worker/db"
)

func DeleteAlbum(ctx context.Context, minioClient *minio.Client, albumId string) error {
	queries, _ := db.Get()

	var albumIdUUID pgtype.UUID
	err := albumIdUUID.Scan(albumId)
	if err != nil {
		return fmt.Errorf("unable to parse album id: %w", err)
	}

	minio.ToObjectInfo(os.Getenv("S3_Bucket"))

	album, err := queries.GetAlbum(ctx, albumIdUUID)

	assets, err := queries.GetAlbumAssets(ctx, albumIdUUID)
	if err != nil {
		return fmt.Errorf("unable to get album assets: %w", err)
	}

	ch := make(chan minio.ObjectInfo)
	go func() {
		defer close(ch)

		for _, asset := range assets {
			ch <- minio.ObjectInfo{
				Key: asset.Original,
			}

			ch <- minio.ObjectInfo{
				Key: asset.View,
			}

			ch <- minio.ObjectInfo{
				Key: asset.Preview,
			}

			ch <- minio.ObjectInfo{
				Key: asset.Thumbnail,
			}
		}
	}()

	for rErr := range minioClient.RemoveObjects(ctx, os.Getenv("S3_BUCKET"), ch, minio.RemoveObjectsOptions{}) {
		if rErr.Err != nil {
			slog.Error("failed to delete object",
				slog.String("key", rErr.ObjectName),
				slog.String("error", rErr.Error()),
			)
		}
	}

	err = minioClient.RemoveObject(ctx, os.Getenv("S3_BUCKET"), album.Cover, minio.RemoveObjectOptions{})
	if err != nil {
		slog.Error("failed to delete object",
			slog.String("key", album.Cover),
			slog.String("error", err.Error()),
		)

		return fmt.Errorf("failed to delete album cover object (%s): %w", album.Cover, err)
	}

	err = queries.MarkAlbumDeleted(ctx, albumIdUUID)
	if err != nil {
		slog.Error("failed to mark album as deleted",
			slog.String("id", albumId),
			slog.String("error", err.Error()))
	}

	assetsIds := make([]pgtype.UUID, len(assets))
	for i, asset := range assets {
		assetsIds[i].Scan(asset.ID)
	}

	results := queries.MarkAssetsDeleted(ctx, assetsIds)
	defer results.Close()
	results.Exec(func(i int, err error) {
		if err != nil {
			slog.Error("failed to mark asset deleted",
				slog.String("id", assetsIds[i].String()),
				slog.String("err", err.Error()),
			)
		}
		slog.Debug("asset is marked as deleted",
			slog.String("id", assetsIds[i].String()),
		)
	})

	return nil
}
