package album

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/wutipong/albums/albums-importer/profile"
	"github.com/wutipong/albums/albums-importer/server/api"
)

func createAlbum(ctx context.Context, profileName string, dryRun bool, name string) (err error) {
	config, err := profile.LoadProfile(ctx, profileName)
	if err != nil {
		return err
	}

	slog.Info("profile",
		slog.String("name", profileName),
		slog.String("url", config.URL),
	)

	serverUrl, err := url.Parse(config.URL)
	if err != nil {
		return fmt.Errorf("Unable to parse server URL: %w", err)
	}
	server := api.ServerConfig{
		URL:     serverUrl,
		DryRun:  dryRun,
		APIKey:  config.APIKey,
		Network: string(config.Network),
	}

	resp, err := api.CreateAlbum(ctx, server, name)
	if err != nil {
		return err
	}

	slog.Info("album created",
		slog.String("id", resp.ID),
		slog.String("name", resp.Name),
	)

	return nil
}
