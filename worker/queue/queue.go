package queue

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/acaloiaro/neoq"
	"github.com/acaloiaro/neoq/backends/memory"
	"github.com/acaloiaro/neoq/handler"
	"github.com/acaloiaro/neoq/jobs"
	"github.com/davidbyttow/govips/v2/vips"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/wutipong/albums/worker/db"
)

var queue neoq.Neoq

// var done = make(chan bool)

func Init(ctx context.Context) error {
	var err error
	err = vips.Startup(nil)
	if err != nil {
		return fmt.Errorf("unable to initialize vips")
	}
	queue, err = neoq.New(ctx,
		neoq.WithBackend(memory.Backend),
		// neoq.WithRecoveryCallback(func(ctx context.Context, err error) error {
		// 	slog.Error("error processing task", slog.String("error", err.Error()))
		// 	return nil
		// }),
	)
	if err != nil {
		return fmt.Errorf("unable to initialize queue")
	}

	// create a handler that listens for new job on the "greetings" queue
	h := handler.New("asset-processing", func(ctx context.Context) (err error) {

		j, _ := jobs.FromContext(ctx)
		id := j.Payload["id"]
		slog.Info("processing asset", slog.Any("id", id))

		idStr, ok := id.(string)
		if !ok {
			return fmt.Errorf("invalid id")
		}

		var uuid pgtype.UUID
		err = uuid.Scan(idStr)
		if err != nil {
			slog.Info("error parsing uuid", slog.String("id", idStr))
			return fmt.Errorf("unable to parse id: %w", err)
		}
		queries, _ := db.Get()

		asset, err := queries.GetAsset(ctx, uuid)
		if err != nil {
			slog.Info("error reading asset.", slog.String("id", idStr))
			return fmt.Errorf("unable to read asset data: %w", err)
		}

		if asset.ProcessStatus != db.ProcessStatusTPending {
			return nil
		}

		asset.ProcessStatus = db.ProcessStatusTProcessing

		_, err = queries.UpdateAssetProcessStatus(ctx,
			db.UpdateAssetProcessStatusParams{
				ID:            uuid,
				ProcessStatus: db.ProcessStatusTProcessing,
			},
		)

		if err != nil {
			return fmt.Errorf("unable to save image metadata: %w", err)
		}

		ProcessAsset(ctx, idStr, queries)

		//done <- true
		return
	}, handler.Concurrency(1))
	return queue.Start(ctx, h)
}

func Shutdown(ctx context.Context) {
	queue.Shutdown(ctx)
	vips.Shutdown()
}

func EnqueueAssetProcessing(ctx context.Context, id string) (status db.ProcessStatusT, err error) {
	var uuid pgtype.UUID

	uuid.Scan(id)

	{
		queries, _ := db.Get()

		status, err = queries.GetAssetProcessStatus(ctx, uuid)
		slog.Info("asset status", slog.Any("status", status))

		slog.Info("enqueueing asset", slog.String("id", id))

		if status != db.ProcessStatusTPending {
			return
		}
	}

	j := &jobs.Job{Queue: "asset-processing", Payload: map[string]any{"id": id}}

	jobId, err := queue.Enqueue(ctx, j)
	if err != nil {
		err = fmt.Errorf("unable to add job: %w", err)
		return
	}

	slog.Info("job added", slog.String("job", jobId))

	return
}
