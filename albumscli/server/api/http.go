package api

import (
	"net/url"

	"github.com/imroc/req/v3"
)

type ServerConfig struct {
	URL    *url.URL
	DryRun bool
	APIKey string
}

type ErrorResponse struct {
	Success string `json:"success"`
	Message string `json:"error"`
}

func (err ErrorResponse) Error() string {
	return err.Message
}

func NewClient(config ServerConfig) *req.Client {
	return req.C().
		SetCommonHeader("x-api-key", config.APIKey).
		SetBaseURL(config.URL.String())
}
