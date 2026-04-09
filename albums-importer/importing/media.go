package importing

import (
	"path/filepath"
	"slices"
	"strings"
)

var mediaExtensions = []string{
	".jpg",
	".jpeg",
	".png",
	".gif",
	".svg",
	".tiff",
	".webp",

	// movie
	".mp4",
	".m4v",
	".webm",
}

func IsMediaFile(path string) bool {
	return slices.Contains(mediaExtensions, strings.ToLower(filepath.Ext(path)))
}
