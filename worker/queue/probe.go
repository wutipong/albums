package queue

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
)

type Stream struct {
	Index       int    `json:"index"`
	CodecName   string `json:"codec_name"`
	CodecType   string `json:"codec_type"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	PixelFormat string `json:"pix_fmt"`
}

type Format struct {
	Filename   string `json:"filename"`
	FormatName string `json:"format_name"`
	Duration   string `json:"duration"`
}

type Probe struct {
	Streams []Stream `json:"streams"`
	Format  Format   `json:"format"`
}

func (info Probe) Video() (s Stream, err error) {
	slog.Debug("probe", "info", info)
	idx := slices.IndexFunc(info.Streams, func(s Stream) bool {
		return s.CodecType == "video"
	})

	if idx == -1 {
		err = fmt.Errorf("no video stream found")
		return
	}

	s = info.Streams[idx]
	return
}

func (info Probe) Audio() (s Stream, err error) {
	slog.Debug("probe", "info", info)
	idx := slices.IndexFunc(info.Streams, func(s Stream) bool {
		return s.CodecType == "audio"
	})

	if idx == -1 {
		err = fmt.Errorf("no audio stream found")
		return
	}

	s = info.Streams[idx]
	return
}

func validateVideo(ctx context.Context, info Probe) bool {
	switch info.Format.FormatName {
	case "mov,mp4,m4a,3gp,3g2,mj2":
		return validateMP4(ctx, info)
	case "matroska,webm":
		return validateWebM(ctx, info)
	}

	return false
}

func validateWebM(ctx context.Context, info Probe) bool {
	videoStream, err := info.Video()
	if err != nil {
		return false
	}

	if videoStream.CodecName != "vp8" &&
		videoStream.CodecName != "vp9" {
		return false
	}

	if videoStream.PixelFormat != "yuv420p" {
		return false
	}

	audioStream, err := info.Video()
	if err != nil {
		return false
	}

	if audioStream.CodecName != "opus" &&
		audioStream.CodecName != "vorbis" {
		return false
	}

	return true

}

func validateMP4(ctx context.Context, info Probe) bool {
	videoStreamIdx := slices.IndexFunc(info.Streams, func(s Stream) bool {
		return s.CodecType == "video"
	})

	if videoStreamIdx == -1 {
		return false
	}

	videoStream := info.Streams[videoStreamIdx]

	if videoStream.CodecName != "h264" {
		return false
	}

	if videoStream.PixelFormat != "yuv420p" {
		return false
	}

	audioStreamIdx := slices.IndexFunc(info.Streams, func(s Stream) bool {
		return s.CodecType == "audio"
	})

	if audioStreamIdx == -1 {
		return false
	}
	audioStream := info.Streams[audioStreamIdx]

	if audioStream.CodecName != "aac" &&
		audioStream.CodecName != "mp3" {
		return false
	}

	return true

}
