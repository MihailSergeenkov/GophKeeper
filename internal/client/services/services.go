package services

import (
	"fmt"
	"time"

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

// getClient сервис получения HTTP клиента.
func getClient(cfg Configurer) *resty.Client {
	return resty.
		New().
		SetTimeout(time.Duration(cfg.GetRequestTimeout()) * time.Second).
		SetRetryCount(cfg.GetRequestRetry())
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
