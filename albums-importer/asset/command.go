package asset

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
	path := ""

	return &cli.Command{
		Name:  "asset",
		Usage: "Manage assets in the server",
		Commands: []*cli.Command{
			{
				Name:    "create",
				Aliases: []string{"new", "n"},
				Usage:   "Create a new asset",
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
						Name:        "path",
						UsageText:   "Asset path",
						Destination: &path,
					}, &cli.StringArg{
						Name:        "album-id",
						UsageText:   "Album ID to add the asset to",
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

					return createAsset(ctx, *profile, dryRun, path, id)
				},
			},
		},
	}
}
