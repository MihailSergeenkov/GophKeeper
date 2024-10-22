package models

import (
	"io"

	"github.com/golang-jwt/jwt/v5"
)

// RegisterUserRequest тип для регистрации пользователя.
type RegisterUserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// CreateUserTokenRequest тип для получения токена доступа пользователя.
type CreateUserTokenRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// CreateUserTokenResponse тип ответа с токеном доступа пользователя.
type CreateUserTokenResponse struct {
	AuthToken string `json:"auth_token"`
}

// User тип пользователя.
type User struct {
	Login    string
	Password []byte
	ID       int
}

// Claims тип для данных токена доступа.
type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

// UserData тип для данных пользователя.
type UserData struct {
	ID          int    `json:"id"`
	Mark        string `json:"mark,omitempty"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type,omitempty"`
}

// AddPasswordRequest тип для добавления пароля пользователя.
type AddPasswordRequest struct {
	Login       string `json:"login"`
	Password    string `json:"password"`
	Mark        string `json:"mark"`
	Description string `json:"description"`
}

// AddCardRequest тип для добавления карты пользователя.
type AddCardRequest struct {
	Number      string `json:"number"`
	Owner       string `json:"owner"`
	ExpiryDate  string `json:"expiry_date"`
	CVV2        string `json:"cvv2"`
	Mark        string `json:"mark"`
	Description string `json:"description"`
}

// AddTextRequest тип для добавления текста пользователя.
type AddTextRequest struct {
	Data        string `json:"data"`
	Mark        string `json:"mark"`
	Description string `json:"description"`
}

// AddFileRequest тип для добавления файла пользователя.
type AddFileRequest struct {
	File        io.Reader
	FileName    string
	FileSize    int64
	Mark        string
	Description string
}

// Password тип для пароля пользователя.
type Password struct {
	ID          int    `json:"id"`
	Login       string `json:"login"`
	Password    string `json:"password"`
	Mark        string `json:"mark"`
	Description string `json:"description"`
}

// Card тип для карты пользователя.
type Card struct {
	ID          int    `json:"id"`
	Number      string `json:"number"`
	Owner       string `json:"owner"`
	ExpiryDate  string `json:"expiry_date"`
	CVV2        string `json:"cvv2"`
	Mark        string `json:"mark"`
	Description string `json:"description"`
}

// Text тип для текста пользователя.
type Text struct {
	ID          int    `json:"id"`
	Data        string `json:"data"`
	Mark        string `json:"mark"`
	Description string `json:"description"`
}

// File тип для файла пользователя.
type File struct {
	File      io.ReadCloser
	FileBytes []byte // специально для теста
}

// EncryptPasswordData тип для шифрованных данных пароля пользователя.
type EncryptPasswordData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// EncryptCardData тип для шифрованных данных карты пользователя.
type EncryptCardData struct {
	Number     string `json:"number"`
	Owner      string `json:"owner"`
	ExpiryDate string `json:"expiry_date"`
	CVV2       string `json:"cvv2"`
}

// EncryptTextData тип для шифрованных данных текста пользователя.
type EncryptTextData struct {
	Data string `json:"data"`
}

// EncryptFileData тип для шифрованных данных файла пользователя.
type EncryptFileData struct {
	FileName string `json:"file_name"`
}
