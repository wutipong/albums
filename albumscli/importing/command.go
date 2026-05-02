package importing

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"slices"

	"github.com/urfave/cli/v3"
	"github.com/wutipong/albums/albumscli/log"
	"github.com/wutipong/albums/albumscli/profile"
	"github.com/wutipong/albums/albumscli/server/api"
	"github.com/wutipong/albums/albumscli/server/types"
)

func Command(profileStr *string, displayLogLevel *string, fileLogLevel *string) *cli.Command {
	sourceDir := ""
	force := false
	dryRun := false

	return &cli.Command{
		Name:  "import",
		Usage: "perform importing assets.",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "force",
				Value:       false,
				Usage:       "Force processing album even if an album with the same name exists.",
				Destination: &force,
			},
			&cli.BoolFlag{
				Name:        "dry-run",
				Value:       false,
				Usage:       "Processing assets without working with the Albums server.",
				Destination: &dryRun,
				Category:    "Processing",
			},
		},
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name:        "source-dir",
				UsageText:   "Source directory path",
				Destination: &sourceDir,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			err := log.Setup(*profileStr, *displayLogLevel, true, *fileLogLevel)
			if err != nil {
				return fmt.Errorf("unable to setup log: %w", err)
			}
			defer log.CleanUp()

			c, err := profile.LoadProfile(ctx, *profileStr)
			if err != nil {
				return fmt.Errorf(
					"unable to load configuration. please run 'albumscli setup' first: %w",
					err,
				)
			}

			slog.Info("Albums instance",
				slog.String("url", c.URL),
			)

			url, err := url.Parse(c.URL)
			if err != nil {
				return fmt.Errorf("invalid Albums url: %w", err)
			}

			server := api.ServerConfig{
				URL:     url,
				DryRun:  dryRun,
				APIKey:  c.APIKey,
				Network: string(c.Network),
			}

			return Process(
				ctx,
				server,
				sourceDir,
				force,
			)
		},
	}
}

func Process(
	ctx context.Context,
	server api.ServerConfig,
	sourceDir string,
	force bool,
) error {
	var albums []types.Album
	var err error

	if !force {
		resp, err := api.GetAlbumList(ctx, server)
		if err != nil {
			err = fmt.Errorf("unable to retrieved existing albums: %w", err)
			return err
		}
		albums = resp.Albums
		slog.Debug("albums", slog.Any("existing albums", albums))
	}

	entries, err := os.ReadDir(sourceDir)
	if err != nil {
		return fmt.Errorf("failed to read source directory: %w", err)
	}

	for _, entry := range entries {
		path := filepath.Join(sourceDir, entry.Name())

		albumPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			slog.Error(
				"failed to determine album name.",
				slog.String("error", err.Error()),
			)
			continue
		}

		slog.Debug("processing path",
			slog.String("path", path),
			slog.String("albumPath", albumPath),
		)

		matchingAlbums := slices.DeleteFunc(slices.Clone(albums),
			func(album types.Album) bool {
				return album.Name != albumPath
			})

		slog.Debug("matching albums",
			slog.Any("album", albumPath),
			slog.Any("existing", albums),
			slog.Any("matchalbums", matchingAlbums),
		)

		if !force && len(matchingAlbums) > 0 {
			slog.Warn(
				"album already exists. skipping.",
				slog.String("name", albumPath),
				slog.String("id", matchingAlbums[0].ID),
			)
			continue
		}

		slog.Info("creating album",
			slog.String("name", albumPath),
			slog.String("entry", entry.Name()),
		)

		if entry.IsDir() {
			err = ProcessDirectory(ctx, server, sourceDir, albumPath)
		} else {
			if !IsArchiveFile(path) {
				slog.Debug("skipping unsupported file",
					slog.String("path", path),
				)
				continue
			}
			err = ProcessArchive(ctx, server, sourceDir, albumPath)
		}

		if err != nil {
			slog.Error(
				"failed upload assets.",
				slog.String("error", err.Error()),
				slog.String("sourceDir", sourceDir),
				slog.String("albumPath", albumPath),
			)
			if errors.Is(err, context.Canceled) {
				return err
			}
			continue
		}
	}

	return err
}
