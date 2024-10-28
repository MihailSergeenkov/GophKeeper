package services

import (
	"errors"
	"net/http"

	"github.com/MihailSergeenkov/GophKeeper/internal/models"
)

// AddFile сервис добавления файла.
func AddFile(cfg Configurer, filePath, mark, description string) error {
	const path = "/user/files"
	client := getClient(cfg)

	addResp := models.AddResponse{}

	resp, err := client.R().
		SetHeader(ContentTypeHeader, FormDataContentType).
		SetHeader(AuthHeader, cfg.GetToken()).
		SetFile("file", filePath).
		SetFormData(map[string]string{
			"mark":        mark,
			"description": description,
		}).
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
		Mark:        mark,
		Description: description,
		Type:        "file",
	}

	if err := cfg.AddData(d); err != nil {
		return failedDumpData(err)
	}

	return nil
}

// GetFile сервис получения файла.
func GetFile(cfg Configurer, id, dir string) error {
	const path = "/user/files/{id}"

	userData, ok := cfg.GetData()[id]
	if !ok {
		return errors.New("file id not found")
	}

	client := getClient(cfg)
	client.SetOutputDirectory(dir)
	resp, err := client.R().
		SetHeader(AuthHeader, cfg.GetToken()).
		SetPathParams(map[string]string{"id": id}).
		SetOutput(userData.Mark).
		Get(cfg.GetServerAPI() + path)

	if err != nil {
		return failedRequest(err)
	}
	if resp.StatusCode() != http.StatusOK {
		return failedResponseStatus(resp.Status())
	}

	return nil
}
