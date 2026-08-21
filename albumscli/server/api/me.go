package api

import (
	"context"

	"github.com/wutipong/albums/albumscli/server/types"
)

type MeResponse struct {
	User    types.User    `json:"user"`
	Session types.Session `json:"session"`
}

func GetMe(ctx context.Context, server ServerConfig) (resp MeResponse, err error) {
	resp, err = Get[MeResponse](ctx, server, "api/me")
	return
}
