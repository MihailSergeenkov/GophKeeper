package services

import (
	"context"
	"errors"
	"testing"

	"github.com/MihailSergeenkov/GophKeeper/internal/server/config"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/models"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/services/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetchUserData(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	storage := mocks.NewMockStorager(mockCtrl)
	fs := mocks.NewMockFileStorager(mockCtrl)
	crypter := mocks.NewMockCrypter(mockCtrl)
	settings := config.Settings{}
	s := NewServices(storage, fs, crypter, &settings)

	ctx := context.Background()

	type storageResponse struct {
		data []models.UserData
		err  error
	}
	tests := []struct {
		name            string
		wantErr         bool
		storageResponse storageResponse
	}{
		{
			name:    "success fetch user data",
			wantErr: false,
			storageResponse: storageResponse{
				data: []models.UserData{
					{
						ID:          1,
						Type:        "card",
						Mark:        "Mark",
						Description: "Description",
					},
				},
				err: nil,
			},
		},
		{
			name:    "fetch user data failed",
			wantErr: true,
			storageResponse: storageResponse{
				data: []models.UserData{},
				err:  errors.New("some error"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage.EXPECT().FetchUserData(ctx).Times(1).Return(test.storageResponse.data, test.storageResponse.err)

			resp, err := s.FetchUserData(ctx)

			if test.wantErr {
				require.Error(t, err)
				assert.ErrorContains(t, err, "failed to fetch user data", "some error")
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.storageResponse.data, resp)
			}
		})
	}
}
