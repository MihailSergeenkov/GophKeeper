package routes

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MihailSergeenkov/GophKeeper/internal/models"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/config"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/routes/mocks"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/storage"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestAuthMiddleware(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	settings, err := config.Setup(false)
	require.NoError(t, err)

	logger := zap.NewNop()
	store := mocks.NewMockStorager(mockCtrl)

	ctx := context.Background()
	userID := 1
	userToken := buildJWTString(t, settings, userID)

	type want struct {
		code int
	}

	tests := []struct {
		name           string
		withAuthHeader bool
		userToken      string
		mockStorage    func()
		want           want
	}{
		{
			name:           "success auth",
			withAuthHeader: true,
			userToken:      userToken,
			mockStorage: func() {
				store.EXPECT().GetUserByID(ctx, userID).Times(1).Return(models.User{}, nil)
			},
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name:           "without auth token",
			withAuthHeader: false,
			userToken:      "",
			mockStorage:    func() {},
			want: want{
				code: http.StatusUnauthorized,
			},
		},
		{
			name:           "failed auth token",
			withAuthHeader: true,
			userToken:      "token",
			mockStorage:    func() {},
			want: want{
				code: http.StatusUnauthorized,
			},
		},
		{
			name:           "when user not found",
			withAuthHeader: true,
			userToken:      userToken,
			mockStorage: func() {
				store.EXPECT().GetUserByID(ctx, userID).Times(1).Return(models.User{}, storage.ErrUserNotFound)
			},
			want: want{
				code: http.StatusUnauthorized,
			},
		},
		{
			name:           "failed get user from storage",
			withAuthHeader: true,
			userToken:      userToken,
			mockStorage: func() {
				store.EXPECT().GetUserByID(ctx, userID).Times(1).Return(models.User{}, errors.New("some error"))
			},
			want: want{
				code: http.StatusUnauthorized,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockStorage()

			someHandler := func(w http.ResponseWriter, r *http.Request) {}

			request := httptest.NewRequest(http.MethodGet, "/", http.NoBody)

			if test.withAuthHeader {
				request.Header.Add("X-Auth-Token", test.userToken)
			}

			w := httptest.NewRecorder()

			m := authMiddleware(settings, logger, store)(http.HandlerFunc(someHandler))
			m.ServeHTTP(w, request)

			res := w.Result()
			defer closeBody(t, res)

			assert.Equal(t, test.want.code, res.StatusCode)
		})
	}
}

func buildJWTString(t *testing.T, settings *config.Settings, userID int) string {
	t.Helper()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, models.Claims{
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(settings.SecretKey))
	require.NoError(t, err)

	return tokenString
}
