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

	"github.com/devusSs/minls/internal/log"
)

func (c *Client) UploadFile(ctx context.Context, filePath string, public bool) (string, error) {
	if ctx == nil {
		return "", errors.New("context cannot be nil")
	}

	err := c.createBucket(ctx, public)
	if err != nil {
		return "", fmt.Errorf("could not create bucket: %w", err)
	}

	ct, err := findContentType(filePath)
	if err != nil {
		return "", fmt.Errorf("could not find content type: %w", err)
	}

	log.Debug(
		"minio - *client.UploadFile",
		slog.String("action", "find_content_type"),
		slog.String("ct", ct),
	)

	fileName, err := randomizeFileName(filePath)
	if err != nil {
		return "", fmt.Errorf("could not randomize file name: %w", err)
	}

	log.Debug(
		"minio - *client.UploadFile",
		slog.String("action", "randomize_file_name"),
		slog.String("file_name", fileName),
	)

	bucketName := bucketNamePrivate
	if public {
		bucketName = bucketNamePublic
	}

	log.Debug(
		"minio - *client.UploadFile",
		slog.String("action", "set_bucket_name"),
		slog.String("bucket_name", bucketName),
	)

	info, err := c.client.FPutObject(ctx, bucketName, fileName, filePath, minio.PutObjectOptions{
		ContentType: ct,
	})
	if err != nil {
		return "", fmt.Errorf("could not fput object: %w", err)
	}

	log.Debug(
		"minio - *client.UploadFile",
		slog.String("action", "f_put_object"),
		slog.String("bucket_name", bucketName),
		slog.String("file_name", fileName),
		slog.String("file_path", filePath),
		slog.Any("info", info),
	)

	if public {
		link := fmt.Sprintf("%s/%s/%s", c.client.EndpointURL().String(), bucketName, info.Key)
		log.Debug(
			"minio - *client.UploadFile",
			slog.String("action", "return"),
			slog.String("warn", "link is public, returning early"),
			slog.String("link", link),
		)
		return link, nil
	}

	link, err := c.client.PresignedGetObject(ctx, bucketName, fileName, objectExpiry, nil)
	if err != nil {
		return "", fmt.Errorf("could not get presigned url: %w", err)
	}

	log.Debug(
		"minio - *client.UploadFile",
		slog.String("action", "presigned_get_object"),
		slog.String("bucket_name", bucketName),
		slog.String("file_name", fileName),
		slog.Duration("object_expiry", objectExpiry),
		slog.String("link", link.String()),
	)

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

	log.Debug(
		"minio - *client.createBucket",
		slog.String("action", "set_bucket_name"),
		slog.String("bucket_name", bucket),
	)

	exists, err := c.client.BucketExists(ctx, bucket)
	if err != nil {
		return fmt.Errorf("could not check if bucket exists: %w", err)
	}

	log.Debug(
		"minio - *client.createBucket",
		slog.String("action", "check_bucket_exists"),
		slog.Bool("bucket_exists", exists),
	)

	if exists {
		log.Debug(
			"minio - *client.createBucket",
			slog.String("action", "check_bucket_exists"),
			slog.String("warn", "bucket exists, skipping"),
		)
		return nil
	}

	err = c.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{
		Region:        bucketRegion,
		ObjectLocking: bucketObjectLocking,
	})
	if err != nil {
		return fmt.Errorf("could not create bucket: %w", err)
	}

	log.Debug(
		"minio - *client.createBucket",
		slog.String("action", "make_bucket"),
		slog.String("info", "created bucket"),
		slog.String("bucket_name", bucket),
		slog.String("bucket_region", bucketRegion),
		slog.Bool("bucket_object_locking", bucketObjectLocking),
	)

	if !public {
		log.Debug(
			"minio - *client.createBucket",
			slog.String("action", "set_policy"),
			slog.String("warn", "bucket not public, skipping"),
		)
		return nil
	}

	policy := fmt.Sprintf(bucketPolicyPublicTemplate, bucket)
	err = c.client.SetBucketPolicy(ctx, bucket, policy)
	if err != nil {
		return fmt.Errorf("could not set public bucket policy: %w", err)
	}

	log.Debug(
		"minio - *client.createBucket",
		slog.String("action", "set_policy"),
		slog.String("policy", policy),
	)

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

	log.Debug("minio - findContentType", slog.Any("mime", mime))

	return mime.String(), nil
}

func randomizeFileName(filePath string) (string, error) {
	file := filepath.Base(filePath)
	log.Debug("minio - randomizeFileName", slog.String("file", file))

	ext := filepath.Ext(file)
	log.Debug("minio - randomizeFileName", slog.String("ext", ext))

	uid, err := uuid.NewUUID()
	if err != nil {
		return "", fmt.Errorf("could not create uuid: %w", err)
	}

	log.Debug("minio - randomizeFileName", slog.Any("uid", uid))

	final := uid.String() + ext
	log.Debug("minio - randomizeFileName", slog.String("final", final))

	return final, nil
}
