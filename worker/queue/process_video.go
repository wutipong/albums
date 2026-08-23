package queue

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/minio/minio-go/v7"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"github.com/wutipong/albums/worker/db"
	vips "github.com/wutipong/albums/worker/vips"
)

const VIDEO_WIDTH = 1280
const VIDEO_HEIGHT = 720

func processVideoAsset(ctx context.Context, minioClient *minio.Client, asset *db.Asset) error {
	slog.Info("process video asset", slog.Any("id", asset.ID))

	err := ctx.Err()
	if err != nil {
		slog.Info("context.", slog.String("error", err.Error()))
		return fmt.Errorf("context cancelled: %w", err)
	}

	s3Obj, err := getObjectFromS3(ctx, minioClient, asset.Original)

	if err != nil {
		return fmt.Errorf("unable to get object from s3: %w", err)
	}
	defer s3Obj.Close()

	originalFile, err := os.CreateTemp("",
		fmt.Sprintf("*.%s", filepath.Base(asset.Filename)),
	)

	if err != nil {
		return fmt.Errorf("unable to create temp file for original asset: %w", err)
	}
	defer os.Remove(originalFile.Name())

	io.Copy(originalFile, s3Obj)

	probe, err := ffmpeg.Probe(originalFile.Name())
	if err != nil {
		return fmt.Errorf("unable to probe original video: %w", err)
	}

	var info Probe
	json.Unmarshal([]byte(probe), &info)

	err = processVideoThumbnail(ctx, minioClient, asset, originalFile, info)
	if err != nil {
		return fmt.Errorf("unable to process video asset thumbnail: %w", err)
	}

	err = processVideoPreview(ctx, minioClient, asset, originalFile, info)
	if err != nil {
		return fmt.Errorf("unable to process video asset preview: %w", err)
	}

	err = processVideoView(ctx, minioClient, asset, originalFile, info)
	if err != nil {
		return fmt.Errorf("unable to process video asset view: %w", err)
	}

	asset.Type = db.AssetTypeTVideo

	return nil
}

func processVideoView(
	ctx context.Context, minioClient *minio.Client, asset *db.Asset,
	originalFile *os.File, _ Probe,
) error {
	slog.Info("process video asset view media", slog.Any("id", asset.ID))
	err := ctx.Err()
	if err != nil {
		return fmt.Errorf("context cancelled: %w", err)
	}

	if asset.View == "" || asset.View == asset.Original {
		asset.View = createAssetKey("mp4")
	}
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

	_, err = minioClient.PutObject(
		ctx,
		os.Getenv("S3_BUCKET"),
		asset.View,
		outputFile,
		-1,
		minio.PutObjectOptions{
			ContentType: "video/mp4",
		},
	)

	return nil
}

func processVideoThumbnail(
	ctx context.Context,
	minioClient *minio.Client,
	asset *db.Asset,
	originalFile *os.File,
	info Probe,
) error {
	slog.Info("process video asset thumbnail", slog.Any("id", asset.ID))
	err := ctx.Err()
	if err != nil {
		return fmt.Errorf("context cancelled: %w", err)
	}

	image, err := extractVideoThumbnail(info, originalFile, asset, "vframes", "1")

	if err != nil {
		return fmt.Errorf("unable to extract image from video: %w", err)
	}

	err = processThumbnail(image)
	if err != nil {
		return fmt.Errorf("Unable to process thumbnail image: %w", err)
	}

	saveParams := vips.DefaultWebpsaveBufferOptions()
	saveParams.Q = THUMBNAIL_QUALITY

	buf, err := image.WebpsaveBuffer(saveParams)
	if err != nil {
		return fmt.Errorf("unable to write preview image: %w", err)
	}

	asset.ThumbnailHeight = int32(image.Height())
	asset.ThumbnailWidth = int32(image.Width())

	if asset.Thumbnail == "" {
		asset.Thumbnail = createAssetKey("webp")
	}

	err = putObjectToS3(ctx, minioClient, asset.Thumbnail, bytes.NewReader(buf), "image/webp")
	if err != nil {
		return fmt.Errorf("unable to put object to S3: %w", err)
	}
	return nil
}

func extractVideoThumbnail(
	info Probe,
	originalFile *os.File,
	asset *db.Asset,
	paramKey string,
	paramValue string,
) (image *vips.Image, error error) {
	duration, err := strconv.ParseFloat(info.Format.Duration, 10)
	if err != nil {
		err = fmt.Errorf("unable to parse duration: %w", err)
		return
	}

	outputFile, err := os.CreateTemp("", "*view.webp")
	if err != nil {
		err = fmt.Errorf("unable to create temp file to transcode: %w", err)
		return
	}
	defer os.Remove(outputFile.Name())

	// save thumbnail at 1/3 duration
	err = ffmpeg.
		Input(originalFile.Name(), ffmpeg.KwArgs{
			"ss": fmt.Sprintf("%f", duration/3),
		}).
		Output(outputFile.Name(), ffmpeg.KwArgs{
			"c:v":    "libwebp",
			paramKey: paramValue,
			"loop":   "0",
			// "quality": fmt.Sprintf("%d", THUMBNAIL_QUALITY),
			// "vf":      fmt.Sprintf("scale=-2:%d", THUMBNAIL_HEIGHT),
		}).OverWriteOutput().ErrorToStdOut().Run()

	if err != nil {
		err = fmt.Errorf("unable to create thumbnail asset for video asset: %w", err)
		return
	}

	videoDuration := time.Duration(duration) * time.Second
	asset.VideoDuration = pgtype.Interval{
		Microseconds: videoDuration.Microseconds(),
		Valid:        true,
	}

	image, err = vips.NewImageFromFile(outputFile.Name(), nil)
	if err != nil {
		err = fmt.Errorf("unable to read extracted file: %w", err)
		return
	}

	return
}

func processVideoPreview(
	ctx context.Context, minioClient *minio.Client, asset *db.Asset,
	originalFile *os.File, info Probe,
) error {
	slog.Info("process video preview", slog.Any("id", asset.ID))
	err := ctx.Err()
	if err != nil {
		return fmt.Errorf("context cancelled: %w", err)
	}

	slog.Info("process video asset thumbnail", slog.Any("id", asset.ID))

	image, err := extractVideoThumbnail(info, originalFile, asset, "t", "5")
	if err != nil {
		return fmt.Errorf("unable to extract preview image: %w", err)
	}

	err = processAnimatedPreview(image)
	if err != nil {
		return fmt.Errorf("unable to process preview image: %w", err)
	}

	saveParams := vips.DefaultWebpsaveBufferOptions()
	saveParams.Q = THUMBNAIL_QUALITY

	if asset.Preview == "" {
		asset.Preview = createAssetKey("webp")
	}

	buf, err := image.WebpsaveBuffer(saveParams)

	err = putObjectToS3(ctx, minioClient, asset.Preview, bytes.NewReader(buf), "image/webp")
	if err != nil {
		return fmt.Errorf("unable to put object to S3: %w", err)
	}
	return nil
}
