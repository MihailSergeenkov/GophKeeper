package services

import (
	"context"
	"errors"
	"testing"

	"github.com/MihailSergeenkov/GophKeeper/internal/server/config"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/services/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPing(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	storage := mocks.NewMockStorager(mockCtrl)
	fs := mocks.NewMockFileStorager(mockCtrl)
	crypter := mocks.NewMockCrypter(mockCtrl)
	settings := config.Settings{}
	s := NewServices(storage, fs, crypter, &settings)

	ctx := context.Background()

	tests := []struct {
		name    string
		wantErr bool
		err     error
	}{
		{
			name:    "success ping",
			wantErr: false,
			err:     nil,
		},
		{
			name:    "ping failed",
			wantErr: true,
			err:     errors.New("some error"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage.EXPECT().Ping(ctx).Times(1).Return(test.err)

			err := s.Ping(ctx)

			if test.wantErr {
				require.Error(t, err)
				assert.ErrorContains(t, err, "failed to ping DB", "some error")
			} else {
				require.NoError(t, err)
			}
		})
	}
}
