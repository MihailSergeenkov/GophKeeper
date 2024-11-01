package services

import (
	"errors"
	"net/http"
	"testing"

	"github.com/MihailSergeenkov/GophKeeper/internal/client/services/mocks"
	"github.com/MihailSergeenkov/GophKeeper/internal/models"
	"github.com/go-resty/resty/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cfg := mocks.NewMockConfigurer(mockCtrl)
	r := mocks.NewMockRequester(mockCtrl)
	s := Init(cfg, r)

	url := "http://some/api"
	req := models.RegisterUserRequest{}

	type postResponse struct {
		resp *resty.Response
		err  error
	}
	tests := []struct {
		name         string
		postResponse postResponse
		wantErr      bool
		errText      string
	}{
		{
			name: "register user success",
			postResponse: postResponse{
				resp: &resty.Response{
					RawResponse: &http.Response{StatusCode: http.StatusOK},
				},
				err: nil,
			},
			wantErr: false,
			errText: "",
		},
		{
			name: "register user failed when response status not 200",
			postResponse: postResponse{
				resp: &resty.Response{
					RawResponse: &http.Response{StatusCode: http.StatusForbidden},
				},
				err: nil,
			},
			wantErr: true,
			errText: "response status",
		},
		{
			name: "register user failed when request failed",
			postResponse: postResponse{
				resp: nil,
				err:  errors.New("some error"),
			},
			wantErr: true,
			errText: "failed request",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg.EXPECT().GetServerAPI().Times(1).Return(url)

			r.EXPECT().Post(url+"/user/register", gomock.Any(), gomock.Any()).
				Times(1).Return(test.postResponse.resp, test.postResponse.err)

			err := s.RegisterUser(req)

			if test.wantErr {
				require.Error(t, err)
				assert.ErrorContains(t, err, test.errText)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestLoginUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cfg := mocks.NewMockConfigurer(mockCtrl)
	r := mocks.NewMockRequester(mockCtrl)
	s := Init(cfg, r)

	url := "http://some/api"
	req := models.CreateUserTokenRequest{}

	type postResponse struct {
		resp *resty.Response
		err  error
	}
	type updateToken struct {
		count int
		err   error
	}
	tests := []struct {
		name         string
		postResponse postResponse
		updateToken  updateToken
		wantErr      bool
		errText      string
	}{
		{
			name: "login user success",
			postResponse: postResponse{
				resp: &resty.Response{
					RawResponse: &http.Response{StatusCode: http.StatusOK},
				},
				err: nil,
			},
			updateToken: updateToken{
				count: 1,
				err:   nil,
			},
			wantErr: false,
			errText: "",
		},
		{
			name: "login user failed",
			postResponse: postResponse{
				resp: &resty.Response{
					RawResponse: &http.Response{StatusCode: http.StatusOK},
				},
				err: nil,
			},
			updateToken: updateToken{
				count: 1,
				err:   errors.New("some error"),
			},
			wantErr: true,
			errText: "failed to update auth token",
		},
		{
			name: "login user failed when response status not 200",
			postResponse: postResponse{
				resp: &resty.Response{
					RawResponse: &http.Response{StatusCode: http.StatusForbidden},
				},
				err: nil,
			},
			updateToken: updateToken{
				count: 0,
				err:   nil,
			},
			wantErr: true,
			errText: "response status",
		},
		{
			name: "login user failed when request failed",
			postResponse: postResponse{
				resp: nil,
				err:  errors.New("some error"),
			},
			updateToken: updateToken{
				count: 0,
				err:   nil,
			},
			wantErr: true,
			errText: "failed request",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg.EXPECT().GetServerAPI().Times(1).Return(url)

			r.EXPECT().Post(url+"/user/token", gomock.Any(), gomock.Any(), gomock.Any()).
				Times(1).Return(test.postResponse.resp, test.postResponse.err)

			cfg.EXPECT().UpdateToken(gomock.Any()).Times(test.updateToken.count).Return(test.updateToken.err)

			err := s.LoginUser(req)

			if test.wantErr {
				require.Error(t, err)
				assert.ErrorContains(t, err, test.errText)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestLogoutUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cfg := mocks.NewMockConfigurer(mockCtrl)
	r := mocks.NewMockRequester(mockCtrl)
	s := Init(cfg, r)

	type updateToken struct {
		err error
	}
	type updateData struct {
		count int
		err   error
	}
	tests := []struct {
		name        string
		updateToken updateToken
		updateData  updateData
		wantErr     bool
		errText     string
	}{
		{
			name: "loguot user success",
			updateToken: updateToken{
				err: nil,
			},
			updateData: updateData{
				count: 1,
				err:   nil,
			},
			wantErr: false,
			errText: "",
		},
		{
			name: "logout user failed when failed update token",
			updateToken: updateToken{
				err: errors.New("some error"),
			},
			updateData: updateData{
				count: 0,
				err:   nil,
			},
			wantErr: true,
			errText: "failed to update auth token",
		},
		{
			name: "logout user failed when failed update data",
			updateToken: updateToken{
				err: nil,
			},
			updateData: updateData{
				count: 1,
				err:   errors.New("some error"),
			},
			wantErr: true,
			errText: "ailed to update data",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg.EXPECT().UpdateToken("").Times(1).Return(test.updateToken.err)
			cfg.EXPECT().UpdateData([]models.UserData{}).Times(test.updateData.count).Return(test.updateData.err)

			err := s.LogoutUser()

			if test.wantErr {
				require.Error(t, err)
				assert.ErrorContains(t, err, test.errText)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
