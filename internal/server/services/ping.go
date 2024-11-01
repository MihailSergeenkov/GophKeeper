package services

import (
	"context"
	"fmt"
)

// Ping функция для проверки работоспособности БД.
func (s *Services) Ping(ctx context.Context) error {
	if err := s.storage.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping DB %w", err)
	}

	return nil
}
