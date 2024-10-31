package services

import (
	"errors"
	"net/http"
	"strings"

	"github.com/MihailSergeenkov/GophKeeper/internal/client/requests"
	"github.com/MihailSergeenkov/GophKeeper/internal/models"
)

// AddFile сервис добавления файла.
func (s *Services) AddFile(filePath, mark, description string) error {
	const path = "/user/files"

	addResp := models.AddResponse{}

	preparedMark := strings.ReplaceAll(strings.ToLower(mark), " ", "_")

	resp, err := s.httpRequests.Post(
		s.cfg.GetServerAPI()+path,
		requests.WithHeader(ContentTypeHeader, FormDataContentType),
		requests.WithHeader(AuthHeader, s.cfg.GetToken()),
		requests.WithFile(filePath),
		requests.WithFormData(map[string]string{
			"mark":        preparedMark,
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
		Mark:        preparedMark,
		Description: description,
		Type:        "file",
	}

	if err := s.cfg.AddData(d); err != nil {
		return failedDumpData(err)
	}

	return nil
}

// GetFile сервис получения файла.
func (s *Services) GetFile(fileMark, dir string) error {
	const path = "/user/files/{fileMark}"

	if _, ok := s.cfg.GetData()[fileMark]; !ok {
		return errors.New("file mark not found")
	}

	resp, err := s.httpRequests.Get(
		s.cfg.GetServerAPI()+path,
		requests.WithHeader(AuthHeader, s.cfg.GetToken()),
		requests.WithPathParams(map[string]string{"fileMark": fileMark}),
		requests.WithOutput(dir+"/"+fileMark),
	)

	if err != nil {
		return failedRequest(err)
	}
	if resp.StatusCode() != http.StatusOK {
		return failedResponseStatus(resp.Status())
	}

	return nil
}
