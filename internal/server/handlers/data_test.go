package handlers

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MihailSergeenkov/GophKeeper/internal/models"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/handlers/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestFetchUserData(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	s := mocks.NewMockServicer(mockCtrl)
	l := mocks.NewMockLogger(mockCtrl)
	handlers := NewHandlers(s, l)

	errSome := errors.New("some error")

	type serviceResponse struct {
		res []models.UserData
		err error
	}

	type want struct {
		code          int
		contentType   string
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
			name: "fetch user data success",
			serviceResponse: serviceResponse{
				res: []models.UserData{
					{
						ID:          1,
						Type:        "card",
						Mark:        "Mark",
						Description: "Description",
					},
				},
				err: nil,
			},
			want: want{
				code:          http.StatusOK,
				contentType:   JSONContentType,
				body:          "[{\"mark\":\"Mark\",\"description\":\"Description\",\"type\":\"card\",\"id\":1}]\n",
				errorLogTimes: 0,
				log:           "",
			},
		},
		{
			name: "fetch user data failed",
			serviceResponse: serviceResponse{
				res: []models.UserData{},
				err: errSome,
			},
			want: want{
				code:          http.StatusInternalServerError,
				contentType:   "",
				body:          "",
				errorLogTimes: 1,
				log:           "failed to fetch user data from DB",
			},
		},
		{
			name: "when user data not found",
			serviceResponse: serviceResponse{
				res: []models.UserData{},
				err: nil,
			},
			want: want{
				code:          http.StatusNoContent,
				contentType:   "",
				body:          "",
				errorLogTimes: 0,
				log:           "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_ = s.EXPECT().FetchUserData(gomock.Any()).Times(1).Return(test.serviceResponse.res, test.serviceResponse.err)
			_ = l.EXPECT().Error(test.want.log, zap.Error(errSome)).Times(test.want.errorLogTimes)

			request := httptest.NewRequest(http.MethodGet, "/api/user/data", http.NoBody)
			w := httptest.NewRecorder()
			handlers.FetchUserData()(w, request)

			res := w.Result()
			defer closeBody(t, res)

			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get(ContentTypeHeader))

			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, test.want.body, string(resBody))
		})
	}
}
