package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MihailSergeenkov/GophKeeper/internal/client/requests"
	"github.com/MihailSergeenkov/GophKeeper/internal/models"
)

// RegisterUser сервис регистрации пользователя.
func (s *Services) RegisterUser(req models.RegisterUserRequest) error {
	const path = "/user/register"

	body, err := json.Marshal(req)
	if err != nil {
		return failedCreateBody(err)
	}

	resp, err := s.httpRequests.Post(
		s.cfg.GetServerAPI()+path,
		requests.WithHeader(ContentTypeHeader, JSONContentType),
		requests.WithBody(body),
	)
	if err != nil {
		return failedRequest(err)
	}
	if resp.StatusCode() != http.StatusOK {
		return failedResponseStatus(resp.Status())
	}

	return nil
}

// LoginUser сервис аутентификации пользователя.
func (s *Services) LoginUser(req models.CreateUserTokenRequest) error {
	const path = "/user/token"

	body, err := json.Marshal(req)
	if err != nil {
		return failedCreateBody(err)
	}

	userToken := models.CreateUserTokenResponse{}

	resp, err := s.httpRequests.Post(
		s.cfg.GetServerAPI()+path,
		requests.WithHeader(ContentTypeHeader, JSONContentType),
		requests.WithBody(body),
		requests.WithResult(&userToken),
	)
	if err != nil {
		return failedRequest(err)
	}
	if resp.StatusCode() != http.StatusOK {
		return failedResponseStatus(resp.Status())
	}

	if err := s.cfg.UpdateToken(userToken.AuthToken); err != nil {
		return fmt.Errorf("failed to update auth token: %w", err)
	}
	return nil
}

// LogoutUser сервис удаления данных пользователя.
func (s *Services) LogoutUser() error {
	if err := s.cfg.UpdateToken(""); err != nil {
		return fmt.Errorf("failed to update auth token: %w", err)
	}

	if err := s.cfg.UpdateData([]models.UserData{}); err != nil {
		return fmt.Errorf("failed to update data: %w", err)
	}

	return nil
}
