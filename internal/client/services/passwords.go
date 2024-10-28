package services

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/MihailSergeenkov/GophKeeper/internal/models"
)

// AddPassword сервис добавления данных логин-пароль.
func AddPassword(cfg Configurer, req models.AddPasswordRequest) error {
	const path = "/user/passwords"
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
		Type:        "password",
	}

	if err := cfg.AddData(d); err != nil {
		return failedDumpData(err)
	}

	return nil
}

// GetPassword сервис получения данных логин-пароль.
func GetPassword(cfg Configurer, id string) (models.Password, error) {
	const path = "/user/passwords/{id}"

	password := models.Password{}

	if _, ok := cfg.GetData()[id]; !ok {
		return password, errors.New("password id not found")
	}

	client := getClient(cfg)
	resp, err := client.R().
		SetHeader(ContentTypeHeader, JSONContentType).
		SetHeader(AuthHeader, cfg.GetToken()).
		SetPathParams(map[string]string{"id": id}).
		SetResult(&password).
		Get(cfg.GetServerAPI() + path)

	if err != nil {
		return password, failedRequest(err)
	}
	if resp.StatusCode() != http.StatusOK {
		return password, failedResponseStatus(resp.Status())
	}

	return password, nil
}
