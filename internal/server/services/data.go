package services

import (
	"context"
	"fmt"

	"github.com/MihailSergeenkov/GophKeeper/internal/models"
)

// FetchUserData функция для получения базовой информации о данных пользователя.
func (s *Services) FetchUserData(ctx context.Context) ([]models.UserData, error) {
	data, err := s.storage.FetchUserData(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user data: %w", err)
	}

	return data, nil
}
