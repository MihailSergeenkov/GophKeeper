package services

import (
	"fmt"
	"net/http"

	"github.com/MihailSergeenkov/GophKeeper/internal/models"
)

// SyncData сервис синхронизации данных с сервера.
func SyncData(cfg Configurer) error {
	const path = "/user/data"
	client := getClient(cfg)

	userData := []models.UserData{}

	resp, err := client.R().
		SetHeader(ContentTypeHeader, JSONContentType).
		SetHeader(AuthHeader, cfg.GetToken()).
		SetResult(&userData).
		Get(cfg.GetServerAPI() + path)

	if err != nil {
		return fmt.Errorf("failed request: %w", err)
	}

	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusNoContent {
		return fmt.Errorf("response status: %s", resp.Status())
	}

	cfg.UpdateData(userData)
	return nil
}

// GetData сервис получения данных из кеша.
func GetData(cfg Configurer) []models.UserData {
	userData := []models.UserData{}

	for _, v := range cfg.GetData() {
		userData = append(userData, v)
	}

	return userData
}
