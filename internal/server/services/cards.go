package services

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/MihailSergeenkov/GophKeeper/internal/models"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/storage"
)

const cardDataType = "card"

var (
	ErrUserNumberInvalid     = errors.New("user number invalid")
	ErrUserOwnerIsTooBig     = errors.New("user owner is too big")
	ErrUserExpiryDateInvalid = errors.New("user expiry date invalid")
	ErrUserCVV2Invalid       = errors.New("user cvv2 invalid")

	cardNumberSize     = 16
	cardOwnerMaxSize   = 100
	cardExpiryDateSize = 7
	cardCVV2Size       = 3
)

// AddCard функция для добавления карты пользователя.
func (s *Services) AddCard(ctx context.Context, req *models.AddCardRequest) (int, error) {
	if err := validateAddCardRequest(req); err != nil {
		return 0, failedValidateFields(err)
	}

	data := models.EncryptCardData{
		Number:     req.Number,
		Owner:      req.Owner,
		ExpiryDate: req.ExpiryDate,
		CVV2:       req.CVV2,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return 0, failedGenerateJSONData(err)
	}

	encData := s.crypter.EncryptData(jsonData)

	id, err := s.storage.AddUserData(ctx, encData, req.Mark, req.Description, cardDataType)
	if err != nil {
		return 0, failedAddUserData(err)
	}

	return id, nil
}

// GetCard функция для получения карты пользователя.
func (s *Services) GetCard(ctx context.Context, id int) (models.Card, error) {
	resp := models.Card{}

	decData, mark, description, err := s.storage.GetUserData(ctx, id, cardDataType)
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

	var encData models.EncryptCardData

	if err = json.Unmarshal(jsonData, &encData); err != nil {
		return resp, failedGenerateData(err)
	}

	resp.ID = id
	resp.Number = encData.Number
	resp.Owner = encData.Owner
	resp.ExpiryDate = encData.ExpiryDate
	resp.CVV2 = encData.CVV2
	resp.Mark = mark
	resp.Description = description

	return resp, nil
}

func validateAddCardRequest(req *models.AddCardRequest) error {
	if len(req.Number) != cardNumberSize {
		return ErrUserNumberInvalid
	}
	if len(req.Owner) > cardOwnerMaxSize {
		return ErrUserOwnerIsTooBig
	}
	if len(req.ExpiryDate) != cardExpiryDateSize {
		return ErrUserExpiryDateInvalid
	}
	if len(req.CVV2) != cardCVV2Size {
		return ErrUserCVV2Invalid
	}
	if len([]rune(req.Mark)) > maxMarkSize {
		return ErrUserMarkIsTooBig
	}
	if len([]rune(req.Description)) > maxDescriptionSize {
		return ErrUserDescriptionIsTooBig
	}

	return nil
}
