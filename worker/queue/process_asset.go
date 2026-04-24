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

const THUMBNAIL_QUALITY = 60

var animationExts = []string{
	".gif",
	".webm",
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
			slog.Info("error processing image asset.", slog.String("error", err.Error()))
			return fmt.Errorf("unable to process image asset: %w", err)
		}
	case "video":
		err = processVideoAsset(ctx, minioClient, &asset)
		if err != nil {
			slog.Info("error proessing video.", slog.String("error", err.Error()))
			return fmt.Errorf("unable to process video asset: %w", err)
		}
	default:
		asset.DeletedAt = pgtype.Timestamptz{Time: time.Now(), Valid: true}
		slog.Info("asset not recongnized will be deleted")
	}

	_, err = queries.UpdateAsset(ctx, db.UpdateAssetParams{
		ID:              uuid,
		Filename:        asset.Filename,
		Type:            asset.Type,
		Original:        asset.Original,
		Preview:         asset.Preview,
		Thumbnail:       asset.Thumbnail,
		View:            asset.View,
		ProcessStatus:   db.ProcessStatusTProcessed,
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
		return fmt.Errorf("unable to save image metadata: %w", err)
	}

	return nil
}

func createAssetKey() string {
	return fmt.Sprintf("public/%s", uuid.NewString())
}
