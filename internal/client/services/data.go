package services

import (
	"fmt"
	"net/http"

	"github.com/MihailSergeenkov/GophKeeper/internal/client/requests"
	"github.com/MihailSergeenkov/GophKeeper/internal/models"
)

// SyncData сервис синхронизации данных с сервера.
func (s *Services) SyncData() error {
	const path = "/user/data"
	userData := []models.UserData{}

	resp, err := s.httpRequests.Get(
		s.cfg.GetServerAPI()+path,
		requests.WithHeader(ContentTypeHeader, JSONContentType),
		requests.WithHeader(AuthHeader, s.cfg.GetToken()),
		requests.WithResult(&userData),
	)
	if err != nil {
		return failedRequest(err)
	}

	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusNoContent {
		return failedResponseStatus(resp.Status())
	}

	if err := s.cfg.UpdateData(userData); err != nil {
		return fmt.Errorf("failed to update data: %w", err)
	}
	return nil
}

// GetData сервис получения данных из кеша.
func (s *Services) GetData() []models.UserData {
	userData := []models.UserData{}

	for _, v := range s.cfg.GetData() {
		userData = append(userData, v)
	}

	return userData
}
