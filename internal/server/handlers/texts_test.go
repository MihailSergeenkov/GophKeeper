package handlers

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/MihailSergeenkov/GophKeeper/internal/server/handlers/mocks"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/models"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/services"
	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestAddText(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	s := mocks.NewMockServicer(mockCtrl)
	l := mocks.NewMockLogger(mockCtrl)
	handlers := NewHandlers(s, l)

	requestBody := `{"data":"test","mark":"test","description":"test"}`
	requestObject := models.AddTextRequest{
		Data:        "test",
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
			name: "add text success",
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
			name: "add text failed",
			serviceResponse: serviceResponse{
				id:  0,
				err: errors.New("some error"),
			},
			want: want{
				code:          http.StatusInternalServerError,
				body:          "",
				errorLogTimes: 1,
				log:           "failed to add text",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s.EXPECT().
				AddText(gomock.Any(), requestObject).
				Times(1).
				Return(test.serviceResponse.id, test.serviceResponse.err)

			l.EXPECT().Error(test.want.log, zap.Error(test.serviceResponse.err)).Times(test.want.errorLogTimes)

			request := httptest.NewRequest(http.MethodPost, "/api/user/texts", strings.NewReader(requestBody))
			w := httptest.NewRecorder()
			handlers.AddText()(w, request)

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

func TestAddTextFailedReadBody(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	s := mocks.NewMockServicer(mockCtrl)
	l := mocks.NewMockLogger(mockCtrl)
	handlers := NewHandlers(s, l)

	requestBody := `{"data":"test","mark":"test","description":"test",adasd}`

	t.Run("failed to read request body", func(t *testing.T) {
		s.EXPECT().AddText(gomock.Any(), gomock.Any()).Times(0)
		l.EXPECT().Error("failed to read request body", gomock.Any()).Times(1)

		request := httptest.NewRequest(http.MethodPost, "/api/user/texts", strings.NewReader(requestBody))
		w := httptest.NewRecorder()
		handlers.AddText()(w, request)

		res := w.Result()
		defer closeBody(t, res)

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
}

func TestGetText(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	s := mocks.NewMockServicer(mockCtrl)
	l := mocks.NewMockLogger(mockCtrl)
	handlers := NewHandlers(s, l)

	textID := 1

	type serviceResponse struct {
		res models.Text
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
			name: "get text success",
			serviceResponse: serviceResponse{
				res: models.Text{
					ID:          1,
					Data:        "test",
					Mark:        "test",
					Description: "test",
				},
				err: nil,
			},
			want: want{
				code:          http.StatusOK,
				body:          "{\"id\":1,\"data\":\"test\",\"mark\":\"test\",\"description\":\"test\"}\n",
				errorLogTimes: 0,
				log:           "",
			},
		},
		{
			name: "text no found",
			serviceResponse: serviceResponse{
				res: models.Text{},
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
			name: "get text failed with some error",
			serviceResponse: serviceResponse{
				res: models.Text{},
				err: errors.New("some error"),
			},
			want: want{
				code:          http.StatusInternalServerError,
				body:          "",
				errorLogTimes: 1,
				log:           "failed to get text",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s.EXPECT().GetText(gomock.Any(), textID).Times(1).
				Return(test.serviceResponse.res, test.serviceResponse.err)

			l.EXPECT().Error(test.want.log, zap.Error(test.serviceResponse.err)).Times(test.want.errorLogTimes)

			// r := chi.NewRouter()
			// r.Get("/api/user/texts/{textID}", handlers.GetText())
			// ts := httptest.NewServer(r)
			// defer ts.Close()

			request := httptest.NewRequest(http.MethodGet, "/api/user/texts/1", http.NoBody)
			// request, _ := http.NewRequest(http.MethodGet, ts.URL+"/api/user/texts/1", http.NoBody)
			w := httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("textID", "1")
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

			// ctx := context.Background()
			// ctx = context.WithValue(rctx, "sdf", "df")

			// request := httptest.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/api/user/texts/1", http.NoBody)

			handlers.GetText()(w, request)

			res := w.Result()
			defer closeBody(t, res)
			// res, err := http.DefaultClient.Do(request)
			// require.NoError(t, err)

			assert.Equal(t, test.want.code, res.StatusCode)

			if http.StatusOK == res.StatusCode {
				resBody, err := io.ReadAll(res.Body)
				require.NoError(t, err)
				assert.Equal(t, test.want.body, string(resBody))
			}
		})
	}
}

func TestGetTextFailedReadParam(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	s := mocks.NewMockServicer(mockCtrl)
	l := mocks.NewMockLogger(mockCtrl)
	handlers := NewHandlers(s, l)

	t.Run("failed to read request param", func(t *testing.T) {
		s.EXPECT().GetText(gomock.Any(), gomock.Any()).Times(0)
		l.EXPECT().Error("failed ID param", gomock.Any()).Times(1)

		request := httptest.NewRequest(http.MethodGet, "/api/user/texts", http.NoBody)
		w := httptest.NewRecorder()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("textID", "adasd")
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

		handlers.GetText()(w, request)

		res := w.Result()
		defer closeBody(t, res)

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
}
