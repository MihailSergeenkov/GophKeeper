package services

import (
	"context"
	"errors"
	"testing"

	"github.com/MihailSergeenkov/GophKeeper/internal/models"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/config"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/services/mocks"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/storage"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddPassword(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	store := mocks.NewMockStorager(mockCtrl)
	fs := mocks.NewMockFileStorager(mockCtrl)
	crypter := mocks.NewMockCrypter(mockCtrl)
	settings := config.Settings{}
	s := NewServices(store, fs, crypter, &settings)

	req := models.AddPasswordRequest{
		Login:       "test",
		Password:    "test",
		Mark:        "test",
		Description: "test",
	}

	ctx := context.Background()
	dataType := "password"
	encData := []byte("some data")

	type sResponse struct {
		id  int
		err error
	}
	tests := []struct {
		name      string
		sResponse sResponse
		wantErr   bool
	}{
		{
			name: "add password success",
			sResponse: sResponse{
				id:  1,
				err: nil,
			},
			wantErr: false,
		},
		{
			name: "add password failed",
			sResponse: sResponse{
				id:  0,
				err: errors.New("some error"),
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			crypter.EXPECT().EncryptData(gomock.Any()).Times(1).Return(encData)
			store.EXPECT().
				AddUserData(ctx, encData, req.Mark, req.Description, dataType).
				Times(1).Return(test.sResponse.id, test.sResponse.err)

			id, err := s.AddPassword(ctx, req)

			if test.wantErr {
				require.Error(t, err)
				assert.ErrorContains(t, err, "failed to add user data", "some error")
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.sResponse.id, id)
			}
		})
	}
}

func TestAddPasswordValidationFailed(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	store := mocks.NewMockStorager(mockCtrl)
	fs := mocks.NewMockFileStorager(mockCtrl)
	crypter := mocks.NewMockCrypter(mockCtrl)
	settings := config.Settings{}
	s := NewServices(store, fs, crypter, &settings)

	ctx := context.Background()

	type arg struct {
		req models.AddPasswordRequest
	}
	type want struct {
		err error
	}
	tests := []struct {
		name string
		arg  arg
		want want
	}{
		{
			name: "when user login very big",
			arg: arg{
				req: models.AddPasswordRequest{
					Login: `testtesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttest
					testtesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttest
					testtesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttest`,
					Password:    "test",
					Mark:        "test",
					Description: "test",
				},
			},
			want: want{
				err: ErrUserLoginIsTooBig,
			},
		},
		{
			name: "when user password very big",
			arg: arg{
				req: models.AddPasswordRequest{
					Login: "test",
					Password: `testtesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttest
					testtesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttest
					testtesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttest`,
					Mark:        "test",
					Description: "test",
				},
			},
			want: want{
				err: ErrUserPasswordIsTooBig,
			},
		},
		{
			name: "when user mark very big",
			arg: arg{
				req: models.AddPasswordRequest{
					Login:       "test",
					Password:    "test",
					Mark:        generateString(150),
					Description: "test",
				},
			},
			want: want{
				err: ErrUserMarkIsTooBig,
			},
		},
		{
			name: "when user description very big",
			arg: arg{
				req: models.AddPasswordRequest{
					Login:       "test",
					Password:    "test",
					Mark:        "test",
					Description: generateString(4000),
				},
			},
			want: want{
				err: ErrUserDescriptionIsTooBig,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			crypter.EXPECT().EncryptData(gomock.Any()).Times(0)
			store.EXPECT().AddUserData(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

			_, err := s.AddPassword(ctx, test.arg.req)

			require.Error(t, err)
			assert.ErrorContains(t, err, test.want.err.Error())
		})
	}
}

func TestGetPassword(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	store := mocks.NewMockStorager(mockCtrl)
	fs := mocks.NewMockFileStorager(mockCtrl)
	crypter := mocks.NewMockCrypter(mockCtrl)
	settings := config.Settings{}
	s := NewServices(store, fs, crypter, &settings)

	ctx := context.Background()
	userDataID := 1
	decData := []byte("some data")
	mark := "test"
	description := "test"
	dataType := "password"

	type cResponse struct {
		jsonData []byte
		err      error
	}
	tests := []struct {
		name      string
		cResponse cResponse
		wantErr   bool
		errText   string
	}{
		{
			name: "get user data success",
			cResponse: cResponse{
				jsonData: []byte(`{"login":"test","password":"test"}`),
				err:      nil,
			},
			wantErr: false,
			errText: "",
		},
		{
			name: "when decrypt data failed",
			cResponse: cResponse{
				jsonData: nil,
				err:      errors.New("some error"),
			},
			wantErr: true,
			errText: "failed to decrypt data",
		},
		{
			name: "when generate user data failed",
			cResponse: cResponse{
				jsonData: []byte(`test`),
				err:      nil,
			},
			wantErr: true,
			errText: "failed to generate data",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store.EXPECT().GetUserData(ctx, userDataID, dataType).
				Times(1).Return(decData, mark, description, nil)

			crypter.EXPECT().DecryptData(decData).Times(1).Return(test.cResponse.jsonData, test.cResponse.err)

			resp, err := s.GetPassword(ctx, userDataID)

			if test.wantErr {
				require.Error(t, err)
				assert.ErrorContains(t, err, test.errText)
			} else {
				require.NoError(t, err)
				assert.Equal(t, models.Password{
					ID:          userDataID,
					Login:       "test",
					Password:    "test",
					Mark:        mark,
					Description: description,
				}, resp)
			}
		})
	}
}

func TestGetPasswordFailedStorage(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	store := mocks.NewMockStorager(mockCtrl)
	fs := mocks.NewMockFileStorager(mockCtrl)
	crypter := mocks.NewMockCrypter(mockCtrl)
	settings := config.Settings{}
	s := NewServices(store, fs, crypter, &settings)

	ctx := context.Background()
	userDataID := 1
	dataType := "password"

	type sResponse struct {
		err error
	}
	tests := []struct {
		name      string
		sResponse sResponse
		errText   string
	}{
		{
			name: "failed to get user data",
			sResponse: sResponse{
				err: errors.New("some error"),
			},
			errText: "failed to get user data",
		},
		{
			name: "when user data not found",
			sResponse: sResponse{
				err: storage.ErrUserDataNotFound,
			},
			errText: "requested data no found",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store.EXPECT().GetUserData(ctx, userDataID, dataType).Times(1).Return([]byte{}, "", "", test.sResponse.err)
			crypter.EXPECT().DecryptData(gomock.Any()).Times(0)

			_, err := s.GetPassword(ctx, userDataID)

			require.Error(t, err)
			assert.ErrorContains(t, err, test.errText, test.sResponse.err.Error())
		})
	}
}
