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

	slog.Info("enqueue asset to asset-processing queue", slog.String("id", id))

	uuid := pgtype.UUID{}
	err = uuid.Scan(id)
	if err != nil {
		err = fmt.Errorf("unable to parse asset id: %w", err)
		return
	}

	quries, _ := db.Get()

	processStatus, err := quries.GetAssetProcessStatus(ctx, uuid)
	if err != nil {
		err = fmt.Errorf("unable to find asset: %w", err)
		return
	}

	slog.Info("adding asset", slog.String("id", id))

	_, err = queue.EnqueueAssetProcessing(ctx, id)
	if err != nil {
		err = fmt.Errorf("unable to enqueue a new job: %w", err)
		return
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
	quries, _ := db.Get()

	assets, err := quries.GetPendingAssets(ctx)

	slog.Info("scan library for unprocessed asset.")

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

func (s *WorkerServiceServer) UpdateAlbumThumbnail(
	ctx context.Context,
	req *pb.UpdateAlbumThumbnailRequest,
) (resp *pb.UpdateAlbumThumbnailResponse, err error) {

	resp = &pb.UpdateAlbumThumbnailResponse{
		Id:      req.Id,
		AssetId: req.AssetId,
	}

	err = queue.EnqueuePopulateAlbumsCover(ctx, req.Id, req.AssetId)
	if err != nil {
		err = fmt.Errorf("unable to unque popluate album cover.")
		return
	}

	return
}
