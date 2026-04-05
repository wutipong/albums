package profile

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	"github.com/wutipong/albums/albums-importer/log"
)

func Command(profile *string, displayLogLevel *string, fileLogLevel *string) *cli.Command {
	return &cli.Command{
		Name:  "profile",
		Usage: "Manage profiles for albums importer",
		Commands: []*cli.Command{
			{
				Name: "setup",
				Usage: "Setup profile file interactively. " +
					"Existing configuration file will be overwritten.",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					err := log.Setup(*profile, *displayLogLevel, false, *fileLogLevel)
					if err != nil {
						return fmt.Errorf("unable to setup log: %w", err)
					}
					defer log.CleanUp()

					return setupProfile(*profile)
				},
			}, {
				Name:  "show",
				Usage: "Show profile configuration",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return showProfile(*profile)
				},
			},
		},
	}
}
