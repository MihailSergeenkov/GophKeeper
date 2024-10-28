package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/MihailSergeenkov/GophKeeper/internal/models"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/services"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// AddPassword обработчик для добавления данных пароля пользователя.
func (h *Handlers) AddPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.AddPasswordRequest

		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			h.logger.Error(readReqErrStr, zap.Error(err))
			return
		}

		id, err := h.services.AddPassword(r.Context(), req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.logger.Error("failed to add password", zap.Error(err))
			return
		}

		w.Header().Set(ContentTypeHeader, JSONContentType)
		w.WriteHeader(http.StatusCreated)

		enc := json.NewEncoder(w)
		if err := enc.Encode(models.AddResponse{ID: id}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.logger.Error(encRespErrStr, zap.Error(err))
			return
		}
	}
}

// GetPassword обработчик для получения данных конкретного пароля пользователя.
func (h *Handlers) GetPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "passwordID")
		passwordID, err := strconv.Atoi(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			h.logger.Error("failed password ID param", zap.Error(err))
			return
		}

		password, err := h.services.GetPassword(r.Context(), passwordID)
		if err != nil {
			if errors.Is(err, services.ErrNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			h.logger.Error("failed to get password", zap.Error(err))
			return
		}

		w.Header().Set(ContentTypeHeader, JSONContentType)
		w.WriteHeader(http.StatusOK)

		enc := json.NewEncoder(w)
		if err := enc.Encode(password); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.logger.Error(encRespErrStr, zap.Error(err))
			return
		}
	}
}
