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
		return failedCreateBody(err)
	}

	resp, err := client.R().
		SetHeader(ContentTypeHeader, JSONContentType).
		SetBody(body).
		Post(cfg.GetServerAPI() + path)

	if err != nil {
		return failedRequest(err)
	}
	if resp.StatusCode() != http.StatusOK {
		return failedResponseStatus(resp.Status())
	}

	return nil
}

// LoginUser сервис аутентификации пользователя.
func LoginUser(cfg Configurer, req models.CreateUserTokenRequest) error {
	const path = "/user/token"
	client := getClient(cfg)

	body, err := json.Marshal(req)
	if err != nil {
		return failedCreateBody(err)
	}

	userToken := models.CreateUserTokenResponse{}

	resp, err := client.R().
		SetHeader(ContentTypeHeader, JSONContentType).
		SetBody(body).
		SetResult(&userToken).
		Post(cfg.GetServerAPI() + path)

	if err != nil {
		return failedRequest(err)
	}
	if resp.StatusCode() != http.StatusOK {
		return failedResponseStatus(resp.Status())
	}

	if err := cfg.UpdateToken(userToken.AuthToken); err != nil {
		return fmt.Errorf("failed to update auth token: %w", err)
	}
	return nil
}

// LogoutUser сервис удаления данных пользователя.
func LogoutUser(cfg Configurer) error {
	if err := cfg.UpdateToken(""); err != nil {
		return fmt.Errorf("failed to update auth token: %w", err)
	}

	if err := cfg.UpdateData([]models.UserData{}); err != nil {
		return fmt.Errorf("failed to update data: %w", err)
	}

	return nil
}
