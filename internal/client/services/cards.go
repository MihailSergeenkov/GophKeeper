package services

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/MihailSergeenkov/GophKeeper/internal/models"
)

// AddCard сервис добавления данных банковской карты.
func AddCard(cfg Configurer, req *models.AddCardRequest) error {
	const path = "/user/cards"
	client := getClient(cfg)

	body, err := json.Marshal(req)
	if err != nil {
		return failedCreateBody(err)
	}

	addResp := models.AddResponse{}

	resp, err := client.R().
		SetHeader(ContentTypeHeader, JSONContentType).
		SetHeader(AuthHeader, cfg.GetToken()).
		SetBody(body).
		SetResult(&addResp).
		Post(cfg.GetServerAPI() + path)

	if err != nil {
		return failedRequest(err)
	}
	if resp.StatusCode() != http.StatusCreated {
		return failedResponseStatus(resp.Status())
	}

	d := models.UserData{
		ID:          addResp.ID,
		Mark:        req.Mark,
		Description: req.Description,
		Type:        "card",
	}

	if err := cfg.AddData(d); err != nil {
		return failedDumpData(err)
	}

	return nil
}

// GetCard сервис получения данных банковской карты.
func GetCard(cfg Configurer, id string) (models.Card, error) {
	const path = "/user/cards/{id}"

	card := models.Card{}

	if _, ok := cfg.GetData()[id]; !ok {
		return card, errors.New("card id not found")
	}

	client := getClient(cfg)
	resp, err := client.R().
		SetHeader(ContentTypeHeader, JSONContentType).
		SetHeader(AuthHeader, cfg.GetToken()).
		SetPathParams(map[string]string{"id": id}).
		SetResult(&card).
		Get(cfg.GetServerAPI() + path)

	if err != nil {
		return card, failedRequest(err)
	}
	if resp.StatusCode() != http.StatusOK {
		return card, failedResponseStatus(resp.Status())
	}

	return card, nil
}
