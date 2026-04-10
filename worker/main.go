package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/lmittmann/tint"
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

	ctx := context.Background()

	err := db.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		slog.Error("unable to conect to database", slog.String("error", err.Error()))
		return
	}
	defer db.Close(ctx)

	err = queue.Init(ctx)
	if err != nil {
		slog.Error("unable to start job queue", slog.String("error", err.Error()))
		return
	}
	defer queue.Shutdown(ctx)

	err = processExistingItems(ctx)
	if err != nil {
		slog.Error("unable to processing pending items", slog.String("error", err.Error()))
	}

	address := os.Getenv("WORKER_ADDRESS")
	if address == "" {
		slog.Error("invalid address")
		return
	}
	lis, err := net.Listen("tcp", address)
	if err != nil {
		slog.Error("unable to start server", slog.String("error", err.Error()))
		return
	}
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterWorkerServiceServer(grpcServer, &service.WorkerServiceServer{})

	if err := grpcServer.Serve(lis); err != nil {
		slog.Error("error running grpc server.", slog.String("error", err.Error()))
	}

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
