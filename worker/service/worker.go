package service

import (
	"context"

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

	return
}

// Notify worker to queue unprocessed asset to processing queue.
func (s *WorkerServiceServer) NotifyScanCache(
	context.Context,
	*pb.NotifyScanCacheRequest,
) (resp *pb.NotifyScanCacheResponse, err error) {
	return
}
