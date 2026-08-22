package queue

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/wutipong/albums/worker/db"
)

func EnqueueProcessAllAssets(ctx context.Context, onlyMissing bool) error {
	queries, _ := db.Get()

	var assets []db.Asset
	var err error

	if onlyMissing {
		assets, err = queries.GetPendingOrFailedAssets(ctx)
	} else {
		assets, err = queries.GetAssetsWithoutUploading(ctx)
	}
	if err != nil {
		return fmt.Errorf("unable to query pending items: %w", err)
	}

	slog.Info("scan library for unprocessed asset.")

	slog.Info("pending tasks found", slog.Int("count", len(assets)))
	if len(assets) == 0 {
		return nil
	}

	for _, asset := range assets {
		slog.Info("adding asset", slog.String("id", asset.ID.String()))

		EnqueueAssetProcessing(ctx, asset.ID.String())
	}

	return nil
}
