package album

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/urfave/cli/v3"
)

func Command(profile *string, displayLogLevel *string, fileLogLevel *string) *cli.Command {
	dryRun := false
	id := ""
	name := ""

	return &cli.Command{
		Name:  "album",
		Usage: "Manage albums in the server",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List all albums in the server. ",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:        "dry-run",
						Value:       false,
						Usage:       "Don't make actual API calls to the server. Useful for testing.",
						Destination: &dryRun,
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return listAlbum(ctx, *profile, dryRun)
				},
			}, {
				Name:  "show",
				Usage: "Show album details",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:        "dry-run",
						Value:       false,
						Usage:       "Don't make actual API calls to the server. Useful for testing.",
						Destination: &dryRun,
					},
				},
				Arguments: []cli.Argument{
					&cli.StringArg{
						Name:        "id",
						UsageText:   "Album ID",
						Destination: &id,
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if displayLogLevel != nil && *displayLogLevel == "debug" {
						slog.SetDefault(slog.New(tint.NewHandler(os.Stderr, &tint.Options{
							Level:      slog.LevelDebug,
							TimeFormat: time.Kitchen,
						})))
					}

					return showAlbum(ctx, *profile, dryRun, id)
				},
			}, {
				Name:    "create",
				Aliases: []string{"new", "n"},
				Usage:   "Create a new album",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:        "dry-run",
						Value:       false,
						Usage:       "Don't make actual API calls to the server. Useful for testing.",
						Destination: &dryRun,
					},
				},
				Arguments: []cli.Argument{
					&cli.StringArg{
						Name:        "name",
						UsageText:   "Album name",
						Destination: &name,
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if displayLogLevel != nil && *displayLogLevel == "debug" {
						slog.SetDefault(slog.New(tint.NewHandler(os.Stderr, &tint.Options{
							Level:      slog.LevelDebug,
							TimeFormat: time.Kitchen,
						})))
					}

					return createAlbum(ctx, *profile, dryRun, name)
				},
			},
		},
	}
}
