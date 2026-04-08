package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/wutipong/albums/worker/db"
	"github.com/wutipong/albums/worker/queue"
	"github.com/wutipong/albums/worker/service/pb"
)

type WorkerServiceServer struct {
	pb.UnimplementedWorkerServiceServer
}

// Notify worker to process specific asset.
func (s *WorkerServiceServer) NotifyProcessAsset(
	ctx context.Context,
	req *pb.NotifyProcessAssetResquest,
) (resp *pb.NotifyProcessAssetResponse, err error) {
	id := req.Id
	resp = &pb.NotifyProcessAssetResponse{
		Id: id,
	}

	uuid := pgtype.UUID{}
	err = uuid.Scan(id)
	if err != nil {
		err = fmt.Errorf("unable to parse asset id: %w", err)
		return
	}

	quries := db.New(db.Connection())
	processStatus, err := quries.GetAssetProcessStatus(ctx, uuid)
	if err != nil {
		err = fmt.Errorf("unable to find asset: %w", err)
		return
	}

	if processStatus != db.ProcessStatusTPending {
		slog.Info("adding asset", slog.String("id", id))

		queue.EnqueueAssetProcessing(ctx, id)
	}

	switch processStatus {
	case db.ProcessStatusTPending:
		resp.Status = pb.AssetStatus_PENDING
	case db.ProcessStatusTProcessing:
		resp.Status = pb.AssetStatus_PROCESSING
	case db.ProcessStatusTProcessed:
		resp.Status = pb.AssetStatus_PROCESSED
	}

	return
}

// Notify worker to queue unprocessed asset to processing queue.
func (s *WorkerServiceServer) NotifyScanCache(
	ctx context.Context,
	req *pb.NotifyScanCacheRequest,
) (resp *pb.NotifyScanCacheResponse, err error) {
	quries := db.New(db.Connection())
	assets, err := quries.GetPendingAssets(ctx)

	resp = &pb.NotifyScanCacheResponse{}
	if len(assets) == 0 {
		return
	}

	for _, asset := range assets {
		slog.Info("adding asset", slog.String("id", asset.ID.String()))

		queue.EnqueueAssetProcessing(ctx, asset.ID.String())
	}

	return
}
