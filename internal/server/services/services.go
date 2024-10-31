package services

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/MihailSergeenkov/GophKeeper/internal/models"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/config"
)

var (
	ErrNotFound                = errors.New("requested data no found")
	ErrUserMarkIsTooBig        = errors.New("user mark is too big")
	ErrUserDescriptionIsTooBig = errors.New("user description is too big")

	maxMarkSize        = 100
	maxDescriptionSize = 3000
)

// Services структура для работы с сервисами приложения.
type Services struct {
	storage     Storager
	fileStorage FileStorager
	crypter     Crypter
	settings    *config.Settings
}

// Storager интерфейс для хранилища данных.
type Storager interface {
	Ping(ctx context.Context) error
	AddUser(ctx context.Context, userLogin string, userPassword []byte) error
	GetUserByLogin(ctx context.Context, userLogin string) (models.User, error)
	FetchUserData(ctx context.Context) ([]models.UserData, error)
	AddUserData(ctx context.Context, encData []byte, mark string, description string, dataType string) (int, error)
	GetUserData(ctx context.Context, id int, dataType string) ([]byte, string, string, error)
	GetFileUserData(ctx context.Context, fileMark string) ([]byte, error)
}

// Crypter интерфейс для криптографии.
type Crypter interface {
	EncryptData(data []byte) []byte
	DecryptData(data []byte) ([]byte, error)
}

// FileStorager интерфейс для файлового хранилища данных.
type FileStorager interface {
	AddFile(ctx context.Context, file io.Reader, objectName string, objectSize int64) error
	GetFile(ctx context.Context, objectName string) (io.ReadCloser, error)
}

// NewServices функция инициализации сервисов приложения.
func NewServices(storage Storager, fileStorage FileStorager, crypter Crypter, settings *config.Settings) *Services {
	return &Services{
		storage:     storage,
		fileStorage: fileStorage,
		crypter:     crypter,
		settings:    settings,
	}
}

// failedValidateFields оберта ошибки валидации.
func failedValidateFields(err error) error {
	return fmt.Errorf("failed to validate fields %w", err)
}

// failedGenerateJSONData оберта ошибки генерации json.
func failedGenerateJSONData(err error) error {
	return fmt.Errorf("failed to generate json data %w", err)
}

// failedAddUserData оберта ошибки добавления данных пользователя.
func failedAddUserData(err error) error {
	return fmt.Errorf("failed to add user data %w", err)
}

// failedGetUserData оберта ошибки получения данных пользователя.
func failedGetUserData(err error) error {
	return fmt.Errorf("failed to get user data %w", err)
}

// failedDecryptData оберта ошибки расшифрования данных пользователя.
func failedDecryptData(err error) error {
	return fmt.Errorf("failed to decrypt data %w", err)
}

// failedGenerateData оберта ошибки генерации данных пользователя.
func failedGenerateData(err error) error {
	return fmt.Errorf("failed to generate data %w", err)
}
