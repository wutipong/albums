package api

import (
	"context"
	"path"

	"github.com/wutipong/albums/albumscli/server/types"
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

type CreateAlbumRequest struct {
	Name string `json:"name"`
}

func CreateAlbum(
	ctx context.Context,
	server ServerConfig,
	name string,
) (resp types.Album, err error) {
	req := CreateAlbumRequest{Name: name}
	resp, err = Post[types.Album](ctx, server, path.Join("api", "album"), req)
	return
}

type DeleteAlbumResponse struct {
	Success bool `json:"success"`
}

func DeleteAlbum(
	ctx context.Context,
	server ServerConfig,
	id string,
) (resp DeleteAlbumResponse, err error) {

	resp, err = Delete[DeleteAlbumResponse](ctx, server, path.Join("api", "album", id))
	return
}

type PopulateAlbumCoverRequest struct {
	AssetID string `json:"asset_id"`
}

type PopulateAlbumCoverResponse struct {
	Status string `json:"status"`
}

func PopulateAlbumCover(
	ctx context.Context,
	server ServerConfig,
	albumID string,
) (resp PopulateAlbumCoverResponse, err error) {
	req := PopulateAlbumCoverRequest{
		AssetID: "",
	}
	resp, err = Post[PopulateAlbumCoverResponse](
		ctx, server, path.Join("api", "album", albumID, "cover"), req)
	return
}
