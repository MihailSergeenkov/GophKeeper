package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MihailSergeenkov/GophKeeper/internal/models"
)

// RegisterUser сервис регистрации пользователя.
func RegisterUser(cfg Configurer, req models.RegisterUserRequest) error {
	const path = "/user/register"
	client := getClient(cfg)

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to create body: %w", err)
	}

	resp, err := client.R().
		SetHeader(ContentTypeHeader, JSONContentType).
		SetBody(body).
		Post(cfg.GetServerAPI() + path)

	if err != nil {
		return fmt.Errorf("failed request: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("response status: %s", resp.Status())
	}

	return nil
}

// LoginUser сервис аутентификации пользователя.
func LoginUser(cfg Configurer, req models.CreateUserTokenRequest) error {
	const path = "/user/token"
	client := getClient(cfg)

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to create body: %w", err)
	}

	userToken := models.CreateUserTokenResponse{}

	resp, err := client.R().
		SetHeader(ContentTypeHeader, JSONContentType).
		SetBody(body).
		SetResult(&userToken).
		Post(cfg.GetServerAPI() + path)

	if err != nil {
		return fmt.Errorf("failed request: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("response status: %s", resp.Status())
	}

	cfg.UpdateToken(userToken.AuthToken)
	return nil
}

// LogoutUser сервис удаления данных пользователя.
func LogoutUser(cfg Configurer) {
	cfg.UpdateToken("")
	cfg.UpdateData([]models.UserData{})
}
