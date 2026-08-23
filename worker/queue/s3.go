package queue

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/minio/minio-go/v7"
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

func putObjectToS3(
	ctx context.Context,
	minioClient *minio.Client,
	key string,
	reader *bytes.Reader,
	contentType string,
) error {
	_, err := minioClient.PutObject(
		ctx, os.Getenv("S3_BUCKET"),
		key,
		reader,
		reader.Size(),
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)

	return err
}

func getObjectFromS3(ctx context.Context,
	minioClient *minio.Client,
	key string) (*minio.Object, error) {

	return minioClient.GetObject(
		ctx,
		os.Getenv("S3_BUCKET"),
		key,
		minio.GetObjectOptions{},
	)
}
