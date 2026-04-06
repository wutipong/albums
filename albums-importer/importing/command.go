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
	"github.com/wutipong/albums/albums-importer/log"
	"github.com/wutipong/albums/albums-importer/profile"
	"github.com/wutipong/albums/albums-importer/server/api"
	"github.com/wutipong/albums/albums-importer/server/types"
)

func Command(profileStr *string, displayLogLevel *string, fileLogLevel *string) *cli.Command {
	sourceDir := ""
	force := false
	dryRun := false
	disableDirectory := true
	disableArchive := true

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
				Usage:       "Processing assets without working with the Immich server.",
				Destination: &dryRun,
				Category:    "Processing",
			},
			&cli.BoolFlag{
				Name:        "disable-directory",
				Value:       false,
				Usage:       "Disable processing media files in directories.",
				Destination: &disableDirectory,
				Category:    "Processing",
			},
			&cli.BoolFlag{
				Name:        "disable-archive",
				Value:       false,
				Usage:       "Disable processing media files in archive files.",
				Destination: &disableArchive,
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
					"unable to load configuration. please run 'immich-importer setup' first: %w",
					err,
				)
			}

			slog.Info("Immich instance",
				slog.String("url", c.URL),
			)

			url, err := url.Parse(c.URL)
			if err != nil {
				return fmt.Errorf("invalid immich url: %w", err)
			}

			server := api.ServerConfig{
				URL:    url,
				DryRun: dryRun,
			}

			return Process(
				ctx,
				server,
				sourceDir,
				force,
				!disableDirectory,
				!disableArchive,
			)
		},
	}
}

func Process(
	ctx context.Context,
	server api.ServerConfig,
	sourceDir string,
	force bool,
	processDirectory bool,
	processArchive bool,

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
	err = filepath.WalkDir(sourceDir,
		func(path string, d os.DirEntry, err error,
		) error {
			if err != nil {
				slog.Warn(
					"failed to access path. skipping.",
					slog.String("path", path),
					slog.String("error", err.Error()),
				)
				return nil
			}

			var assetIds []string

			albumPath, err := filepath.Rel(sourceDir, path)
			if err != nil {
				slog.Error(
					"failed to determine album name.",
					slog.String("error", err.Error()),
				)
				return nil
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
				)
				return nil
			}

			if d.IsDir() {
				if !processDirectory {
					slog.Debug("skipping directory",
						slog.String("path", path),
					)
					return nil
				}
				err = ProcessDirectory(ctx, server, sourceDir, albumPath)
			} else {
				if !processArchive {
					slog.Debug("skipping file",
						slog.String("path", path),
					)
					return nil
				}
				if !IsArchiveFile(path) {
					slog.Debug("skipping unsupported file",
						slog.String("path", path),
					)
					return nil
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
				return nil
			}

			if len(assetIds) == 0 {
				slog.Debug(
					"no assets uploaded. skip create album.",
					slog.String("name", albumPath),
				)
				return nil
			}

			if len(matchingAlbums) > 0 {
				slog.Info(
					"album already exists. update existing album.",
					slog.String("name", albumPath),
				)

				var albumIds []string
				for _, album := range matchingAlbums {
					albumIds = append(albumIds, album.ID)
				}

				return nil
			}

			return nil
		})

	return err
}
