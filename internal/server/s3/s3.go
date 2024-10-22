package s3

import (
	"context"
	"fmt"
	"io"

	"github.com/MihailSergeenkov/GophKeeper/internal/server/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/encrypt"
)

var contentType = "application/octet-stream"

// S3 структура для работы с S3 хранилищем приложения.
type S3 struct {
	client   *minio.Client
	settings *config.S3Settings
}

// NewClient функция инициализации S3 хранилища приложения.
func NewClient(ctx context.Context, settings *config.S3Settings) (*S3, error) {
	minioClient, err := minio.New(settings.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(settings.AccessKeyID, settings.SecretAccessKey, ""),
		Secure: settings.UseSSL,
		Region: settings.Region,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize minio client object %w", err)
	}

	exists, err := minioClient.BucketExists(ctx, settings.BucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to check backet exist %w", err)
	}
	if exists {
		return &S3{client: minioClient, settings: settings}, nil
	}

	err = minioClient.MakeBucket(ctx, settings.BucketName, minio.MakeBucketOptions{Region: "us-east-1"})
	if err != nil {
		return nil, fmt.Errorf("failed to create backet %w", err)
	}

	return &S3{client: minioClient, settings: settings}, nil
}

func (fs S3) AddFile(ctx context.Context, file io.Reader, objectName string, objectSize int64) error {
	options := minio.PutObjectOptions{ContentType: contentType}

	if fs.settings.SecureFiles {
		encryption := encrypt.DefaultPBKDF([]byte(fs.settings.SecretPassword), []byte(fs.settings.BucketName+objectName))
		options.ServerSideEncryption = encryption
	}

	_, err := fs.client.PutObject(ctx, fs.settings.BucketName, objectName, file, objectSize, options)
	if err != nil {
		return fmt.Errorf("failed to put file %w", err)
	}

	return nil
}

func (fs S3) GetFile(ctx context.Context, objectName string) (io.ReadCloser, error) {
	options := minio.GetObjectOptions{}

	if fs.settings.SecureFiles {
		encryption := encrypt.DefaultPBKDF([]byte(fs.settings.SecretPassword), []byte(fs.settings.BucketName+objectName))
		options.ServerSideEncryption = encryption
	}

	file, err := fs.client.GetObject(ctx, fs.settings.BucketName, objectName, options)
	if err != nil {
		return nil, fmt.Errorf("failed to get file %w", err)
	}

	return file, nil
}
