package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/MihailSergeenkov/GophKeeper/internal/models"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/config"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/storage"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserValidationFields = errors.New("some fields have not been validated")
	ErrUserLoginExist       = errors.New("user already exist")
	ErrUserLoginCreds       = errors.New("user has invalid login or password")
)

// RegisterUser зарегистрировать нового пользователя.
func (s *Services) RegisterUser(ctx context.Context, req models.RegisterUserRequest) error {
	if err := validateRegisterUserRequest(req); err != nil {
		return fmt.Errorf("failed to validate fields %w", err)
	}

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return fmt.Errorf("failed to hash passwords %w", err)
	}

	err = s.storage.AddUser(ctx, req.Login, hashedPassword)
	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) && pgxError.Code == pgerrcode.UniqueViolation {
			return ErrUserLoginExist
		}
		return fmt.Errorf("failed to add user %w", err)
	}

	return nil
}

// CreateUserToken получить ключ доступа для пользователя.
func (s *Services) CreateUserToken(
	ctx context.Context,
	req models.CreateUserTokenRequest,
) (models.CreateUserTokenResponse, error) {
	resp := models.CreateUserTokenResponse{}

	user, err := s.storage.GetUserByLogin(ctx, req.Login)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return resp, ErrUserLoginCreds
		}
		return resp, fmt.Errorf("failed to get user from DB %w", err)
	}

	if err := verifyPassword(user.Password, req.Password); err != nil {
		return resp, ErrUserLoginCreds
	}

	authToken, err := buildJWTString(s.settings, user.ID)
	if err != nil {
		return resp, fmt.Errorf("failed to build auth token: %w", err)
	}

	resp.AuthToken = authToken

	return resp, nil
}

func validateRegisterUserRequest(req models.RegisterUserRequest) error {
	if req.Login == "" {
		return ErrUserValidationFields
	}
	if req.Password == "" {
		return ErrUserValidationFields
	}

	return nil
}

func hashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("could not hash password %w", err)
	}

	return hashedPassword, nil
}

func verifyPassword(hashedPassword []byte, candidatePassword string) error {
	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(candidatePassword))
	if err != nil {
		return fmt.Errorf("failed to compare: %w", err)
	}

	return nil
}

func buildJWTString(settings *config.Settings, userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, models.Claims{
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(settings.SecretKey))
	if err != nil {
		return "", fmt.Errorf("failed to signed token: %w", err)
	}

	return tokenString, nil
}
