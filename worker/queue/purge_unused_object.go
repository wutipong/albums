package queue

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/bits-and-blooms/bloom/v3"
	"github.com/minio/minio-go/v7"
	"github.com/wutipong/albums/worker/db"
)

const OBJECT_LIMIT = 100

func PurgeUnusedObjectfunc(
	ctx context.Context, minioClient *minio.Client,
) error {

	slog.Info("Purge unused object from s3")
	queries, _ := db.Get()

	filter := bloom.NewWithEstimates(1000000, 0.01)
	for i := 0; ; i += OBJECT_LIMIT {
		params := db.GetAssetsWithObjectsParams{
			Limit:  OBJECT_LIMIT,
			Offset: int64(i),
		}
		assets, err := queries.GetAssetsWithObjects(ctx, params)
		if err != nil {
			return fmt.Errorf("unable to retrive asset value: %w", err)
		}
		if len(assets) == 0 {
			break
		}
		for _, asset := range assets {
			filter.AddString(asset.Original)
			filter.AddString(asset.Preview)
			filter.AddString(asset.View)
			filter.AddString(asset.Thumbnail)
		}
	}

	objectsCh := make(chan minio.ObjectInfo)
	go func() {
		defer close(objectsCh)
		for object := range minioClient.ListObjects(
			ctx, os.Getenv("S3_BUCKET"),
			minio.ListObjectsOptions{Recursive: true},
		) {
			slog.Debug("current object", slog.String("key", object.Key))
			if object.Err != nil {
				slog.Error("error reading object", slog.String("key", object.Key))
				continue
			}
			if filter.TestString(object.Key) {
				slog.Debug("object key found, skip", slog.String("key", object.Key))

				continue
			}

			if exitst, err := queries.IsObjectInUse(ctx, object.Key); exitst || err != nil {
				slog.Debug("object key found in the database, skip", slog.String("key", object.Key))
				continue
			}

			slog.Debug("add object to delete batch", slog.String("key", object.Key))

			objectsCh <- object
		}
	}()

	opts := minio.RemoveObjectsOptions{}
	for rErr := range minioClient.RemoveObjects(ctx, os.Getenv("S3_BUCKET"), objectsCh, opts) {
		if rErr.Err != nil {
			slog.Error("failed to delete object",
				slog.String("key", rErr.ObjectName),
				slog.String("error", rErr.Error()),
			)
		}
	}
	return nil
}
