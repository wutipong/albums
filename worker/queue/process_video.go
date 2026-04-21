package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jackc/pgx/v5/pgtype"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"github.com/wutipong/albums/worker/db"
)

const VIDEO_WIDTH = 1280
const VIDEO_HEIGHT = 720

func processVideoAsset(ctx context.Context, s3Client *s3.Client, asset *db.Asset) error {
	slog.Info("process video asset", slog.Any("id", asset.ID))

	err := ctx.Err()
	if err != nil {
		slog.Info("context.", slog.String("error", err.Error()))
		return fmt.Errorf("context cancelled: %w", err)
	}

	s3Obj, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Key:    aws.String(asset.Original),
		Bucket: aws.String(os.Getenv("S3_BUCKET")),
	})

	if err != nil {
		return fmt.Errorf("unable to get object from s3: %w", err)
	}
	defer s3Obj.Body.Close()

	originalFile, err := os.CreateTemp("",
		fmt.Sprintf("*.%s", filepath.Base(asset.Filename)),
	)

	if err != nil {
		return fmt.Errorf("unable to create temp file for original asset: %w", err)
	}
	defer os.Remove(originalFile.Name())

	io.Copy(originalFile, s3Obj.Body)

	probe, err := ffmpeg.Probe(originalFile.Name())
	if err != nil {
		return fmt.Errorf("unable to probe original video: %w", err)
	}

	var info Probe
	json.Unmarshal([]byte(probe), &info)

	err = processVideoThumbnail(ctx, s3Client, asset, originalFile, info)
	if err != nil {
		return fmt.Errorf("unable to process video asset thumbnail: %w", err)
	}

	err = processVideoPreview(ctx, s3Client, asset, originalFile, info)
	if err != nil {
		return fmt.Errorf("unable to process video asset preview: %w", err)
	}

	err = processVideoView(ctx, s3Client, asset, originalFile, info)
	if err != nil {
		return fmt.Errorf("unable to process video asset view: %w", err)
	}

	asset.Type = db.AssetTypeTVideo

	return nil
}

func processVideoView(
	ctx context.Context, s3Client *s3.Client, asset *db.Asset,
	originalFile *os.File, _ Probe,
) error {
	slog.Info("process video asset view media", slog.Any("id", asset.ID))
	err := ctx.Err()
	if err != nil {
		return fmt.Errorf("context cancelled: %w", err)
	}

	asset.View = createAssetKey()
	outputFile, err := os.CreateTemp("", "*view.mp4")
	if err != nil {
		return fmt.Errorf("unable to create temp file to transcode: %w", err)
	}
	defer os.Remove(outputFile.Name())

	err = ffmpeg.Input(originalFile.Name()).
		Output(outputFile.Name(), ffmpeg.KwArgs{
			"vf": fmt.Sprintf(
				"scale=%d:%d:force_original_aspect_ratio=decrease,scale=trunc(iw/2)*2:trunc(ih/2)*2",
				VIDEO_WIDTH, VIDEO_HEIGHT,
			),
			"c:v":      "libx264",
			"preset":   "superfast",
			"crf":      "30",
			"pix_fmt":  "yuv420p",    // Fixes browser incompatibility [6]
			"c:a":      "aac",        // Standard audio
			"movflags": "+faststart", // Enables progressive loading [5]
		}).OverWriteOutput().ErrorToStdOut().Run()

	if err != nil {
		return fmt.Errorf("unable to create view asset for video asset: %w", err)
	}

	probe, err := ffmpeg.Probe(outputFile.Name())
	if err != nil {
		return fmt.Errorf("unable to probe original video: %w", err)
	}

	var viewInfo Probe
	json.Unmarshal([]byte(probe), &viewInfo)

	viewVideoStream, err := viewInfo.Video()
	if err != nil {
		return fmt.Errorf("unable to get video stream from video asset: %w", err)
	}
	asset.ViewWidth = int32(viewVideoStream.Width)
	asset.ViewHeight = int32(viewVideoStream.Height)

	outputFile.Seek(0, io.SeekStart)

	_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(os.Getenv("S3_BUCKET")),
		Body:        outputFile,
		Key:         aws.String(asset.View),
		ContentType: aws.String("video/mp4"),
	})

	return nil
}

func processVideoThumbnail(
	ctx context.Context, s3Client *s3.Client, asset *db.Asset, originalFile *os.File, info Probe,
) error {
	slog.Info("process video asset thumbnail", slog.Any("id", asset.ID))
	err := ctx.Err()
	if err != nil {
		return fmt.Errorf("context cancelled: %w", err)
	}

	duration, err := strconv.ParseFloat(info.Format.Duration, 10)
	if err != nil {
		return fmt.Errorf("unable to parse duration: %w", err)
	}

	outputFile, err := os.CreateTemp("", "*view.webp")
	if err != nil {
		return fmt.Errorf("unable to create temp file to transcode: %w", err)
	}
	defer os.Remove(outputFile.Name())

	// save thumbnail at 1/3 duration
	err = ffmpeg.
		Input(originalFile.Name(), ffmpeg.KwArgs{
			"ss": fmt.Sprintf("%f", duration/3),
		}).
		Output(outputFile.Name(), ffmpeg.KwArgs{
			"c:v":     "libwebp",
			"vframes": "1",
			"quality": fmt.Sprintf("%d", THUMBNAIL_QUALITY),
		}).OverWriteOutput().ErrorToStdOut().Run()

	if err != nil {
		return fmt.Errorf("unable to create thumbnail asset for video asset: %w", err)
	}

	videoDuration := time.Duration(duration) * time.Second
	asset.VideoDuration = pgtype.Interval{
		Microseconds: videoDuration.Microseconds(),
		Valid:        true,
	}

	_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(os.Getenv("S3_BUCKET")),
		Body:        outputFile,
		Key:         aws.String(asset.Thumbnail),
		ContentType: aws.String("image/webp"),
	})
	return nil
}

func processVideoPreview(
	ctx context.Context, s3Client *s3.Client, asset *db.Asset,
	originalFile *os.File, info Probe,
) error {
	slog.Info("process video preview", slog.Any("id", asset.ID))
	err := ctx.Err()
	if err != nil {
		return fmt.Errorf("context cancelled: %w", err)
	}
	outputFile, err := os.CreateTemp("", "*view.webp")
	if err != nil {
		return fmt.Errorf("unable to create temp file to transcode: %w", err)
	}
	defer os.Remove(outputFile.Name())

	duration, err := strconv.ParseFloat(info.Format.Duration, 10)
	if err != nil {
		return fmt.Errorf("unable to parse duration: %w", err)
	}

	// save preview at 1/3 duration, 5 seconds-long in 5 fps.
	err = ffmpeg.
		Input(originalFile.Name(), ffmpeg.KwArgs{
			"ss": fmt.Sprintf("%f", duration/3),
		}).
		Output(outputFile.Name(), ffmpeg.KwArgs{
			"c:v":     "libwebp",
			"t":       "5",
			"loop":    "0",
			"quality": fmt.Sprintf("%d", THUMBNAIL_QUALITY),
			"vf":      "fps=5",
		}).OverWriteOutput().ErrorToStdOut().Run()

	if err != nil {
		return fmt.Errorf("unable to create thumbnail asset for video asset: %w", err)
	}

	asset.Preview = createAssetKey()
	_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(os.Getenv("S3_BUCKET")),
		Body:        outputFile,
		Key:         aws.String(asset.Preview),
		ContentType: aws.String("image/webp"),
	})

	return nil
}
