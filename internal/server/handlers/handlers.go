package handlers

import (
	"context"

	"github.com/MihailSergeenkov/GophKeeper/internal/models"

	"go.uber.org/zap/zapcore"
)

var (
	encRespErrStr     = "error encoding response"
	readReqErrStr     = "failed to read request body"
	ContentTypeHeader = "Content-Type"
	JSONContentType   = "application/json"
)

// Handlers структура для работы с обработчиками HTTP запросов приложения.
type Handlers struct {
	services Servicer
	logger   Logger
}

// Servicer интерфейс для сервисов приложения.
type Servicer interface {
	Ping(ctx context.Context) error
	RegisterUser(ctx context.Context, req models.RegisterUserRequest) error
	CreateUserToken(ctx context.Context, req models.CreateUserTokenRequest) (models.CreateUserTokenResponse, error)
	FetchUserData(ctx context.Context) ([]models.UserData, error)
	AddPassword(ctx context.Context, req models.AddPasswordRequest) (int, error)
	GetPassword(ctx context.Context, id int) (models.Password, error)
	AddCard(ctx context.Context, req *models.AddCardRequest) (int, error)
	GetCard(ctx context.Context, id int) (models.Card, error)
	AddText(ctx context.Context, req models.AddTextRequest) (int, error)
	GetText(ctx context.Context, id int) (models.Text, error)
	AddFile(ctx context.Context, req models.AddFileRequest) (int, error)
	GetFile(ctx context.Context, id int) (models.File, error)
}

// Logger интерфейс для логгера приложения.
type Logger interface {
	Error(msg string, fields ...zapcore.Field)
}

// NewHandlers функция инициализации обработчиков HTTP запросов приложения.
func NewHandlers(services Servicer, logger Logger) *Handlers {
	return &Handlers{
		services: services,
		logger:   logger,
	}
}
