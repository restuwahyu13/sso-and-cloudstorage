package packages

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
)

type minioConfig struct {
	MinioClient     *minio.Client
	endpoint        string
	accessKeyId     string
	secretAccessKey string
	secureProtocol  string
}

func NewMinio() *minioConfig {
	config := minioConfig{}
	config.endpoint = GetString("MINIO_ENDPOINT")
	config.accessKeyId = GetString("MINIO_ACCESS_KEY_ID")
	config.secretAccessKey = GetString("MINIO_SECRET_ACCESS_KEY")
	config.secureProtocol = GetString("MINIO_SECURE")

	secureProtocol, _ := strconv.ParseBool(config.secureProtocol)
	client, _ := minio.New(config.endpoint, &minio.Options{
		Creds:           credentials.NewStaticV4(config.accessKeyId, config.secureProtocol, ""),
		Secure:          secureProtocol,
		TrailingHeaders: true,
	})

	return &minioConfig{MinioClient: client}
}

func (h *minioConfig) bucketExists(ctx context.Context, bucketName string) (bool, error) {
	res, err := h.MinioClient.BucketExists(ctx, bucketName)

	if err != nil {
		defer logrus.Errorf("bucketExists error: %s", err.Error())
		return res, err
	}

	return res, nil
}

// MakeBucket creates a new bucket with bucketName with a context to control cancellations and timeouts
func (h *minioConfig) MakeBucket(ctx context.Context, bucketName string) (interface{}, error) {
	checkBucket, err := h.bucketExists(ctx, bucketName)

	if err != nil {
		return nil, err
	}

	if !checkBucket {
		err := h.MinioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})

		if err != nil {
			defer logrus.Errorf("MakeBucket error: %s", err.Error())
			return nil, err
		}
	}

	return fmt.Sprintf("Created bucket name %s success", bucketName), nil
}

// ListObjects returns objects list after evaluating the passed options
func (h *minioConfig) ListObjects(ctx context.Context, bucketName string) (res <-chan minio.ObjectInfo, err error) {
	checkBucket, err := h.bucketExists(ctx, bucketName)
	if err != nil {
		return res, err
	}

	if checkBucket {
		minioRes := h.MinioClient.ListObjects(ctx, bucketName, minio.ListObjectsOptions{WithVersions: true, WithMetadata: true})
		res = minioRes
		err = nil
	}

	return res, nil
}

// ListBuckets list all buckets owned by this authenticated user
func (h *minioConfig) ListBucket(ctx context.Context) ([]minio.BucketInfo, error) {
	res, err := h.MinioClient.ListBuckets(ctx)

	if err != nil {
		defer logrus.Errorf("ListBucket error: %s", err.Error())
		return res, err
	}

	return res, nil
}

// GetObject wrapper function that accepts a request context
func (h *minioConfig) GetObject(ctx context.Context, bucketName, objectName string) (res *minio.Object, err error) {
	checkBucket, err := h.bucketExists(ctx, bucketName)
	if err != nil {
		return res, err
	}

	if checkBucket {
		minioRes, minioErr := h.MinioClient.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{Checksum: true})
		if minioErr != nil {
			defer logrus.Errorf("GetObject error: %s", minioErr.Error())
			return res, minioErr
		}

		res = minioRes
		err = nil
	}

	return res, nil
}

// PutObject creates an object in a bucket.
func (h *minioConfig) PutObject(ctx context.Context, bucketName, objectName string, reader io.Reader) (res minio.UploadInfo, err error) {
	checkBucket, err := h.bucketExists(ctx, bucketName)
	if err != nil {
		return res, err
	}

	if checkBucket {
		minioRes, minioErr := h.MinioClient.PutObject(ctx, bucketName, objectName, reader, int64(5242880), minio.PutObjectOptions{})
		if minioErr != nil {
			defer logrus.Errorf("PutObject error: %s", minioErr.Error())
			return res, minioErr
		}

		res = minioRes
		err = nil
	}

	return res, nil
}

// FPutObject - Create an object in a bucket, with contents from file at filePath. Allows request cancellation.
func (h *minioConfig) FGetObject(ctx context.Context, bucketName, objectName, filePath string) error {
	checkBucket, err := h.bucketExists(ctx, bucketName)
	if err != nil {
		return err
	}

	if checkBucket {
		err := h.MinioClient.FGetObject(ctx, bucketName, objectName, filePath, minio.GetObjectOptions{Checksum: true})
		if err != nil {
			defer logrus.Errorf("FGetObject error: %s", err.Error())
			return err
		}
	}

	return nil
}

// FPutObject - Create an object in a bucket, with contents from file at filePath. Allows request cancellation.
func (h *minioConfig) FPutObject(ctx context.Context, bucketName, objectName, filePath string) (res minio.UploadInfo, err error) {
	checkBucket, err := h.bucketExists(ctx, bucketName)
	if err != nil {
		return res, err
	}

	if checkBucket {
		minioRes, minioErr := h.MinioClient.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{})
		if minioErr != nil {
			defer logrus.Errorf("FPutObject error: %s", minioErr.Error())
			return res, minioErr
		}

		res = minioRes
		err = nil
	}

	return res, nil
}

// RemoveBucket deletes the bucket name
func (h *minioConfig) RemoveBucket(ctx context.Context, bucketName string) error {
	checkBucket, err := h.bucketExists(ctx, bucketName)
	if err != nil {
		return err
	}

	if checkBucket {
		err := h.MinioClient.RemoveBucket(ctx, bucketName)
		if err != nil {
			defer logrus.Errorf("RemoveBucket error: %s", err.Error())
			return err
		}
	}

	return nil
}

// RemoveObject removes an object from a bucket
func (h *minioConfig) RemoveObject(ctx context.Context, bucketName, objectName string) error {
	checkBucket, err := h.bucketExists(ctx, bucketName)
	if err != nil {
		return err
	}

	if checkBucket {
		err := h.MinioClient.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{ForceDelete: true})
		if err != nil {
			defer logrus.Errorf("RemoveObject error: %s", err.Error())
			return err
		}
	}

	return nil
}
