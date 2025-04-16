package minio

import (
	"context"
	"errors"
	"fmt"
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

	err := c.createBucket(ctx, public)
	if err != nil {
		return "", fmt.Errorf("could not create bucket: %w", err)
	}

	ct, err := findContentType(filePath)
	if err != nil {
		return "", fmt.Errorf("could not find content type: %w", err)
	}

	fileName, err := randomizeFileName(filePath)
	if err != nil {
		return "", fmt.Errorf("could not randomize file name: %w", err)
	}

	// TODO: rework this
	bucketName := bucketNamePrivate
	if public {
		bucketName = bucketNamePublic
	}

	info, err := c.client.FPutObject(ctx, bucketName, fileName, filePath, minio.PutObjectOptions{
		ContentType: ct,
	})
	if err != nil {
		return "", fmt.Errorf("could not fput object: %w", err)
	}

	if public {
		link := fmt.Sprintf("%s/%s/%s", c.client.EndpointURL().String(), bucketName, info.Key)
		return link, nil
	}

	link, err := c.client.PresignedGetObject(ctx, bucketName, fileName, objectExpiry, nil)
	if err != nil {
		return "", fmt.Errorf("could not get presigned url: %w", err)
	}

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

	exists, err := c.client.BucketExists(ctx, bucket)
	if err != nil {
		return fmt.Errorf("could not check if bucket exists: %w", err)
	}

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

	if !public {
		return nil
	}

	err = c.client.SetBucketPolicy(ctx, bucket, fmt.Sprintf(bucketPolicyPublicTemplate, bucket))
	if err != nil {
		return fmt.Errorf("could not set public bucket policy: %w", err)
	}

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
