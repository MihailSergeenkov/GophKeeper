package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/MihailSergeenkov/GophKeeper/internal/server/models"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/services"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// AddCard обработчик для добавления данных карты пользователя.
func (h *Handlers) AddCard() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.AddCardRequest

		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			h.logger.Error(readReqErrStr, zap.Error(err))
			return
		}

		id, err := h.services.AddCard(r.Context(), req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.logger.Error("failed to add card", zap.Error(err))
			return
		}

		w.Header().Set(ContentTypeHeader, JSONContentType)
		w.WriteHeader(http.StatusCreated)

		enc := json.NewEncoder(w)
		if err := enc.Encode(models.UserData{ID: id}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.logger.Error(encRespErrStr, zap.Error(err))
			return
		}
	}
}

// GetCard обработчик для получения данных конкретной карты пользователя.
func (h *Handlers) GetCard() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "cardID")
		cardID, err := strconv.Atoi(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			h.logger.Error("failed ID param", zap.Error(err))
			return
		}

		card, err := h.services.GetCard(r.Context(), cardID)
		if err != nil {
			if errors.Is(err, services.ErrNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			h.logger.Error("failed to get card", zap.Error(err))
			return
		}

		w.Header().Set(ContentTypeHeader, JSONContentType)
		w.WriteHeader(http.StatusOK)

		enc := json.NewEncoder(w)
		if err := enc.Encode(card); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.logger.Error(encRespErrStr, zap.Error(err))
			return
		}
	}
}
