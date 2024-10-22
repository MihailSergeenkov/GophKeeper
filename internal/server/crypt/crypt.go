package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"fmt"

	"github.com/MihailSergeenkov/GophKeeper/internal/server/config"
)

// Crypt структура для работы с функциями криптографии приложения.
type Crypt struct {
	settings *config.Settings
	aesgcm   cipher.AEAD
	nonce    []byte
}

// NewCrypt функция инициализации криптографии приложения.
func NewCrypt(settings *config.Settings) (*Crypt, error) {
	key := sha256.Sum256([]byte(settings.SecretKey))

	aesblock, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher block %w", err)
	}
	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher aead %w", err)
	}

	nonce := key[len(key)-aesgcm.NonceSize():]

	return &Crypt{
		settings: settings,
		aesgcm:   aesgcm,
		nonce:    nonce,
	}, nil
}

// EncryptData функция зашифровки данных.
func (c Crypt) EncryptData(data []byte) []byte {
	encrypted := c.aesgcm.Seal(nil, c.nonce, data, nil)
	return encrypted
}

// DecryptData функция расшифровки данных.
func (c Crypt) DecryptData(encrypted []byte) ([]byte, error) {
	decrypted, err := c.aesgcm.Open(nil, c.nonce, encrypted, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data %w", err)
	}

	return decrypted, nil
}
