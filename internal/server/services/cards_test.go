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

func TestAddCard(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	store := mocks.NewMockStorager(mockCtrl)
	fs := mocks.NewMockFileStorager(mockCtrl)
	crypter := mocks.NewMockCrypter(mockCtrl)
	settings := config.Settings{}
	s := NewServices(store, fs, crypter, &settings)

	req := models.AddCardRequest{
		Number:      "1234123412341234",
		Owner:       "test",
		ExpiryDate:  "11/2300",
		CVV2:        "777",
		Mark:        "test",
		Description: "test",
	}

	ctx := context.Background()
	dataType := "card"
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
			name: "add card success",
			sResponse: sResponse{
				id:  1,
				err: nil,
			},
			wantErr: false,
		},
		{
			name: "add card failed",
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

			id, err := s.AddCard(ctx, &req)

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

func TestAddCardValidationFailed(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	store := mocks.NewMockStorager(mockCtrl)
	fs := mocks.NewMockFileStorager(mockCtrl)
	crypter := mocks.NewMockCrypter(mockCtrl)
	settings := config.Settings{}
	s := NewServices(store, fs, crypter, &settings)

	ctx := context.Background()

	type arg struct {
		req models.AddCardRequest
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
			name: "when user number invalid",
			arg: arg{
				req: models.AddCardRequest{
					Number:      "test",
					Owner:       "test",
					ExpiryDate:  "11/2300",
					CVV2:        "777",
					Mark:        "test",
					Description: "test",
				},
			},
			want: want{
				err: ErrUserNumberInvalid,
			},
		},
		{
			name: "when user owner very big",
			arg: arg{
				req: models.AddCardRequest{
					Number: "1234123412341234",
					Owner: `testtesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttest
					testtesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttest
					testtesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttest`,
					ExpiryDate:  "11/2300",
					CVV2:        "777",
					Mark:        "test",
					Description: "test",
				},
			},
			want: want{
				err: ErrUserOwnerIsTooBig,
			},
		},
		{
			name: "when user expiry date invalid",
			arg: arg{
				req: models.AddCardRequest{
					Number:      "1234123412341234",
					Owner:       "test",
					ExpiryDate:  "test",
					CVV2:        "777",
					Mark:        "test",
					Description: "test",
				},
			},
			want: want{
				err: ErrUserExpiryDateInvalid,
			},
		},
		{
			name: "when user cvv2 invalid",
			arg: arg{
				req: models.AddCardRequest{
					Number:      "1234123412341234",
					Owner:       "test",
					ExpiryDate:  "11/2300",
					CVV2:        "test",
					Mark:        "test",
					Description: "test",
				},
			},
			want: want{
				err: ErrUserCVV2Invalid,
			},
		},
		{
			name: "when user mark very big",
			arg: arg{
				req: models.AddCardRequest{
					Number:      "1234123412341234",
					Owner:       "test",
					ExpiryDate:  "11/2300",
					CVV2:        "777",
					Mark:        "testtesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttest",
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
				req: models.AddCardRequest{
					Number:      "1234123412341234",
					Owner:       "test",
					ExpiryDate:  "11/2300",
					CVV2:        "777",
					Mark:        "test",
					Description: "testtesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttest",
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

			_, err := s.AddCard(ctx, &test.arg.req)

			require.Error(t, err)
			assert.ErrorContains(t, err, test.want.err.Error())
		})
	}
}

func TestGetCard(t *testing.T) {
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
	dataType := "card"

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
				jsonData: []byte(`{"number":"1234123412341234","owner":"test","expiry_date":"11/2300","cvv2":"777"}`),
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

			resp, err := s.GetCard(ctx, userDataID)

			if test.wantErr {
				require.Error(t, err)
				assert.ErrorContains(t, err, test.errText)
			} else {
				require.NoError(t, err)
				assert.Equal(t, models.Card{
					ID:          userDataID,
					Number:      "1234123412341234",
					Owner:       "test",
					ExpiryDate:  "11/2300",
					CVV2:        "777",
					Mark:        mark,
					Description: description,
				}, resp)
			}
		})
	}
}

func TestGetCardFailedStorage(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	store := mocks.NewMockStorager(mockCtrl)
	fs := mocks.NewMockFileStorager(mockCtrl)
	crypter := mocks.NewMockCrypter(mockCtrl)
	settings := config.Settings{}
	s := NewServices(store, fs, crypter, &settings)

	ctx := context.Background()
	userDataID := 1
	dataType := "card"

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

			_, err := s.GetCard(ctx, userDataID)

			require.Error(t, err)
			assert.ErrorContains(t, err, test.errText, test.sResponse.err.Error())
		})
	}
}
