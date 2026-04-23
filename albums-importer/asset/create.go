package asset

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"os"

	"github.com/wutipong/albums/albums-importer/profile"
	"github.com/wutipong/albums/albums-importer/server/api"
)

func createAsset(ctx context.Context, profileName string, dryRun bool, path string, albumID string) error {
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

	slog.Info("creating asset",
		slog.String("path", path),
		slog.String("album_id", albumID),
	)

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("Unable to open file: %w", err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("unable to retrieve filestat: %w", err)
	}

	resp, err := api.PostAsset(ctx, server, albumID, "", path, file, stat.Size())
	if err != nil {
		return err
	}

	slog.Info("asset created",
		slog.String("id", resp.Asset.ID),
		slog.String("filename", resp.Asset.Filename),
	)

	return nil
}
