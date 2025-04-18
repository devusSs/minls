package cli

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/devusSs/minls/internal/clip"
	"github.com/devusSs/minls/internal/env"
	"github.com/devusSs/minls/internal/log"
	"github.com/devusSs/minls/internal/minio"
	"github.com/devusSs/minls/internal/yourls"
)

func Upload() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err := initialize()
	if err != nil {
		return fmt.Errorf("could not initialize cli: %w", err)
	}

	log.Debug("cli - Upload", slog.String("action", "initialized"))

	env, err := env.Load()
	if err != nil {
		return fmt.Errorf("could not load env: %w", err)
	}

	log.Debug("cli - Upload", slog.String("action", "loaded_env"), slog.Any("env", env))

	if len(os.Args) != neededArgs {
		return fmt.Errorf("missing upload filepath or policy, len(os.Args) = %d", len(os.Args))
	}

	log.Debug("cli - Upload", slog.String("action", "checked_args"), slog.Int("args", len(os.Args)))

	fp, err := getFilePath()
	if err != nil {
		return fmt.Errorf("could not get file path: %w", err)
	}

	log.Debug("cli - Upload", slog.String("action", "got_file_path"), slog.String("fp", fp))

	p, err := getPolicy()
	if err != nil {
		return fmt.Errorf("could not get policy: %w", err)
	}

	log.Debug("cli - Upload", slog.String("action", "got_policy"), slog.String("p", p))

	mc, err := minio.NewClient(env.MinioAccessKey, env.MinioAccessSecret, env.MinioEndpoint)
	if err != nil {
		return fmt.Errorf("could not create minio client: %w", err)
	}

	log.Debug("cli - Upload", slog.String("action", "minio_client_init"))

	log.Debug(
		"cli - Upload",
		slog.String("action", "uploading_to_minio"),
		slog.String("fp", fp),
		slog.String("p", p),
	)

	minioLink, err := mc.UploadFile(ctx, fp, p == "public")
	if err != nil {
		return fmt.Errorf("could not upload file: %w", err)
	}

	log.Info(
		"cli - Upload",
		slog.String("action", "uploaded_to_minio"),
		slog.String("minio_link", minioLink),
	)

	yc := yourls.NewClient(env.YOURLSEndpoint, env.YOURLSSignature)

	log.Debug("cli - Upload", slog.String("action", "yourls_client_init"))

	link, err := yc.Shorten(ctx, minioLink)
	if err != nil {
		return fmt.Errorf("could not shorten link: %w", err)
	}

	log.Info("cli - Upload", slog.String("action", "shortened_url"), slog.String("link", link))

	// TODO: setup & write to storage

	err = clip.Init()
	if err != nil {
		return fmt.Errorf("clipboard could not be initialized: %w", err)
	}

	log.Debug("cli - Upload", slog.String("action", "clip_init"))

	err = clip.Write(link)
	if err != nil {
		return fmt.Errorf("could not write to clipboard: %w", err)
	}

	log.Info("cli - Upload", slog.String("action", "clip_write"), slog.String("link", link))

	return nil
}

const neededArgs = 4

func getFilePath() (string, error) {
	fp := os.Args[2]
	if fp == "" {
		return "", errors.New("no file path provided")
	}

	log.Debug("getFilePath", slog.String("action", "check_fp"), slog.String("fp", fp))

	_, err := os.Stat(fp)
	if err != nil {
		return "", fmt.Errorf("file could not be found: %w", err)
	}

	return fp, nil
}

func getPolicy() (string, error) {
	p := os.Args[3]
	if p != "public" && p != "private" {
		return "", fmt.Errorf("invalid policy provided: %s", p)
	}

	log.Debug("getPolicy", slog.String("action", "check_p"), slog.String("p", p))

	return p, nil
}
