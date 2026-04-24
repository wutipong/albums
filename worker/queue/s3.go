package queue

import (
	"fmt"
	"net/url"
	"strings"
)

func GetMinioEndpoint(input string) (endpoint string, secure bool, err error) {
	endpointUrl, err := url.Parse(input)
	if err != nil {
		err = fmt.Errorf("unable to parse input url")
		return
	}

	scheme := endpointUrl.Scheme

	if scheme == "https" {
		secure = true
	} else {
		secure = false
	}
	endpoint = input
	endpoint = strings.TrimPrefix(endpoint, scheme)
	endpoint = strings.TrimPrefix(endpoint, "://")

	return
}
