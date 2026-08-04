package queue

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/minio/minio-go/v7"
	"github.com/wutipong/albums/worker/db"
)

const THUMBNAIL_QUALITY = 80

var animationExts = []string{
	".gif",
	".webp",
}

func hasAnimationExt(ext string) bool {
	return slices.Contains(animationExts, strings.ToLower(ext))
}

func ProcessAsset(ctx context.Context, minioClient *minio.Client, id string) error {
	slog.Info("processing asset", slog.String("id", id))
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

	switch asset.Type {
	case "image":
		err = processImageAsset(ctx, minioClient, &asset)
		if err != nil {
			slog.Error("error processing image asset.",
				slog.String("error", err.Error()),
				slog.String("id", id),
			)
		}
	case "video":
		err = processVideoAsset(ctx, minioClient, &asset)
		if err != nil {
			slog.Error("error processing video asset.",
				slog.String("error", err.Error()),
				slog.String("id", id),
			)
		}
	default:
		asset.DeletedAt = pgtype.Timestamptz{Time: time.Now(), Valid: true}
		slog.Info("asset not recongnized will be deleted")
	}

	if err == nil {
		asset.ProcessStatus = db.ProcessStatusTProcessed
	} else {
		asset.ProcessStatus = db.ProcessStatusTFailed
	}

	_, e := queries.UpdateAsset(ctx, db.UpdateAssetParams{
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

	if e != nil {
		slog.Error("update asset fails.", slog.String("error", e.Error()))
		return fmt.Errorf("unable to save asset metadata: %w", e)
	}

	if err != nil {
		return fmt.Errorf("unable to process asset: %w", err)
	}

	slog.Info("asset processed successfully", slog.String("id", id))

	return nil
}

func createAssetKey(extension string) string {
	return fmt.Sprintf("public/%s.%s", uuid.NewString(), extension)
}
