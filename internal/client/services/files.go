package services

import (
	"errors"
	"net/http"

	"github.com/MihailSergeenkov/GophKeeper/internal/client/requests"
	"github.com/MihailSergeenkov/GophKeeper/internal/models"
)

// AddFile сервис добавления файла.
func (s *Services) AddFile(filePath, mark, description string) error {
	const path = "/user/files"

	addResp := models.AddResponse{}

	resp, err := s.httpRequests.Post(
		s.cfg.GetServerAPI()+path,
		requests.WithHeader(ContentTypeHeader, FormDataContentType),
		requests.WithHeader(AuthHeader, s.cfg.GetToken()),
		requests.WithFile(filePath),
		requests.WithFormData(map[string]string{
			"mark":        mark,
			"description": description,
		}),
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
		Mark:        mark,
		Description: description,
		Type:        "file",
	}

	if err := s.cfg.AddData(d); err != nil {
		return failedDumpData(err)
	}

	return nil
}

// GetFile сервис получения файла.
func (s *Services) GetFile(id, dir string) error {
	const path = "/user/files/{id}"

	userData, ok := s.cfg.GetData()[id]
	if !ok {
		return errors.New("file id not found")
	}

	resp, err := s.httpRequests.Get(
		s.cfg.GetServerAPI()+path,
		requests.WithHeader(AuthHeader, s.cfg.GetToken()),
		requests.WithPathParams(map[string]string{"id": id}),
		requests.WithOutput(dir+"/"+userData.Mark),
	)

	if err != nil {
		return failedRequest(err)
	}
	if resp.StatusCode() != http.StatusOK {
		return failedResponseStatus(resp.Status())
	}

	return nil
}
