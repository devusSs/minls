package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/devusSs/minls/internal/clip"
	"github.com/devusSs/minls/internal/env"
	"github.com/devusSs/minls/internal/minio"
	"github.com/devusSs/minls/internal/yourls"
)

// TODO: readd logging
func Upload() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err := initialize()
	if err != nil {
		return fmt.Errorf("could not initialize cli: %w", err)
	}

	env, err := env.Load()
	if err != nil {
		return fmt.Errorf("could not load env: %w", err)
	}

	if len(os.Args) != 4 {
		return fmt.Errorf("missing upload filepath or policy, len(os.Args) = %d", len(os.Args))
	}

	fp, err := getFilePath()
	if err != nil {
		return fmt.Errorf("could not get file path: %w", err)
	}

	p, err := getPolicy()
	if err != nil {
		return fmt.Errorf("could not get policy: %w", err)
	}

	mc, err := minio.NewClient(env.MinioAccessKey, env.MinioAccessSecret, env.MinioEndpoint)
	if err != nil {
		return fmt.Errorf("could not create minio client: %w", err)
	}

	minioLink, err := mc.UploadFile(ctx, fp, p == "public")
	if err != nil {
		return fmt.Errorf("could not upload file: %w", err)
	}

	yc := yourls.NewClient(env.YOURLSEndpoint, env.YOURLSSignature)
	link, err := yc.Shorten(ctx, minioLink)
	if err != nil {
		return fmt.Errorf("could not shorten link: %w", err)
	}

	// TODO: setup & write to storage

	err = clip.Init()
	if err != nil {
		return fmt.Errorf("clipboard could not be initialized: %w", err)
	}

	err = clip.Write(link)
	if err != nil {
		return fmt.Errorf("could not write to clipboard: %w", err)
	}

	return nil
}

func getFilePath() (string, error) {
	fp := os.Args[2]
	if fp == "" {
		return "", errors.New("no file path provided")
	}

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

	return p, nil
}
