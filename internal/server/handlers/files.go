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

// AddFile обработчик для добавления файла пользователя.
func (h *Handlers) AddFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.logger.Error("failed to add file", zap.Error(err))
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			h.logger.Error("failed to retrieving the file", zap.Error(err))
			return
		}
		defer file.Close()

		req := models.AddFileRequest{
			Mark:        r.PostFormValue("mark"),
			Description: r.PostFormValue("description"),
			File:        file,
			FileName:    header.Filename,
			FileSize:    header.Size,
		}

		id, err := h.services.AddFile(r.Context(), req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.logger.Error("failed to add file", zap.Error(err))
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

// GetFile обработчик для получения конкретного файла пользователя.
func (h *Handlers) GetFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "fileID")
		fileID, err := strconv.Atoi(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			h.logger.Error("failed ID param", zap.Error(err))
			return
		}

		file, err := h.services.GetFile(r.Context(), fileID)
		if err != nil {
			if errors.Is(err, services.ErrNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			h.logger.Error("failed to get file", zap.Error(err))
			return
		}
		// defer file.File.Close()

		// fileBytes, err := io.ReadAll(file.File)
		// if err != nil {
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	h.logger.Error("failed to prepare file", zap.Error(err))
		// 	return
		// }

		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)

		// w.Write(fileBytes)
		w.Write(file.FileBytes)
	}
}