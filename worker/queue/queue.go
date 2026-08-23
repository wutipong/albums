package queue

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/acaloiaro/neoq"
	"github.com/acaloiaro/neoq/backends/postgres"
	"github.com/acaloiaro/neoq/handler"
	"github.com/acaloiaro/neoq/jobs"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/wutipong/albums/worker/db"
)

var queue neoq.Neoq

// var done = make(chan bool)

const MAX_RETRIES = 2 // maximum number of retries for a job

func Init(ctx context.Context) error {
	var err error

	queue, err = neoq.New(ctx,
		neoq.WithBackend(postgres.Backend),
		postgres.WithConnectionString(os.Getenv("DATABASE_URL")),
	)
	if err != nil {
		return fmt.Errorf("unable to initialize queue")
	}

	endpoint, secure, err := GetMinioEndpoint(os.Getenv("AWS_ENDPOINT_URL"))
	if err != nil {
		return fmt.Errorf("unable to get endpoint: %w", err)
	}

	slog.Info("using S3 endpoint.", slog.String("endpoint", endpoint))

	accessKeyId := os.Getenv("AWS_ACCESS_KEY_ID")
	secret := os.Getenv("AWS_SECRET_ACCESS_KEY")

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:        credentials.NewStaticV4(accessKeyId, secret, ""),
		Secure:       secure,
		BucketLookup: minio.BucketLookupPath,
	})

	// create a handler that listens for new job on the "greetings" queue
	h := handler.New("asset-processing", func(ctx context.Context) (err error) {
		j, _ := jobs.FromContext(ctx)
		command := j.Payload["command"]
		switch command {
		case "process-asset":
			{
				id := j.Payload["id"]
				idStr := id.(string)
				slog.Info("job", slog.Any("id", id), slog.Any("command", command))
				err = ProcessAsset(ctx, minioClient, idStr)
			}

		case "populate-album-cover":
			err = PopulateAlbumCover(
				ctx,
				minioClient,
				j.Payload["albumId"].(string),
				j.Payload["assetId"].(string),
			)

		case "populate-image-embedding":
			err = PopulateImageEmbedding(
				ctx,
				minioClient,
				j.Payload["assetId"].(string),
			)
		}

		if err != nil {
			slog.Error("failed to process asset:", slog.String("error", err.Error()))
		}

		//done <- true
		return
	}, handler.Concurrency(1))
	return queue.Start(ctx, h)
}

func Shutdown(ctx context.Context) {
	queue.Shutdown(ctx)
}

func EnqueueAssetProcessing(ctx context.Context, id string) (status db.ProcessStatusT, err error) {
	var uuid pgtype.UUID

	uuid.Scan(id)

	queries, _ := db.Get()
	maxRetries := MAX_RETRIES
	status, err = queries.GetAssetProcessStatus(ctx, uuid)
	slog.Info("asset status", slog.Any("status", status))

	slog.Info("enqueueing asset", slog.String("id", id))

	j := &jobs.Job{
		Queue: "asset-processing",
		Payload: map[string]any{
			"command":    "process-asset",
			"id":         id,
			"created_at": time.Now().UTC(),
		},
		MaxRetries: &maxRetries,
	}

	jobId, err := queue.Enqueue(ctx, j)
	if err != nil {
		err = fmt.Errorf("unable to add job: %w", err)
		return
	}

	slog.Info(
		"job added",
		slog.String("job", jobId),
		slog.String("command", "process-asset"),
	)

	return
}

func EnqueueUpdateAlbumCover(ctx context.Context, albumId string, assetId string) error {
	maxRetries := MAX_RETRIES
	j := &jobs.Job{
		Queue: "asset-processing",
		Payload: map[string]any{
			"command":    "populate-album-cover",
			"albumId":    albumId,
			"assetId":    assetId,
			"created_at": time.Now().UTC(),
		},
		MaxRetries: &maxRetries,
	}

	jobId, err := queue.Enqueue(ctx, j)
	if err != nil {
		return fmt.Errorf("unable to add job: %w", err)
	}

	slog.Info(
		"job added",
		slog.String("job", jobId),
		slog.String("albumId", albumId),
		slog.String("assetId", assetId),
		slog.String("command", "populate-album-cover"),
	)

	return nil
}

func EnqueuePopulateImageEmbedding(ctx context.Context, id string) error {
	maxRetries := MAX_RETRIES
	j := &jobs.Job{
		Queue: "asset-processing",
		Payload: map[string]any{
			"command":    "populate-image-embedding",
			"assetId":    id,
			"created_at": time.Now().UTC(),
		},
		MaxRetries: &maxRetries,
	}

	jobId, err := queue.Enqueue(ctx, j)
	if err != nil {
		return fmt.Errorf("unable to add job: %w", err)
	}

	slog.Info(
		"job added",
		slog.String("job", jobId),
		slog.String("assetId", id),
		slog.String("command", "populate-image-embedding"),
	)

	return nil
}

func EnqueuePopulateAlbumsCover(ctx context.Context, missingOnly bool) error {
	queries, _ := db.Get()

	var albums []db.Album
	var err error

	if missingOnly {
		albums, err = queries.GetAlbumsWithoutCover(ctx)
	} else {
		albums, err = queries.GetAllAlbum(ctx)
	}

	if err != nil {
		return fmt.Errorf("unable to populate albums without cover image: %w", err)
	}

	for _, album := range albums {
		portrait := os.Getenv("COVER_ASPECT") == "portrait"

		var asset db.Asset

		if portrait {
			asset, err = queries.GetAlbumPortraitAssetForCover(ctx, album.ID)
		} else {
			asset, err = queries.GetAlbumLandscapeAssetForCover(ctx, album.ID)
		}

		if errors.Is(err, sql.ErrNoRows) {
			asset, err = queries.GetAlbumAssetForCover(ctx, album.ID)
		}

		if err != nil {
			slog.Warn(
				"unable to find image asset for cover",
				slog.String("album_id", album.ID.String()),
				slog.String("name", album.Name),
			)

			continue
		}

		err = EnqueueUpdateAlbumCover(ctx, album.ID.String(), asset.ID.String())

		if err != nil {
			slog.Error(
				"unable add update album cover job",
				slog.Any("id", album.ID),
				slog.String("error", err.Error()),
			)
			return err
		}

	}

	return nil
}
