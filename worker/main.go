package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/urfave/cli/v3"
	"github.com/wutipong/albums/worker/db"
	"github.com/wutipong/albums/worker/queue"
	"github.com/wutipong/albums/worker/service"
	"github.com/wutipong/albums/worker/service/pb"
	"google.golang.org/grpc"
)

//go:generate protoc --go_out=. --go-grpc_out=. -I/workspaces/grpc worker.proto clip.proto
//go:generate sqlc generate
//go:generate vipsgen -out ./vips

func main() {
	id := ""
	processPending := false
	debug := false
	cmd := &cli.Command{
		Name:  "worker",
		Usage: "process assets to albums",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return performWork(ctx, processPending)
		},

		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "process-pending",
				Usage:       "process pending items in the queue.",
				Destination: &processPending,
			},
			&cli.BoolFlag{
				Name:        "debug",
				Usage:       "enable debug logging.",
				Destination: &debug,
			},
		},
		Before: func(ctx context.Context, c *cli.Command) (cctx context.Context, err error) {
			level := slog.LevelInfo
			if debug {
				level = slog.LevelDebug
			}
			slog.SetDefault(slog.New(
				slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
					Level: level,
				}),
			))

			return
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
					return queue.EnqueuePopulateAlbumsCover(ctx, false)
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

func performWork(ctx context.Context, processPending bool) error {
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

	if processPending {
		err = queue.EnqueueProcessAllAssets(ctx, true)
		if err != nil {
			return fmt.Errorf("unable to processing pending items :%w", err)
		}
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
		e := queue.EnqueuePopulateImageEmbedding(ctx, asset.ID.String())
		if e != nil {
			return fmt.Errorf("unable to add new populate image embedding job: %w", err)
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

	endpoint, secure, err := queue.GetMinioEndpoint(os.Getenv("AWS_ENDPOINT_URL"))
	if err != nil {
		return fmt.Errorf("unable to get endpoint: %w", err)
	}

	accessKeyId := os.Getenv("AWS_ACCESS_KEY_ID")
	secret := os.Getenv("AWS_SECRET_ACCESS_KEY")

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:        credentials.NewStaticV4(accessKeyId, secret, ""),
		Secure:       secure,
		BucketLookup: minio.BucketLookupPath,
	})
	if err != nil {
		return fmt.Errorf("unable to create minio client: %w", err)
	}

	return queue.ProcessAsset(ctx, minioClient, id)
}
