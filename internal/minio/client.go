package minio

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Client struct {
	logger *slog.Logger
	client *minio.Client
}

func NewClient(
	accessKey string,
	accessSecret string,
	endpoint string,
	logger *slog.Logger,
) (*Client, error) {
	if logger == nil {
		return nil, errors.New("logger cannot be nil")
	}

	secure := strings.Contains(endpoint, "https://")

	c, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, accessSecret, ""),
		Secure: secure,
	})
	if err != nil {
		return nil, err
	}

	if !secure {
		logger.Warn("minio - new client", slog.String("msg", "insecure endpoint"))
	}

	return &Client{
		logger: logger,
		client: c,
	}, nil
}
