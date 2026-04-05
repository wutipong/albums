package api

import (
	"context"
	"path"

	"github.com/wutipong/albums/albums-importer/server/types"
)

type AlbumListResponse struct {
	Albums []types.Album `json:"albums"`
}

func GetAlbumList(ctx context.Context, server ServerConfig) (resp AlbumListResponse, err error) {
	resp, err = Get[AlbumListResponse](ctx, server, "api/album")
	return
}

type AlbumDetailResponse struct {
	types.Album
	Assets []string `json:"assets"`
}

func GetAlbumDetail(ctx context.Context, server ServerConfig, albumID string) (resp AlbumDetailResponse, err error) {
	resp, err = Get[AlbumDetailResponse](ctx, server, path.Join("api", "album", albumID))
	return
}
