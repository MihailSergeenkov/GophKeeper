package services

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/MihailSergeenkov/GophKeeper/internal/client/requests"
	"github.com/MihailSergeenkov/GophKeeper/internal/models"
)

// AddCard сервис добавления данных банковской карты.
func (s *Services) AddCard(req *models.AddCardRequest) error {
	const path = "/user/cards"

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
		Type:        "card",
	}

	if err := s.cfg.AddData(d); err != nil {
		return failedDumpData(err)
	}

	return nil
}

// GetCard сервис получения данных банковской карты.
func (s *Services) GetCard(id string) (models.Card, error) {
	const path = "/user/cards/{id}"

	card := models.Card{}

	if _, ok := s.cfg.GetData()[id]; !ok {
		return card, errors.New("card id not found")
	}

	resp, err := s.httpRequests.Get(
		s.cfg.GetServerAPI()+path,
		requests.WithHeader(ContentTypeHeader, JSONContentType),
		requests.WithHeader(AuthHeader, s.cfg.GetToken()),
		requests.WithPathParams(map[string]string{"id": id}),
		requests.WithResult(&card),
	)
	if err != nil {
		return card, failedRequest(err)
	}
	if resp.StatusCode() != http.StatusOK {
		return card, failedResponseStatus(resp.Status())
	}

	return card, nil
}
