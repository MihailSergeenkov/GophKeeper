package services

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/MihailSergeenkov/GophKeeper/internal/client/requests"
	"github.com/MihailSergeenkov/GophKeeper/internal/models"
)

// AddPassword сервис добавления данных логин-пароль.
func (s *Services) AddPassword(req models.AddPasswordRequest) error {
	const path = "/user/passwords"

	body, err := json.Marshal(req)
	if err != nil {
		return failedCreateBody(err)
	}

	addResp := models.AddResponse{}

	resp, err := s.httpRequests.Post(
		s.cfg.GetServerAPI()+path,
		requests.WithHeader(ContentTypeHeader, JSONContentType),
		requests.WithHeader(AuthHeader, s.cfg.GetToken()),
		requests.WithBody(body),
		requests.WithResult(&addResp),
	)
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

	if err := s.cfg.AddData(d); err != nil {
		return failedDumpData(err)
	}

	return nil
}

// GetPassword сервис получения данных логин-пароль.
func (s *Services) GetPassword(id string) (models.Password, error) {
	const path = "/user/passwords/{id}"

	password := models.Password{}

	if _, ok := s.cfg.GetData()[id]; !ok {
		return password, errors.New("password id not found")
	}

	resp, err := s.httpRequests.Get(
		s.cfg.GetServerAPI()+path,
		requests.WithHeader(ContentTypeHeader, JSONContentType),
		requests.WithHeader(AuthHeader, s.cfg.GetToken()),
		requests.WithPathParams(map[string]string{"id": id}),
		requests.WithResult(&password),
	)
	if err != nil {
		return password, failedRequest(err)
	}
	if resp.StatusCode() != http.StatusOK {
		return password, failedResponseStatus(resp.Status())
	}

	return password, nil
}
