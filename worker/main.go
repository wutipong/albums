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

//go:generate protoc --go_out=. --go-grpc_out=. -I/workspaces/grpc worker.proto
//go:generate sqlc generate

func main() {
	slog.SetDefault(slog.New(tint.NewHandler(os.Stderr, &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: time.Kitchen,
	})))

	cmd := &cli.Command{
		Name:  "albums-importer",
		Usage: "import assets to albums",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return performWork(ctx)
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
