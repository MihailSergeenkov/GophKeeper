package routes

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/MihailSergeenkov/GophKeeper/internal/server/config"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/constants"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/models"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/storage"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

func authMiddleware(settings *config.Settings, l *zap.Logger, s Storager) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authToken := r.Header.Get("X-Auth-Token")
			if authToken == "" {
				w.WriteHeader(http.StatusUnauthorized)
				l.Error("failed to get auth token")
				return
			}

			userID := getUserID(settings, authToken)
			if userID == -1 {
				w.WriteHeader(http.StatusUnauthorized)
				l.Error("failed to parse auth token")
				return
			}

			_, err := s.GetUserByID(r.Context(), userID)
			if err != nil {
				if errors.Is(err, storage.ErrUserNotFound) {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				w.WriteHeader(http.StatusUnauthorized)
				l.Error("failed to get user from DB", zap.Error(err))
				return
			}

			newContext := context.WithValue(r.Context(), constants.KeyUserID, userID)
			newRequest := r.WithContext(newContext)
			next.ServeHTTP(w, newRequest)
		})
	}
}

func getUserID(settings *config.Settings, tokenString string) int {
	claims := &models.Claims{}

	_, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(settings.SecretKey), nil
	})

	if err != nil {
		return -1
	}

	return claims.UserID
}
