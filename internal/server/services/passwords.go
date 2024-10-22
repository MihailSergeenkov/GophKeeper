package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/MihailSergeenkov/GophKeeper/internal/server/models"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/storage"
)

const passwordDataType = "password"

var (
	ErrUserLoginIsTooBig    = errors.New("user login is too big")
	ErrUserPasswordIsTooBig = errors.New("user password is too big")
)

// AddPassword функция для добавления пароля пользователя.
func (s *Services) AddPassword(ctx context.Context, req models.AddPasswordRequest) (int, error) {
	if err := validateAddPasswordRequest(req); err != nil {
		return 0, fmt.Errorf("failed to validate fields %w", err)
	}

	data := models.EncryptPasswordData{
		Login:    req.Login,
		Password: req.Password,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return 0, fmt.Errorf("failed to generate json data %w", err)
	}

	encData := s.crypter.EncryptData(jsonData)

	id, err := s.storage.AddUserData(ctx, encData, req.Mark, req.Description, passwordDataType)
	if err != nil {
		return 0, fmt.Errorf("failed to add user data %w", err)
	}

	return id, nil
}

// GetPassword функция для получения пароля пользователя.
func (s *Services) GetPassword(ctx context.Context, id int) (models.Password, error) {
	resp := models.Password{}

	decData, mark, description, err := s.storage.GetUserData(ctx, id, passwordDataType)
	if err != nil {
		if errors.Is(err, storage.ErrUserDataNotFound) {
			return resp, ErrNotFound
		}

		return resp, fmt.Errorf("failed to get user data %w", err)
	}

	jsonData, err := s.crypter.DecryptData(decData)
	if err != nil {
		return resp, fmt.Errorf("failed to decrypt data %w", err)
	}

	var encData models.EncryptPasswordData

	if err = json.Unmarshal(jsonData, &encData); err != nil {
		return resp, fmt.Errorf("failed to generate data %w", err)
	}

	resp.ID = id
	resp.Login = encData.Login
	resp.Password = encData.Password
	resp.Mark = mark
	resp.Description = description

	return resp, nil
}

func validateAddPasswordRequest(req models.AddPasswordRequest) error {
	if len([]rune(req.Login)) > 100 {
		return ErrUserLoginIsTooBig
	}
	if len([]rune(req.Password)) > 100 {
		return ErrUserPasswordIsTooBig
	}
	if len([]rune(req.Mark)) > 50 {
		return ErrUserMarkIsTooBig
	}
	if len([]rune(req.Description)) > 50 {
		return ErrUserDescriptionIsTooBig
	}

	return nil
}
