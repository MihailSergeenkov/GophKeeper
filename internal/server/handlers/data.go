package handlers

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

// FetchUserData обработчик для получения базовой информации о данных пользователя.
func (h *Handlers) FetchUserData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := h.services.FetchUserData(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.logger.Error("failed to fetch user data from DB", zap.Error(err))
			return
		}

		if len(data) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		w.Header().Set(ContentTypeHeader, JSONContentType)
		w.WriteHeader(http.StatusOK)

		enc := json.NewEncoder(w)
		if err := enc.Encode(data); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.logger.Error(encRespErrStr, zap.Error(err))
			return
		}
	}
}
