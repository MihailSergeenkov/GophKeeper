package services

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/MihailSergeenkov/GophKeeper/internal/models"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/storage"
)

const passwordDataType = "password"

var (
	ErrUserLoginIsTooBig    = errors.New("user login is too big")
	ErrUserPasswordIsTooBig = errors.New("user password is too big")

	passwordLoginMaxSize    = 100
	passwordPasswordMaxSize = 100
)

// AddPassword функция для добавления пароля пользователя.
func (s *Services) AddPassword(ctx context.Context, req models.AddPasswordRequest) (int, error) {
	if err := validateAddPasswordRequest(req); err != nil {
		return 0, failedValidateFields(err)
	}

	data := models.EncryptPasswordData{
		Login:    req.Login,
		Password: req.Password,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return 0, failedGenerateJSONData(err)
	}

	encData := s.crypter.EncryptData(jsonData)

	id, err := s.storage.AddUserData(ctx, encData, req.Mark, req.Description, passwordDataType)
	if err != nil {
		return 0, failedAddUserData(err)
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

		return resp, failedGetUserData(err)
	}

	jsonData, err := s.crypter.DecryptData(decData)
	if err != nil {
		return resp, failedDecryptData(err)
	}

	var encData models.EncryptPasswordData

	if err = json.Unmarshal(jsonData, &encData); err != nil {
		return resp, failedGenerateData(err)
	}

	resp.ID = id
	resp.Login = encData.Login
	resp.Password = encData.Password
	resp.Mark = mark
	resp.Description = description

	return resp, nil
}

func validateAddPasswordRequest(req models.AddPasswordRequest) error {
	if len([]rune(req.Login)) > passwordLoginMaxSize {
		return ErrUserLoginIsTooBig
	}
	if len([]rune(req.Password)) > passwordPasswordMaxSize {
		return ErrUserPasswordIsTooBig
	}
	if len([]rune(req.Mark)) > maxMarkSize {
		return ErrUserMarkIsTooBig
	}
	if len([]rune(req.Description)) > maxDescriptionSize {
		return ErrUserDescriptionIsTooBig
	}

	return nil
}
