package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/MihailSergeenkov/GophKeeper/internal/models"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/storage"
)

const textDataType = "text"

var ErrUserTextDataIsTooBig = errors.New("user text data is too big")

// AddText функция для добавления текста пользователя.
func (s *Services) AddText(ctx context.Context, req models.AddTextRequest) (int, error) {
	if err := validateAddTextRequest(req); err != nil {
		return 0, fmt.Errorf("failed to validate fields %w", err)
	}

	data := models.EncryptTextData{
		Data: req.Data,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return 0, fmt.Errorf("failed to generate json data %w", err)
	}

	encData := s.crypter.EncryptData(jsonData)

	id, err := s.storage.AddUserData(ctx, encData, req.Mark, req.Description, textDataType)
	if err != nil {
		return 0, fmt.Errorf("failed to add user data %w", err)
	}

	return id, nil
}

// v функция для получения текста пользователя.
func (s *Services) GetText(ctx context.Context, id int) (models.Text, error) {
	resp := models.Text{}

	decData, mark, description, err := s.storage.GetUserData(ctx, id, textDataType)
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

	var encData models.EncryptTextData

	if err = json.Unmarshal(jsonData, &encData); err != nil {
		return resp, fmt.Errorf("failed to generate data %w", err)
	}

	resp.ID = id
	resp.Data = encData.Data
	resp.Mark = mark
	resp.Description = description

	return resp, nil
}

func validateAddTextRequest(req models.AddTextRequest) error {
	if len([]rune(req.Data)) > 1000 {
		return ErrUserTextDataIsTooBig
	}
	if len([]rune(req.Mark)) > 50 {
		return ErrUserMarkIsTooBig
	}
	if len([]rune(req.Description)) > 50 {
		return ErrUserDescriptionIsTooBig
	}

	return nil
}
