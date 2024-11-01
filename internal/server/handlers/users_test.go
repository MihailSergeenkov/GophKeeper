package handlers

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/MihailSergeenkov/GophKeeper/internal/models"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/handlers/mocks"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/services"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestRegisterUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	s := mocks.NewMockServicer(mockCtrl)
	l := mocks.NewMockLogger(mockCtrl)
	handlers := NewHandlers(s, l)

	requestBody := `{"login":"test","password":"test"}`
	requestObject := models.RegisterUserRequest{
		Login:    "test",
		Password: "test",
	}

	type serviceResponse struct {
		err error
	}

	type want struct {
		code          int
		errorLogTimes int
		log           string
	}

	tests := []struct {
		name            string
		serviceResponse serviceResponse
		want            want
	}{
		{
			name: "register user success",
			serviceResponse: serviceResponse{
				err: nil,
			},
			want: want{
				code:          http.StatusOK,
				errorLogTimes: 0,
				log:           "",
			},
		},
		{
			name: "register user failed with ErrUserValidationFields",
			serviceResponse: serviceResponse{
				err: services.ErrUserValidationFields,
			},
			want: want{
				code:          http.StatusBadRequest,
				errorLogTimes: 0,
				log:           "",
			},
		},
		{
			name: "register user failed with ErrUserLoginExist",
			serviceResponse: serviceResponse{
				err: services.ErrUserLoginExist,
			},
			want: want{
				code:          http.StatusConflict,
				errorLogTimes: 0,
				log:           "",
			},
		},
		{
			name: "register user failed with some error",
			serviceResponse: serviceResponse{
				err: errors.New("some error"),
			},
			want: want{
				code:          http.StatusInternalServerError,
				errorLogTimes: 1,
				log:           "failed to register user",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s.EXPECT().RegisterUser(gomock.Any(), requestObject).Times(1).Return(test.serviceResponse.err)
			l.EXPECT().Error(test.want.log, zap.Error(test.serviceResponse.err)).Times(test.want.errorLogTimes)

			request := httptest.NewRequest(http.MethodPost, "/api/user/register", strings.NewReader(requestBody))
			w := httptest.NewRecorder()
			handlers.RegisterUser()(w, request)

			res := w.Result()
			defer closeBody(t, res)

			assert.Equal(t, test.want.code, res.StatusCode)
		})
	}
}

func TestFailedReadBodyRegisterUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	s := mocks.NewMockServicer(mockCtrl)
	l := mocks.NewMockLogger(mockCtrl)
	handlers := NewHandlers(s, l)

	requestBody := `{"login":"test","password":"test",adasd}`

	t.Run("failed to read request body", func(t *testing.T) {
		s.EXPECT().RegisterUser(gomock.Any(), gomock.Any()).Times(0)
		l.EXPECT().Error("failed to read request body", gomock.Any()).Times(1)

		request := httptest.NewRequest(http.MethodPost, "/api/user/register", strings.NewReader(requestBody))
		w := httptest.NewRecorder()
		handlers.RegisterUser()(w, request)

		res := w.Result()
		defer closeBody(t, res)

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
}

func TestCreateUserToken(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	s := mocks.NewMockServicer(mockCtrl)
	l := mocks.NewMockLogger(mockCtrl)
	handlers := NewHandlers(s, l)

	requestBody := `{"login":"test","password":"test"}`
	requestObject := models.CreateUserTokenRequest{
		Login:    "test",
		Password: "test",
	}

	type serviceResponse struct {
		res models.CreateUserTokenResponse
		err error
	}

	type want struct {
		code          int
		body          string
		errorLogTimes int
		log           string
	}

	tests := []struct {
		name            string
		serviceResponse serviceResponse
		want            want
	}{
		{
			name: "create user token success",
			serviceResponse: serviceResponse{
				res: models.CreateUserTokenResponse{
					AuthToken: "qwerty",
				},
				err: nil,
			},
			want: want{
				code:          http.StatusOK,
				body:          "{\"auth_token\":\"qwerty\"}\n",
				errorLogTimes: 0,
				log:           "",
			},
		},
		{
			name: "create user token failed with ErrUserLoginCreds",
			serviceResponse: serviceResponse{
				res: models.CreateUserTokenResponse{},
				err: services.ErrUserLoginCreds,
			},
			want: want{
				code:          http.StatusUnauthorized,
				body:          "",
				errorLogTimes: 0,
				log:           "",
			},
		},
		{
			name: "create user token failed with some error",
			serviceResponse: serviceResponse{
				res: models.CreateUserTokenResponse{},
				err: errors.New("some error"),
			},
			want: want{
				code:          http.StatusInternalServerError,
				body:          "",
				errorLogTimes: 1,
				log:           "failed to create user token",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s.EXPECT().
				CreateUserToken(gomock.Any(), requestObject).
				Times(1).
				Return(test.serviceResponse.res, test.serviceResponse.err)

			l.EXPECT().Error(test.want.log, zap.Error(test.serviceResponse.err)).Times(test.want.errorLogTimes)

			request := httptest.NewRequest(http.MethodPost, "/api/user/token", strings.NewReader(requestBody))
			w := httptest.NewRecorder()
			handlers.CreateUserToken()(w, request)

			res := w.Result()
			defer closeBody(t, res)

			assert.Equal(t, test.want.code, res.StatusCode)

			if http.StatusOK == res.StatusCode {
				resBody, err := io.ReadAll(res.Body)
				require.NoError(t, err)
				assert.Equal(t, test.want.body, string(resBody))
			}
		})
	}
}

func TestFailedReadBodyLoginUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	s := mocks.NewMockServicer(mockCtrl)
	l := mocks.NewMockLogger(mockCtrl)
	handlers := NewHandlers(s, l)

	requestBody := `{"login":"test","password":"test",adasd}`

	t.Run("failed to read request body", func(t *testing.T) {
		s.EXPECT().CreateUserToken(gomock.Any(), gomock.Any()).Times(0)
		l.EXPECT().Error("failed to read request body", gomock.Any()).Times(1)

		request := httptest.NewRequest(http.MethodPost, "/api/user/token", strings.NewReader(requestBody))
		w := httptest.NewRecorder()
		handlers.CreateUserToken()(w, request)

		res := w.Result()
		defer closeBody(t, res)

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
}
