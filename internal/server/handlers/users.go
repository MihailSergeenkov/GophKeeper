package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/MihailSergeenkov/GophKeeper/internal/models"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/services"
	"go.uber.org/zap"
)

// RegisterUser обработчик для регистрации пользователя.
func (h *Handlers) RegisterUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.RegisterUserRequest

		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			h.logger.Error(readReqErrStr, zap.Error(err))
			return
		}

		err := h.services.RegisterUser(r.Context(), req)

		if err != nil {
			if errors.Is(err, services.ErrUserValidationFields) {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if errors.Is(err, services.ErrUserLoginExist) {
				w.WriteHeader(http.StatusConflict)
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			h.logger.Error("failed to register user", zap.Error(err))
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// CreateUserToken обработчик для запроса ключа доступа пользователя.
func (h *Handlers) CreateUserToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.CreateUserTokenRequest

		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			h.logger.Error(readReqErrStr, zap.Error(err))
			return
		}

		resp, err := h.services.CreateUserToken(r.Context(), req)

		if err != nil {
			if errors.Is(err, services.ErrUserLoginCreds) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			h.logger.Error("failed to create user token", zap.Error(err))
			return
		}

		w.Header().Set(ContentTypeHeader, JSONContentType)
		w.WriteHeader(http.StatusOK)

		enc := json.NewEncoder(w)
		if err := enc.Encode(resp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.logger.Error(encRespErrStr, zap.Error(err))
			return
		}
	}
}
