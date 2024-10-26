package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MihailSergeenkov/GophKeeper/internal/models"
)

// AddPassword сервис добавления данных логин-пароль.
func AddPassword(cfg Configurer, req models.AddPasswordRequest) error {
	const path = "/user/passwords"
	client := getClient(cfg)

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to create body: %w", err)
	}

	addResp := models.AddResponse{}

	resp, err := client.R().
		SetHeader(ContentTypeHeader, JSONContentType).
		SetHeader(AuthHeader, cfg.GetToken()).
		SetBody(body).
		SetResult(&addResp).
		Post(cfg.GetServerAPI() + path)

	if err != nil {
		return fmt.Errorf("failed request: %w", err)
	}
	if resp.StatusCode() != http.StatusCreated {
		return fmt.Errorf("response status: %s", resp.Status())
	}

	d := models.UserData{
		ID:          addResp.ID,
		Mark:        req.Mark,
		Description: req.Description,
		Type:        "password",
	}

	if err := cfg.AddData(d); err != nil {
		return fmt.Errorf("failed to dump data: %w", err)
	}

	return nil
}

// GetPassword сервис получения данных логин-пароль.
func GetPassword(cfg Configurer, id string) (models.Password, error) {
	const path = "/user/passwords/{id}"

	password := models.Password{}

	if _, ok := cfg.GetData()[id]; !ok {
		return password, fmt.Errorf("password id not found")
	}

	client := getClient(cfg)
	resp, err := client.R().
		SetHeader(ContentTypeHeader, JSONContentType).
		SetHeader(AuthHeader, cfg.GetToken()).
		SetPathParams(map[string]string{"id": id}).
		SetResult(&password).
		Get(cfg.GetServerAPI() + path)

	if err != nil {
		return password, fmt.Errorf("failed request: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return password, fmt.Errorf("response status: %s", resp.Status())
	}

	return password, nil
}
