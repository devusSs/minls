package env

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"

	"github.com/devusSs/minls/internal/log"
)

type Env struct {
	MinioEndpoint     string `json:"minio_endpoint,omitempty"`
	MinioAccessKey    string `json:"minio_access_key,omitempty"`
	MinioAccessSecret string `json:"minio_access_secret,omitempty"`
	YOURLSEndpoint    string `json:"yourls_endpoint,omitempty"`
	YOURLSSignature   string `json:"yourls_signature,omitempty"`
}

func Load() (*Env, error) {
	exe, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("could not find executable: %w", err)
	}

	log.Debug("env - Load", slog.String("action", "find_executable"), slog.String("exe_path", exe))

	envFilePath := filepath.Join(filepath.Dir(exe), ".env")
	log.Debug(
		"env - Load",
		slog.String("action", "godotenv_load"),
		slog.String("env_file_path", envFilePath),
	)

	// we can ignore the load error since we dont mind
	// whether an .env file is actually present or not
	err = godotenv.Load(envFilePath)
	log.Debug("env - Load", slog.String("action", "godotenv_load"), slog.Any("err", err))

	env := &Env{}

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

	log.Debug("env - loadKey", slog.String("key", key), slog.String("v", v))

	return v, nil
}
