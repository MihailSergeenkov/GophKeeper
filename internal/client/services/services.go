package services

import (
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
