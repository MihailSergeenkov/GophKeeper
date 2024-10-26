package services

import (
	"context"
	"errors"
	"io"

	"github.com/MihailSergeenkov/GophKeeper/internal/models"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/config"
)

var (
	ErrNotFound                = errors.New("requested data no found")
	ErrUserMarkIsTooBig        = errors.New("user mark is too big")
	ErrUserDescriptionIsTooBig = errors.New("user description is too big")
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
