package minio

import (
	"log/slog"
	"strings"

	"github.com/devusSs/minls/internal/log"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
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
		log.Warn("minio - new client", slog.String("msg", "insecure endpoint"))
	}

	return &Client{
		client: c,
	}, nil
}
