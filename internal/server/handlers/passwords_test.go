package handlers

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/MihailSergeenkov/GophKeeper/internal/models"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/config"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/handlers/mocks"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/routes"
	rMocks "github.com/MihailSergeenkov/GophKeeper/internal/server/routes/mocks"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/services"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestAddPassword(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	s := mocks.NewMockServicer(mockCtrl)
	l := mocks.NewMockLogger(mockCtrl)
	handlers := NewHandlers(s, l)

	requestBody := `{"login":"test","password":"test","mark":"test","description":"test"}`
	requestObject := models.AddPasswordRequest{
		Login:       "test",
		Password:    "test",
		Mark:        "test",
		Description: "test",
	}

	type serviceResponse struct {
		id  int
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
			name: "add password success",
			serviceResponse: serviceResponse{
				id:  1,
				err: nil,
			},
			want: want{
				code:          http.StatusCreated,
				body:          "{\"id\":1}\n",
				errorLogTimes: 0,
				log:           "",
			},
		},
		{
			name: "add password failed",
			serviceResponse: serviceResponse{
				id:  0,
				err: errors.New("some error"),
			},
			want: want{
				code:          http.StatusInternalServerError,
				body:          "",
				errorLogTimes: 1,
				log:           "failed to add password",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s.EXPECT().
				AddPassword(gomock.Any(), requestObject).
				Times(1).
				Return(test.serviceResponse.id, test.serviceResponse.err)

			l.EXPECT().Error(test.want.log, zap.Error(test.serviceResponse.err)).Times(test.want.errorLogTimes)

			request := httptest.NewRequest(http.MethodPost, "/api/user/passwords", strings.NewReader(requestBody))
			w := httptest.NewRecorder()
			handlers.AddPassword()(w, request)

			res := w.Result()
			defer closeBody(t, res)

			assert.Equal(t, test.want.code, res.StatusCode)

			if http.StatusCreated == res.StatusCode {
				resBody, err := io.ReadAll(res.Body)
				require.NoError(t, err)
				assert.Equal(t, test.want.body, string(resBody))
			}
		})
	}
}

func TestAddPasswordFailedReadBody(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	s := mocks.NewMockServicer(mockCtrl)
	l := mocks.NewMockLogger(mockCtrl)
	handlers := NewHandlers(s, l)

	requestBody := `{"login":"test","password":"test","mark":"test","description":"test",adasd}`

	t.Run("failed to read request body", func(t *testing.T) {
		s.EXPECT().AddPassword(gomock.Any(), gomock.Any()).Times(0)
		l.EXPECT().Error("failed to read request body", gomock.Any()).Times(1)

		request := httptest.NewRequest(http.MethodPost, "/api/user/passwords", strings.NewReader(requestBody))
		w := httptest.NewRecorder()
		handlers.AddPassword()(w, request)

		res := w.Result()
		defer closeBody(t, res)

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
}

func TestGetPassword(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	s := mocks.NewMockServicer(mockCtrl)
	l := mocks.NewMockLogger(mockCtrl)
	handlers := NewHandlers(s, l)

	settings, err := config.Setup(false)
	require.NoError(t, err)
	storage := rMocks.NewMockStorager(mockCtrl)

	passwordID := 1

	r := routes.NewRouter(handlers, settings, zap.NewNop(), storage)
	ts := httptest.NewServer(r)
	defer ts.Close()

	type serviceResponse struct {
		res models.Password
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
			name: "get password success",
			serviceResponse: serviceResponse{
				res: models.Password{
					ID:          1,
					Login:       "test",
					Password:    "test",
					Mark:        "test",
					Description: "test",
				},
				err: nil,
			},
			want: want{
				code:          http.StatusOK,
				body:          `{"id":1,"login":"test","password":"test","mark":"test","description":"test"}` + "\n",
				errorLogTimes: 0,
				log:           "",
			},
		},
		{
			name: "password no found",
			serviceResponse: serviceResponse{
				res: models.Password{},
				err: services.ErrNotFound,
			},
			want: want{
				code:          http.StatusNotFound,
				body:          "",
				errorLogTimes: 0,
				log:           "",
			},
		},
		{
			name: "get password failed with some error",
			serviceResponse: serviceResponse{
				res: models.Password{},
				err: errors.New("some error"),
			},
			want: want{
				code:          http.StatusInternalServerError,
				body:          "",
				errorLogTimes: 1,
				log:           "failed to get password",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s.EXPECT().GetPassword(gomock.Any(), passwordID).Times(1).
				Return(test.serviceResponse.res, test.serviceResponse.err)

			l.EXPECT().Error(test.want.log, zap.Error(test.serviceResponse.err)).Times(test.want.errorLogTimes)
			storage.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Times(1)

			res, resBody := testRequest(t, ts, http.MethodGet, "/api/user/passwords/1")
			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.body, resBody)
		})
	}
}

func TestGetPasswordFailedReadParam(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	s := mocks.NewMockServicer(mockCtrl)
	l := mocks.NewMockLogger(mockCtrl)
	handlers := NewHandlers(s, l)

	settings, err := config.Setup(false)
	require.NoError(t, err)
	storage := rMocks.NewMockStorager(mockCtrl)

	r := routes.NewRouter(handlers, settings, zap.NewNop(), storage)
	ts := httptest.NewServer(r)
	defer ts.Close()

	t.Run("failed to read request param", func(t *testing.T) {
		s.EXPECT().GetPassword(gomock.Any(), gomock.Any()).Times(0)
		l.EXPECT().Error("failed ID param", gomock.Any()).Times(1)
		storage.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Times(1)

		res, _ := testRequest(t, ts, http.MethodGet, "/api/user/passwords/adasd")
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
}
