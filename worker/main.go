package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/urfave/cli/v3"
	"github.com/wutipong/albums/worker/db"
	"github.com/wutipong/albums/worker/queue"
	"github.com/wutipong/albums/worker/service"
	"github.com/wutipong/albums/worker/service/pb"
	"google.golang.org/grpc"
)

//go:generate protoc --go_out=. --go-grpc_out=. -I/workspaces/grpc worker.proto clip.proto
//go:generate sqlc generate

func main() {
	slog.SetDefault(slog.New(tint.NewHandler(os.Stderr, &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: time.Kitchen,
	})))

	id := ""

	cmd := &cli.Command{
		Name:  "worker",
		Usage: "process assets to albums",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return performWork(ctx)
		},
		Commands: []*cli.Command{
			{
				Name:  "single",
				Usage: "immediately process single asset",
				Arguments: []cli.Argument{
					&cli.StringArg{
						Name:        "id",
						UsageText:   "id of the asset to process",
						Destination: &id,
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return processSingle(ctx, id)
				},
			}, {
				Name:  "populate-albums-cover",
				Usage: "update albums without cover with one from randomly picked asset.",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					err := db.Connect(ctx, os.Getenv("DATABASE_URL"))
					if err != nil {
						return fmt.Errorf("unable to connect to the database: %w", err)
					}
					defer db.Close(ctx)
					return queue.PopulateAlbumsCover(ctx)
				},
			}, {
				Name:  "populate-image-embedding",
				Usage: "update image embedding for all assets.",
				Action: func(ctx context.Context, c *cli.Command) error {
					return populateImageEmbeddings(ctx)
				},
			},
		},
	}

	ctx := context.Background()
	if err := cmd.Run(ctx, os.Args); err != nil {
		slog.Error("operation failed", slog.String("error", err.Error()))
	}
}

func performWork(ctx context.Context) error {
	err := db.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return fmt.Errorf("unable to connect to the database: %w", err)
	}
	defer db.Close(ctx)

	err = queue.Init(ctx)
	if err != nil {
		return fmt.Errorf("unable to start job queue: %w", err)
	}
	defer queue.Shutdown(ctx)

	err = processExistingItems(ctx)
	if err != nil {
		return fmt.Errorf("unable to processing pending items :%w", err)
	}

	address := os.Getenv("WORKER_ADDRESS")
	if address == "" {
		return fmt.Errorf("invalid worker address")
	}
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("unable to start server: %w", err)
	}
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterWorkerServiceServer(grpcServer, &service.WorkerServiceServer{})

	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("error running grpc server: %w", err)
	}
	return nil
}

func processExistingItems(ctx context.Context) error {
	quries, _ := db.Get()

	assets, err := quries.GetPendingAssets(ctx)
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

		queue.EnqueueAssetProcessing(ctx, asset.ID.String())
	}

	return nil
}

func populateImageEmbeddings(ctx context.Context) error {
	err := db.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return fmt.Errorf("unable to connect to the database: %w", err)
	}
	defer db.Close(ctx)

	queries, _ := db.Get()
	assets, err := queries.GetImageAssetsWithoutEmbedding(ctx)
	if err != nil {
		return fmt.Errorf("unable to retrieve assets with embedding missing :%w", err)
	}

	for _, asset := range assets {
		err = queue.PopulateImageEmbedding(ctx, &asset, nil)
		if err != nil {
			slog.Error("unable to populate embedding", slog.String("error", err.Error()))
			continue
		}

		_, err = queries.UpdateAsset(ctx, db.UpdateAssetParams{
			ID:              asset.ID,
			Filename:        asset.Filename,
			Checksum:        asset.Checksum,
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
	}

	return nil
}

func processSingle(ctx context.Context, id string) error {
	err := db.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return fmt.Errorf("unable to connect to the database: %w", err)
	}
	defer db.Close(ctx)

	return queue.ProcessAsset(ctx, id)
}
