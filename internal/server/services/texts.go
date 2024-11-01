package services

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/MihailSergeenkov/GophKeeper/internal/models"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/storage"
)

const textDataType = "text"

var (
	ErrUserTextDataIsTooBig = errors.New("user text data is too big")

	textDataMaxSize = 1000
)

// AddText функция для добавления текста пользователя.
func (s *Services) AddText(ctx context.Context, req models.AddTextRequest) (int, error) {
	if err := validateAddTextRequest(req); err != nil {
		return 0, failedValidateFields(err)
	}

	data := models.EncryptTextData{
		Data: req.Data,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return 0, failedGenerateJSONData(err)
	}

	encData := s.crypter.EncryptData(jsonData)

	id, err := s.storage.AddUserData(ctx, encData, req.Mark, req.Description, textDataType)
	if err != nil {
		return 0, failedAddUserData(err)
	}

	return id, nil
}

// GetText функция для получения текста пользователя.
func (s *Services) GetText(ctx context.Context, id int) (models.Text, error) {
	resp := models.Text{}

	decData, mark, description, err := s.storage.GetUserData(ctx, id, textDataType)
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

	var encData models.EncryptTextData

	if err = json.Unmarshal(jsonData, &encData); err != nil {
		return resp, failedGenerateData(err)
	}

	resp.ID = id
	resp.Data = encData.Data
	resp.Mark = mark
	resp.Description = description

	return resp, nil
}

func validateAddTextRequest(req models.AddTextRequest) error {
	if len([]rune(req.Data)) > textDataMaxSize {
		return ErrUserTextDataIsTooBig
	}
	if len([]rune(req.Mark)) > maxMarkSize {
		return ErrUserMarkIsTooBig
	}
	if len([]rune(req.Description)) > maxDescriptionSize {
		return ErrUserDescriptionIsTooBig
	}

	return nil
}
