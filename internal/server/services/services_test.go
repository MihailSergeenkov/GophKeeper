package services

import (
	"testing"

	"github.com/MihailSergeenkov/GophKeeper/internal/server/config"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/services/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewServices(t *testing.T) {
	t.Run("init services", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		storage := mocks.NewMockStorager(mockCtrl)
		fs := mocks.NewMockFileStorager(mockCtrl)
		crypter := mocks.NewMockCrypter(mockCtrl)
		settings, err := config.Setup(false)
		require.NoError(t, err)

		s := NewServices(storage, fs, crypter, settings)
		assert.Equal(t, storage, s.storage)
		assert.Equal(t, settings, s.settings)
	})
}
