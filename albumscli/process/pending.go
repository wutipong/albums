package process

import (
	"context"
	"fmt"
	"path"

	"github.com/wutipong/albums/albumscli/server/api"
)

func processPending(ctx context.Context, server api.ServerConfig, dryRun bool) error {
	if dryRun {
		return nil
	}
	_, err := api.Get[any](ctx, server, path.Join("api", "asset", "scan"))
	if err != nil {
		return fmt.Errorf("unable to queue process pending asset command: %w", err)
	}

	return nil
}
