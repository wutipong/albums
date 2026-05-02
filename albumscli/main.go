package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/lmittmann/tint"
	"github.com/urfave/cli/v3"
	"github.com/wutipong/albums/albumscli/album"
	"github.com/wutipong/albums/albumscli/importing"
	"github.com/wutipong/albums/albumscli/log"
	"github.com/wutipong/albums/albumscli/process"
	"github.com/wutipong/albums/albumscli/profile"
)

func main() {
	slog.SetDefault(slog.New(tint.NewHandler(os.Stderr, &tint.Options{
		Level:      slog.LevelInfo,
		TimeFormat: time.Kitchen,
	})))

	profileStr := "default"
	debug := false

	cmd := &cli.Command{
		Name:  "albumscli",
		Usage: "import assets to albums",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "profile",
				Value:       "default",
				Usage:       "profile of albums server.",
				Destination: &profileStr,
				Category:    "albums Server",
			},
			&cli.BoolFlag{
				Name:        "debug",
				Value:       false,
				Usage:       "enable loggin debug message",
				Destination: &debug,
			},
		},
		Commands: []*cli.Command{
			profile.Command(&profileStr),
			album.Command(&profileStr),
			process.Command(&profileStr),
			importing.Command(&profileStr),
		},
		Before: func(ctx context.Context, c *cli.Command) (context.Context, error) {
			level := slog.LevelInfo.String()
			if debug {
				level = slog.LevelDebug.String()
			}
			log.Setup(profileStr, level, true, level)
			return ctx, nil

		},
		After: func(ctx context.Context, c *cli.Command) error {
			log.CleanUp()
			return nil
		},
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if err := cmd.Run(ctx, os.Args); err != nil {
		slog.Error("Error running command", slog.String("error", err.Error()))
	}
}
