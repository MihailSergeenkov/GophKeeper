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

func TestAddCard(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	s := mocks.NewMockServicer(mockCtrl)
	l := mocks.NewMockLogger(mockCtrl)
	handlers := NewHandlers(s, l)

	requestBody := `
	{"number":"1234123412341234","owner":"test","expiry_date":"11/2300","cvv2":"777","mark":"test","description":"test"}
	`
	requestObject := models.AddCardRequest{
		Number:      "1234123412341234",
		Owner:       "test",
		ExpiryDate:  "11/2300",
		CVV2:        "777",
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
			name: "add card success",
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
			name: "add card failed",
			serviceResponse: serviceResponse{
				id:  0,
				err: errors.New("some error"),
			},
			want: want{
				code:          http.StatusInternalServerError,
				body:          "",
				errorLogTimes: 1,
				log:           "failed to add card",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s.EXPECT().
				AddCard(gomock.Any(), &requestObject).
				Times(1).
				Return(test.serviceResponse.id, test.serviceResponse.err)

			l.EXPECT().Error(test.want.log, zap.Error(test.serviceResponse.err)).Times(test.want.errorLogTimes)

			request := httptest.NewRequest(http.MethodPost, "/api/user/cards", strings.NewReader(requestBody))
			w := httptest.NewRecorder()
			handlers.AddCard()(w, request)

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

func TestAddCardFailedReadBody(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	s := mocks.NewMockServicer(mockCtrl)
	l := mocks.NewMockLogger(mockCtrl)
	handlers := NewHandlers(s, l)

	requestBody := `{"number":"1234123412341234","mark":"test","description":"test",adasd}`

	t.Run("failed to read request body", func(t *testing.T) {
		s.EXPECT().AddCard(gomock.Any(), gomock.Any()).Times(0)
		l.EXPECT().Error("failed to read request body", gomock.Any()).Times(1)

		request := httptest.NewRequest(http.MethodPost, "/api/user/cards", strings.NewReader(requestBody))
		w := httptest.NewRecorder()
		handlers.AddCard()(w, request)

		res := w.Result()
		defer closeBody(t, res)

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
}

func TestGetCard(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	s := mocks.NewMockServicer(mockCtrl)
	l := mocks.NewMockLogger(mockCtrl)
	handlers := NewHandlers(s, l)

	settings, err := config.Setup(false)
	require.NoError(t, err)
	storage := rMocks.NewMockStorager(mockCtrl)

	cardID := 1

	r := routes.NewRouter(handlers, settings, zap.NewNop(), storage)
	ts := httptest.NewServer(r)
	defer ts.Close()

	type serviceResponse struct {
		res models.Card
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
			name: "get card success",
			serviceResponse: serviceResponse{
				res: models.Card{
					ID:          1,
					Number:      "1234123412341234",
					Owner:       "test",
					ExpiryDate:  "11/2300",
					CVV2:        "777",
					Mark:        "test",
					Description: "test",
				},
				err: nil,
			},
			want: want{
				code:          http.StatusOK,
				body:          `{"number":"1234123412341234","owner":"test","expiry_date":"11/2300","cvv2":"777","mark":"test","description":"test","id":1}` + "\n", //nolint:lll // Исключение
				errorLogTimes: 0,
				log:           "",
			},
		},
		{
			name: "card no found",
			serviceResponse: serviceResponse{
				res: models.Card{},
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
			name: "get card failed with some error",
			serviceResponse: serviceResponse{
				res: models.Card{},
				err: errors.New("some error"),
			},
			want: want{
				code:          http.StatusInternalServerError,
				body:          "",
				errorLogTimes: 1,
				log:           "failed to get card",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s.EXPECT().GetCard(gomock.Any(), cardID).Times(1).
				Return(test.serviceResponse.res, test.serviceResponse.err)

			l.EXPECT().Error(test.want.log, zap.Error(test.serviceResponse.err)).Times(test.want.errorLogTimes)
			storage.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Times(1)

			res, resBody := testGetRequest(t, ts, "/api/user/cards/1")
			closeBody(t, res)

			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.body, resBody)
		})
	}
}

func TestGetCardFailedReadParam(t *testing.T) {
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
		s.EXPECT().GetCard(gomock.Any(), gomock.Any()).Times(0)
		l.EXPECT().Error("failed card ID param", gomock.Any()).Times(1)
		storage.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Times(1)

		res, _ := testGetRequest(t, ts, "/api/user/cards/adasd")
		closeBody(t, res)

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
}
