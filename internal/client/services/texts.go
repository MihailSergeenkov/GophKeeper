package services

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/MihailSergeenkov/GophKeeper/internal/client/requests"
	"github.com/MihailSergeenkov/GophKeeper/internal/models"
)

// AddText сервис добавления текста.
func (s *Services) AddText(req models.AddTextRequest) error {
	const path = "/user/texts"

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
		Type:        "text",
	}

	if err := s.cfg.AddData(d); err != nil {
		return failedDumpData(err)
	}

	return nil
}

// GetText сервис получения текста.
func (s *Services) GetText(id string) (models.Text, error) {
	const path = "/user/texts/{id}"

	text := models.Text{}

	if _, ok := s.cfg.GetData()[id]; !ok {
		return text, errors.New("text id not found")
	}

	resp, err := s.httpRequests.Get(
		s.cfg.GetServerAPI()+path,
		requests.WithHeader(ContentTypeHeader, JSONContentType),
		requests.WithHeader(AuthHeader, s.cfg.GetToken()),
		requests.WithPathParams(map[string]string{"id": id}),
		requests.WithResult(&text),
	)
	if err != nil {
		return text, failedRequest(err)
	}
	if resp.StatusCode() != http.StatusOK {
		return text, failedResponseStatus(resp.Status())
	}

	return text, nil
}
