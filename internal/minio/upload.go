package minio

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

func (c *Client) UploadFile(ctx context.Context, filePath string, public bool) (string, error) {
	if ctx == nil {
		return "", errors.New("context cannot be nil")
	}

	c.logger.Info(
		"minio - uploading file",
		slog.String("file_path", filePath),
		slog.Bool("public", public),
	)

	err := c.createBucket(ctx, public)
	if err != nil {
		return "", fmt.Errorf("could not create bucket: %w", err)
	}

	ct, err := findContentType(filePath)
	if err != nil {
		return "", fmt.Errorf("could not find content type: %w", err)
	}

	c.logger.Info("minio - uploading file", slog.String("content-type", ct))

	fileName, err := randomizeFileName(filePath)
	if err != nil {
		return "", fmt.Errorf("could not randomize file name: %w", err)
	}

	c.logger.Info("minio - uploading file", slog.String("radomized_file_name", fileName))

	// TODO: rework this
	bucketName := bucketNamePrivate
	if public {
		bucketName = bucketNamePublic
	}

	c.logger.Debug("minio - uploading file", slog.String("bucket_name", bucketName))

	info, err := c.client.FPutObject(ctx, bucketName, fileName, filePath, minio.PutObjectOptions{
		ContentType: ct,
	})
	if err != nil {
		return "", fmt.Errorf("could not fput object: %w", err)
	}

	c.logger.Debug("minio - uploading file", slog.Any("info", info))

	if public {
		link := fmt.Sprintf("%s/%s/%s", c.client.EndpointURL().String(), bucketName, info.Key)
		c.logger.Info("minio - uploading file", slog.String("link", link))
		return link, nil
	}

	link, err := c.client.PresignedGetObject(ctx, bucketName, fileName, objectExpiry, nil)
	if err != nil {
		return "", fmt.Errorf("could not get presigned url: %w", err)
	}

	c.logger.Info("minio - uploading file", slog.String("link", link.String()))

	return link.String(), nil
}

func (c *Client) createBucket(ctx context.Context, public bool) error {
	if ctx == nil {
		return errors.New("context cannot be nil")
	}

	bucket := bucketNamePrivate
	if public {
		bucket = bucketNamePublic
	}

	c.logger.Debug("minio - create bucket", slog.String("bucket_name", bucket))

	exists, err := c.client.BucketExists(ctx, bucket)
	if err != nil {
		return fmt.Errorf("could not check if bucket exists: %w", err)
	}

	c.logger.Debug("minio - create bucket", slog.Bool("bucket_exists", exists))

	if exists {
		return nil
	}

	err = c.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{
		Region:        bucketRegion,
		ObjectLocking: bucketObjectLocking,
	})
	if err != nil {
		return fmt.Errorf("could not create bucket: %w", err)
	}

	c.logger.Debug("minio - create bucket", slog.String("msg", "created bucket"))

	if !public {
		return nil
	}

	err = c.client.SetBucketPolicy(ctx, bucket, fmt.Sprintf(bucketPolicyPublicTemplate, bucket))
	if err != nil {
		return fmt.Errorf("could not set public bucket policy: %w", err)
	}

	c.logger.Debug("minio - create bucket", slog.String("msg", "set public bucket policy"))

	return nil
}

const (
	bucketNamePublic           = "minls-public"
	bucketNamePrivate          = "minls-private"
	bucketRegion               = "us-east-1"
	bucketObjectLocking        = false
	bucketPolicyPublicTemplate = `{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Sid": "AddPerm",
				"Effect": "Allow",
				"Principal": {
					"AWS": [
						"*"
					]
				},
				"Action": [
					"s3:GetObject"
				],
				"Resource": [
					"arn:aws:s3:::%s/*"
				]
			}
		]
	}`
	objectExpiry = 7 * 24 * time.Hour
)

func findContentType(filePath string) (string, error) {
	mime, err := mimetype.DetectFile(filePath)
	if err != nil {
		return "", err
	}

	return mime.String(), nil
}

func randomizeFileName(filePath string) (string, error) {
	file := filepath.Base(filePath)
	ext := filepath.Ext(file)

	uid, err := uuid.NewUUID()
	if err != nil {
		return "", fmt.Errorf("could not create uuid: %w", err)
	}

	return uid.String() + ext, nil
}
