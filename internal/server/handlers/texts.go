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

// AddText обработчик для добавления текстовых данных пользователя.
func (h *Handlers) AddText() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.AddTextRequest

		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			h.logger.Error(readReqErrStr, zap.Error(err))
			return
		}

		id, err := h.services.AddText(r.Context(), req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.logger.Error("failed to add text", zap.Error(err))
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

// GetText обработчик для получения конкретных текстовых данных пользователя.
func (h *Handlers) GetText() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "textID")
		textID, err := strconv.Atoi(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			h.logger.Error("failed ID param", zap.Error(err))
			return
		}

		text, err := h.services.GetText(r.Context(), textID)
		if err != nil {
			if errors.Is(err, services.ErrNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			h.logger.Error("failed to get text", zap.Error(err))
			return
		}

		w.Header().Set(ContentTypeHeader, JSONContentType)
		w.WriteHeader(http.StatusOK)

		enc := json.NewEncoder(w)
		if err := enc.Encode(text); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.logger.Error(encRespErrStr, zap.Error(err))
			return
		}
	}
}
