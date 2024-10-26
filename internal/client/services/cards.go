package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MihailSergeenkov/GophKeeper/internal/models"
)

// AddCard сервис добавления данных банковской карты.
func AddCard(cfg Configurer, req models.AddCardRequest) error {
	const path = "/user/cards"
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
		Type:        "card",
	}

	if err := cfg.AddData(d); err != nil {
		return fmt.Errorf("failed to dump data: %w", err)
	}

	return nil
}

// GetCard сервис получения данных банковской карты.
func GetCard(cfg Configurer, id string) (models.Card, error) {
	const path = "/user/cards/{id}"

	card := models.Card{}

	if _, ok := cfg.GetData()[id]; !ok {
		return card, fmt.Errorf("card id not found")
	}

	client := getClient(cfg)
	resp, err := client.R().
		SetHeader(ContentTypeHeader, JSONContentType).
		SetHeader(AuthHeader, cfg.GetToken()).
		SetPathParams(map[string]string{"id": id}).
		SetResult(&card).
		Get(cfg.GetServerAPI() + path)

	if err != nil {
		return card, fmt.Errorf("failed request: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return card, fmt.Errorf("response status: %s", resp.Status())
	}

	return card, nil
}
