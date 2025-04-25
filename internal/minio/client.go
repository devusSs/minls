package minio

import (
	"log/slog"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/devusSs/minls/internal/log"
)

type Client struct {
	client *minio.Client
}

func NewClient(
	accessKey string,
	accessSecret string,
	endpoint string,
) (*Client, error) {
	secure := strings.Contains(endpoint, "https://")
	endpoint = strings.Replace(endpoint, "https://", "", 1)
	endpoint = strings.Replace(endpoint, "http://", "", 1)

	c, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, accessSecret, ""),
		Secure: secure,
	})
	if err != nil {
		return nil, err
	}

	if !secure {
		log.Warn("minio - NewClient", slog.String("warn", "endpoint not secure"))
	}

	return &Client{
		client: c,
	}, nil
}
