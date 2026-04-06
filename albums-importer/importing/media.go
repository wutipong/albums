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
}

func IsMediaFile(path string) bool {
	return slices.Contains(mediaExtensions, strings.ToLower(filepath.Ext(path)))
}
