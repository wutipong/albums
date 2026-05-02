package importing

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/wutipong/albums/albumscli/server/api"
	"github.com/wutipong/albums/albumscli/server/types"
)

func ProcessDirectory(
	ctx context.Context,
	server api.ServerConfig,
	album types.Album,
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

	filepath.WalkDir(filepath.Join(sourceDir, path), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			slog.Error(
				"failed to access path",
				slog.String("path", path),
				slog.String("error", err.Error()),
			)
		}
		if d.IsDir() {
			return nil
		}

		if IsMediaFile(path) {
			err = processMediaFile(ctx, server, sourceDir, path, album)
		} else if IsArchiveFile(path) {
			err = processArchive(ctx, server, sourceDir, path, album)
		}

		if err != nil {
			slog.Error(
				"failed to process file",
				slog.String("path", path),
				slog.String("error", err.Error()),
			)
			return nil
		}
		return nil
	})

	return nil
}

func processMediaFile(
	ctx context.Context,
	server api.ServerConfig,
	sourceDir string,
	path string,
	album types.Album,
) error {
	slog.Debug("processing media file",
		slog.String("sourceDir", sourceDir),
		slog.String("path", path),
	)

	if !IsMediaFile(path) {
		slog.Debug("skipping non-media file",
			slog.String("path", path),
		)
		return nil
	}

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", path, err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat file %s: %w", path, err)
	}

	asset, err := api.PostAsset(
		ctx,
		server,
		album.ID,
		path,
		file.Name(),
		file,
		info.Size(),
	)
	if err != nil {
		if errors.Is(err, ErrDuplicateAsset) {
			slog.Warn(
				"asset already exists. skipping file.",
				slog.String("path", path),
			)
			return nil
		}
		return fmt.Errorf("failed to upload asset for file %s: %w", path, err)
	}

	slog.Info("uploaded asset", slog.Any("asset", asset))

	return nil
}

func processArchive(
	ctx context.Context,
	server api.ServerConfig,
	sourceDir string,
	path string,
	album types.Album,
) error {
	slog.Debug("processing archive file",
		slog.String("sourceDir", sourceDir),
		slog.String("path", path),
	)

	archiveFile, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open archive: %s: %w.",
			path,
			err,
		)
	}
	defer archiveFile.Close()

	err = WalkArchive(ctx, server, album.ID, path, archiveFile)
	if err != nil {
		return fmt.Errorf("failed to process archive %s: %w", path, err)
	}

	return nil
}
