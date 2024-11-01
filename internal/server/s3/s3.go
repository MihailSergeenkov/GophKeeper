package s3

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/MihailSergeenkov/GophKeeper/internal/server/config"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/constants"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/encrypt"
)

var (
	contentType = "application/octet-stream"
	pBacketName = "backetforuserid"
)

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

	return &S3{client: minioClient, settings: settings}, nil
}

// AddFile функция добавления файла в S3 хранилище.
func (fs S3) AddFile(ctx context.Context, file io.Reader, objectName string, objectSize int64) error {
	backetName, err := fs.CheckOrCreateBacket(ctx)
	if err != nil {
		return fmt.Errorf("failed to work with backet: %w", err)
	}

	options := minio.PutObjectOptions{ContentType: contentType}

	if fs.settings.SecureFiles {
		encryption := encrypt.DefaultPBKDF([]byte(fs.settings.SecretPassword), []byte(backetName+objectName))
		options.ServerSideEncryption = encryption
	}

	_, err = fs.client.PutObject(ctx, backetName, objectName, file, objectSize, options)
	if err != nil {
		return fmt.Errorf("failed to put file %w", err)
	}

	return nil
}

// GetFile функция получения файла из S3 хранилища.
func (fs S3) GetFile(ctx context.Context, objectName string) (io.ReadCloser, error) {
	backetName, err := fs.CheckOrCreateBacket(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to work with backet: %w", err)
	}

	options := minio.GetObjectOptions{}

	if fs.settings.SecureFiles {
		encryption := encrypt.DefaultPBKDF([]byte(fs.settings.SecretPassword), []byte(backetName+objectName))
		options.ServerSideEncryption = encryption
	}

	file, err := fs.client.GetObject(ctx, backetName, objectName, options)
	if err != nil {
		return nil, fmt.Errorf("failed to get file %w", err)
	}

	return file, nil
}

// CheckOrCreateBacket функция проверки существования и создания бакета в S3 хранилище.
func (fs S3) CheckOrCreateBacket(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(constants.KeyUserID).(int)
	if !ok {
		return "", errors.New("failed to fetch user id from context")
	}
	backetName := pBacketName + strconv.Itoa(userID)

	exists, err := fs.client.BucketExists(ctx, backetName)
	if err != nil {
		return "", fmt.Errorf("failed to check backet exist %w", err)
	}
	if exists {
		return backetName, nil
	}

	err = fs.client.MakeBucket(ctx, backetName, minio.MakeBucketOptions{Region: fs.settings.Region})
	if err != nil {
		return "", fmt.Errorf("failed to create backet %w", err)
	}

	return backetName, nil
}
