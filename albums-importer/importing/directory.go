package importing

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"slices"

	"github.com/wutipong/albums/albums-importer/server/api"
)

func ProcessDirectory(
	ctx context.Context,
	server api.ServerConfig,
	sourceDir string,
	path string,
) error {
	if ctx.Err() != nil {
		return fmt.Errorf("context error: %w", ctx.Err())
	}
	slog.Debug("processing directory",
		slog.String("sourceDir", sourceDir),
		slog.String("path", path),
	)
	entries, err := os.ReadDir(filepath.Join(sourceDir, path))
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	slog.Debug("directory entries", slog.Any("size", len(entries)))

	files := slices.DeleteFunc(
		entries,
		func(d os.DirEntry) bool {
			if d.IsDir() {
				return true
			}
			if !IsMediaFile(d.Name()) {
				return true
			}
			return false
		},
	)

	slog.Debug("media files", slog.Any("size", len(files)))
	if len(files) == 0 {
		slog.Info("no media files found in directory. skipping.",
			slog.String("path", path),
		)
		return nil
	}

	album, err := api.CreateAlbum(ctx, server, path)
	if err != nil {
		return fmt.Errorf("failed to create album for directory %s: %w", path, err)
	}

	slog.Info("album created for directory",
		slog.String("album", album.Name),
		slog.String("id", album.ID),
	)

	for _, file := range files {
		slog.Info(
			"creating asset",
			slog.String("path", path),
			slog.String("entry", file.Name()),
		)

		info, e := file.Info()
		if e != nil {
			return fmt.Errorf(
				"Unable to read image file propery: %s: %w",
				file.Name(),
				e,
			)
		}

		reader, e := os.Open(filepath.Join(sourceDir, path, file.Name()))
		if e != nil {
			return fmt.Errorf(
				"failed to open image file %s: %w",
				file.Name(),
				e,
			)
		}

		defer reader.Close()

		asset, err := api.PostAsset(
			ctx,
			server,
			album.ID,
			path,
			file.Name(),
			reader,
			info.ModTime(),
		)
		if err != nil {
			slog.Error("failed to upload asset",
				slog.String("album", path),
				slog.String("entry", file.Name()),
				slog.String("error", err.Error()),
			)
			continue
		}

		slog.Info("uploaded asset", slog.Any("asset", asset))
	}
	return nil
}
