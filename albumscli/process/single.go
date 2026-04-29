package process

import (
	"context"
	"fmt"
	"path"

	"github.com/wutipong/albums/albumscli/server/api"
)

func processSingle(ctx context.Context, server api.ServerConfig, dryRun bool, id string) error {
	if dryRun {
		return nil
	}
	_, err := api.Get[any](ctx, server, path.Join("api", "asset", id, "process"))
	if err != nil {
		return fmt.Errorf("unable to queue process asset command: %w", err)
	}

	return nil
}
