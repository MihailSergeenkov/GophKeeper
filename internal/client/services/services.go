package services

import (
	"fmt"

	"github.com/MihailSergeenkov/GophKeeper/internal/client/requests"
	"github.com/MihailSergeenkov/GophKeeper/internal/models"
	"github.com/go-resty/resty/v2"
)

const (
	ContentTypeHeader   = "Content-Type"
	AuthHeader          = "X-AUTH-TOKEN"
	JSONContentType     = "application/json"
	FormDataContentType = "multipart/form-data"
)

// Configurer интерфейс для конфигурации.
type Configurer interface {
	GetServerAPI() string
	GetRequestRetry() int
	GetRequestTimeout() int
	GetToken() string
	GetData() map[string]models.UserData
	UpdateToken(token string) error
	UpdateData(data []models.UserData) error
	AddData(data models.UserData) error
}

// Requester интерфейс для HTTP запросов.
type Requester interface {
	Get(url string, opts ...requests.RequestOptionFunc) (*resty.Response, error)
	Post(url string, opts ...requests.RequestOptionFunc) (*resty.Response, error)
}

// Services структура для работы с сервисами клиента.
type Services struct {
	cfg          Configurer
	httpRequests Requester
}

// Init функция инициализации сервисов клиента.
func Init(cfg Configurer, httpRequests Requester) *Services {
	return &Services{
		cfg:          cfg,
		httpRequests: httpRequests,
	}
}

// failedRequest обертка ошибки запроса.
func failedRequest(err error) error {
	return fmt.Errorf("failed request: %w", err)
}

// failedResponseStatus обертка ошибки неправильного статуса.
func failedResponseStatus(status string) error {
	return fmt.Errorf("response status: %s", status)
}

// failedCreateBody обертка ошибки генерации тела запроса.
func failedCreateBody(err error) error {
	return fmt.Errorf("failed to create body: %w", err)
}

// failedDumpData обертка ошибки кеширования данных.
func failedDumpData(err error) error {
	return fmt.Errorf("failed to dump data: %w", err)
}
