package services

import (
	"fmt"
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
		return fmt.Errorf("failed request: %w", err)
	}
	if resp.StatusCode() != http.StatusCreated {
		return fmt.Errorf("response status: %s", resp.Status())
	}

	d := models.UserData{
		ID:          addResp.ID,
		Mark:        mark,
		Description: description,
		Type:        "file",
	}

	if err := cfg.AddData(d); err != nil {
		return fmt.Errorf("failed to dump data: %w", err)
	}

	return nil
}

// GetFile сервис получения файла.
func GetFile(cfg Configurer, id, dir string) error {
	const path = "/user/files/{id}"

	userData, ok := cfg.GetData()[id]
	if !ok {
		return fmt.Errorf("file id not found")
	}

	client := getClient(cfg)
	client.SetOutputDirectory(dir)
	resp, err := client.R().
		SetHeader(AuthHeader, cfg.GetToken()).
		SetPathParams(map[string]string{"id": id}).
		SetOutput(userData.Mark).
		Get(cfg.GetServerAPI() + path)

	if err != nil {
		return fmt.Errorf("failed request: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("response status: %s", resp.Status())
	}

	return nil
}
