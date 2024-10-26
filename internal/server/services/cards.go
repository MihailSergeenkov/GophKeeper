package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/MihailSergeenkov/GophKeeper/internal/models"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/storage"
)

const cardDataType = "card"

var (
	ErrUserNumberInvalid     = errors.New("user number invalid")
	ErrUserOwnerIsTooBig     = errors.New("user owner is too big")
	ErrUserExpiryDateInvalid = errors.New("user expiry date invalid")
	ErrUserCVV2Invalid       = errors.New("user cvv2 invalid")
)

// AddCard функция для добавления карты пользователя.
func (s *Services) AddCard(ctx context.Context, req models.AddCardRequest) (int, error) {
	if err := validateAddCardRequest(req); err != nil {
		return 0, fmt.Errorf("failed to validate fields %w", err)
	}

	data := models.EncryptCardData{
		Number:     req.Number,
		Owner:      req.Owner,
		ExpiryDate: req.ExpiryDate,
		CVV2:       req.CVV2,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return 0, fmt.Errorf("failed to generate json data %w", err)
	}

	encData := s.crypter.EncryptData(jsonData)

	id, err := s.storage.AddUserData(ctx, encData, req.Mark, req.Description, cardDataType)
	if err != nil {
		return 0, fmt.Errorf("failed to add user data %w", err)
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

		return resp, fmt.Errorf("failed to get user data %w", err)
	}

	jsonData, err := s.crypter.DecryptData(decData)
	if err != nil {
		return resp, fmt.Errorf("failed to decrypt data %w", err)
	}

	var encData models.EncryptCardData

	if err = json.Unmarshal(jsonData, &encData); err != nil {
		return resp, fmt.Errorf("failed to generate data %w", err)
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

func validateAddCardRequest(req models.AddCardRequest) error {
	if len(req.Number) != 16 {
		return ErrUserNumberInvalid
	}
	if len(req.Owner) > 100 {
		return ErrUserOwnerIsTooBig
	}
	if len(req.ExpiryDate) != 7 {
		return ErrUserExpiryDateInvalid
	}
	if len(req.CVV2) != 3 {
		return ErrUserCVV2Invalid
	}
	if len([]rune(req.Mark)) > 50 {
		return ErrUserMarkIsTooBig
	}
	if len([]rune(req.Description)) > 50 {
		return ErrUserDescriptionIsTooBig
	}

	return nil
}
