package service

import (
	"context"

	"github.com/wutipong/albums/worker/service/definition"
)

type WorkerServiceServer struct {
	definition.UnimplementedWorkerServiceServer
}

// Notify worker to process specific asset.
func (s *WorkerServiceServer) NotifyProcessAsset(
	ctx context.Context,
	req *definition.NotifyProcessAssetResquest,
) (resp *definition.NotifyProcessAssetResponse, err error) {

	return
}

// Notify worker to queue unprocessed asset to processing queue.
func (s *WorkerServiceServer) NotifyScanCache(
	context.Context,
	*definition.NotifyScanCacheRequest,
) (resp *definition.NotifyScanCacheResponse, err error) {
	return
}
