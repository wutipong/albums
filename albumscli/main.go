package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/lmittmann/tint"
	"github.com/urfave/cli/v3"
	"github.com/wutipong/albums/albumscli/album"
	"github.com/wutipong/albums/albumscli/importing"
	"github.com/wutipong/albums/albumscli/process"
	"github.com/wutipong/albums/albumscli/profile"
)

func main() {
	slog.SetDefault(slog.New(tint.NewHandler(os.Stderr, &tint.Options{
		Level:      slog.LevelInfo,
		TimeFormat: time.Kitchen,
	})))

	displayLogLevelStr := "warn"
	fileLogLevelStr := "info"
	profileStr := "default"

	cmd := &cli.Command{
		Name:  "albumscli",
		Usage: "import assets to albums",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "display-log",
				Value:       "warn",
				Usage:       "Minimum log-level on display (debug, info, warn, error).",
				Destination: &displayLogLevelStr,
				Category:    "Logging",
			},
			&cli.StringFlag{
				Name:        "file-log",
				Value:       "info",
				Usage:       "Minimum log-level in log file (debug, info, warn, error).",
				Destination: &fileLogLevelStr,
				Category:    "Logging",
			},
			&cli.StringFlag{
				Name:        "profile",
				Value:       "default",
				Usage:       "profile of albums server.",
				Destination: &profileStr,
				Category:    "albums Server",
			},
		},
		Commands: []*cli.Command{
			profile.Command(&profileStr, &displayLogLevelStr, &fileLogLevelStr),
			album.Command(&profileStr, &displayLogLevelStr, &fileLogLevelStr),
			process.Command(&profileStr, &displayLogLevelStr, &fileLogLevelStr),
			importing.Command(&profileStr, &displayLogLevelStr, &fileLogLevelStr),
		},
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if err := cmd.Run(ctx, os.Args); err != nil {
		log.Fatal(err)
	}
}
