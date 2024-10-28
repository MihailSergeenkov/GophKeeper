package services

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/MihailSergeenkov/GophKeeper/internal/models"
)

// AddText сервис добавления текста.
func AddText(cfg Configurer, req models.AddTextRequest) error {
	const path = "/user/texts"
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
		Type:        "text",
	}

	if err := cfg.AddData(d); err != nil {
		return failedDumpData(err)
	}

	return nil
}

// GetText сервис получения текста.
func GetText(cfg Configurer, id string) (models.Text, error) {
	const path = "/user/texts/{id}"

	text := models.Text{}

	if _, ok := cfg.GetData()[id]; !ok {
		return text, errors.New("text id not found")
	}

	client := getClient(cfg)
	resp, err := client.R().
		SetHeader(ContentTypeHeader, JSONContentType).
		SetHeader(AuthHeader, cfg.GetToken()).
		SetPathParams(map[string]string{"id": id}).
		SetResult(&text).
		Get(cfg.GetServerAPI() + path)

	if err != nil {
		return text, failedRequest(err)
	}
	if resp.StatusCode() != http.StatusOK {
		return text, failedResponseStatus(resp.Status())
	}

	return text, nil
}
