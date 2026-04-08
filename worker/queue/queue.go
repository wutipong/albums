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
	queue, err = neoq.New(ctx, neoq.WithBackend(memory.Backend))
	if err != nil {
		return fmt.Errorf("unable to initialize queue")
	}

	// create a handler that listens for new job on the "greetings" queue
	h := handler.New("asset-processing", func(ctx context.Context) (err error) {
		j, _ := jobs.FromContext(ctx)
		id := j.Payload["id"]
		idStr, ok := id.(string)

		if !ok {
			return fmt.Errorf("invalid id")
		}

		var uuid pgtype.UUID
		uuid.Scan(id)
		queries := db.New(db.Connection())
		asset, err := queries.GetAsset(ctx, uuid)
		if err != nil {
			return fmt.Errorf("unable to read asset data: %w", err)
		}

		if asset.ProcessStatus != db.ProcessStatusTPending {
			return nil
		}

		asset.ProcessStatus = db.ProcessStatusTProcessing

		_, err = queries.UpdateAsset(ctx, db.UpdateAssetParams{
			ID:            uuid,
			Filename:      asset.Filename,
			Checksum:      asset.Checksum,
			Type:          asset.Type,
			Original:      asset.Original,
			Preview:       asset.Preview,
			Thumbnail:     asset.Thumbnail,
			View:          asset.View,
			ProcessStatus: asset.ProcessStatus,
		})

		if err != nil {
			return fmt.Errorf("unable to save image metadata: %w", err)
		}

		ProcessAsset(ctx, idStr)

		//done <- true
		return
	})
	return queue.Start(ctx, h)
}

func Shutdown(ctx context.Context) {
	queue.Shutdown(ctx)
	vips.Shutdown()
}

func EnqueueAssetProcessing(ctx context.Context, id string) (status db.ProcessStatusT, err error) {
	var uuid pgtype.UUID

	uuid.Scan(id)
	quries := db.New(db.Connection())
	status, err = quries.GetAssetProcessStatus(ctx, uuid)

	if status != db.ProcessStatusTPending {
		return
	}

	status = db.ProcessStatusTProcessing
	_, err = quries.UpdateAssetProcessStatus(ctx,
		db.UpdateAssetProcessStatusParams{
			ID:            uuid,
			ProcessStatus: status,
		},
	)

	if err != nil {
		err = fmt.Errorf("unable to update status: %w", err)
		return
	}

	j := &jobs.Job{Queue: "asset-processing", Payload: map[string]any{"id": id}}

	jobId, err := queue.Enqueue(context.Background(), j)
	if err != nil {
		err = fmt.Errorf("unable to add job: %w", err)
		return
	}

	slog.Info("job added", slog.String("job", jobId))

	return
}
