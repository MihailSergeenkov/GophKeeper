package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/MihailSergeenkov/GophKeeper/internal/server/models"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/storage"
)

const fileDataType = "file"

// AddFile функция для добавления файла пользователя.
func (s *Services) AddFile(ctx context.Context, req models.AddFileRequest) (int, error) {
	if err := validateAddFileRequest(req); err != nil {
		return 0, fmt.Errorf("failed to validate fields %w", err)
	}

	// не работает
	f, _ := io.ReadAll(req.File)
	encFile := s.crypter.EncryptData(f)
	log.Print(encFile)
	//

	if err := s.fileStorage.AddFile(ctx, bytes.NewReader(encFile), req.FileName, req.FileSize); err != nil {
		return 0, fmt.Errorf("failed to add file to filestorage %w", err)
	}

	data := models.EncryptFileData{
		FileName: req.FileName,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return 0, fmt.Errorf("failed to generate json data %w", err)
	}

	encData := s.crypter.EncryptData(jsonData)

	id, err := s.storage.AddUserData(ctx, encData, req.Mark, req.Description, fileDataType)
	if err != nil {
		return 0, fmt.Errorf("failed to add user data %w", err)
	}

	return id, nil
}

// GetFile функция для получения файла пользователя в виде массива байт.
func (s *Services) GetFile(ctx context.Context, id int) (models.File, error) {
	var resp models.File
	decData, _, _, err := s.storage.GetUserData(ctx, id, fileDataType)
	if err != nil {
		if errors.Is(err, storage.ErrUserDataNotFound) {
			return resp, ErrNotFound
		}

		return resp, fmt.Errorf("failed to get user data %w", err)
	}

	jsonData, err := s.crypter.DecryptData(decData)
	if err != nil {
		return resp, fmt.Errorf("failed to decrypt data %w", err)
	}

	var encData models.EncryptFileData

	if err = json.Unmarshal(jsonData, &encData); err != nil {
		return resp, fmt.Errorf("failed to generate data %w", err)
	}

	log.Print(encData.FileName)
	file, err := s.fileStorage.GetFile(ctx, encData.FileName)
	if err != nil {
		return resp, fmt.Errorf("failed to get file from filestorage %w", err)
	}

	// не работает
	f, _ := io.ReadAll(file)
	log.Print(f) // видно что возвращается из стора усеченное кол-во байт
	defer file.Close()
	decFile, err := s.crypter.DecryptData(f)
	if err != nil {
		return resp, fmt.Errorf("failed to decrypt file %w", err)
	}
	//

	// resp.File = file
	resp.FileBytes = decFile

	return resp, nil
}

func validateAddFileRequest(req models.AddFileRequest) error {
	if len([]rune(req.Mark)) > 50 {
		return ErrUserMarkIsTooBig
	}
	if len([]rune(req.Description)) > 50 {
		return ErrUserDescriptionIsTooBig
	}

	return nil
}
