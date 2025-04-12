package env

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	MinioEndpoint     string `json:"minio_endpoint,omitempty"`
	MinioAccessKey    string `json:"minio_access_key,omitempty"`
	MinioAccessSecret string `json:"minio_access_secret,omitempty"`
	YOURLSEndpoint    string `json:"yourls_endpoint,omitempty"`
	YOURLSSignature   string `json:"yourls_signature,omitempty"`
}

func Load() (*Env, error) {
	_ = godotenv.Load()

	env := &Env{}
	var err error

	env.MinioEndpoint, err = loadKey("MINIO_ENDPOINT")
	if err != nil {
		return nil, fmt.Errorf("could not get MINIO_ENDPOINT: %w", err)
	}

	env.MinioAccessKey, err = loadKey("MINIO_ACCESS_KEY")
	if err != nil {
		return nil, fmt.Errorf("could not get MINIO_ACCESS_KEY: %w", err)
	}

	env.MinioAccessSecret, err = loadKey("MINIO_ACCESS_SECRET")
	if err != nil {
		return nil, fmt.Errorf("could not get MINIO_ACCESS_SECRET: %w", err)
	}

	env.YOURLSEndpoint, err = loadKey("YOURLS_ENDPOINT")
	if err != nil {
		return nil, fmt.Errorf("could not get YOURLS_ENDPOINT: %w", err)
	}

	env.YOURLSSignature, err = loadKey("YOURLS_SIGNATURE")
	if err != nil {
		return nil, fmt.Errorf("could not get YOURLS_SIGNATURE: %w", err)
	}

	return env, nil
}

func loadKey(key string) (string, error) {
	v := os.Getenv(key)
	if v == "" {
		return "", fmt.Errorf("key %s could not be found", key)
	}

	return v, nil
}
