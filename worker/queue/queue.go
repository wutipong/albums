package queue

import (
	"context"
	"fmt"
	"log/slog"
	"os"

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

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewEnvAWS(),
		Secure: secure,
	})
	if err != nil {
		return fmt.Errorf("unable to create minio client: %w", err)
	}

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
				j.Payload["albumId"].(string),
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

	status, err = queries.GetAssetProcessStatus(ctx, uuid)
	slog.Info("asset status", slog.Any("status", status))

	slog.Info("enqueueing asset", slog.String("id", id))

	if status != db.ProcessStatusTPending {
		return
	}

	j := &jobs.Job{
		Queue: "asset-processing",
		Payload: map[string]any{
			"command": "process-asset",
			"id":      id,
		},
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

func EnqueuePopulateAlbumsCover(ctx context.Context, albumId string, assetId string) error {
	j := &jobs.Job{
		Queue: "asset-processing",
		Payload: map[string]any{
			"command": "populate-album-cover",
			"albumId": albumId,
			"assetId": assetId,
		},
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
