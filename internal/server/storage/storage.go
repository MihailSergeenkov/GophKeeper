package storage

import (
	"context"
	"embed"
	"errors"
	"fmt"

	"github.com/MihailSergeenkov/GophKeeper/internal/server/constants"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/models"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrUserDataNotFound = errors.New("user data not found")
)

const failedScanStr = "failed to scan a response row: %w"

// DBPooler интерфейс к пулу БД.
type DBPooler interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Ping(ctx context.Context) error
	Close()
}

// Storage структура postgresql БД.
type Storage struct {
	pool   DBPooler
	logger *zap.Logger
}

// NewStorage инициализирует postgresql БД.
func NewStorage(ctx context.Context, logger *zap.Logger, dbURI string) (*Storage, error) {
	if err := runMigrations(dbURI); err != nil {
		return nil, fmt.Errorf("failed to run DB migrations: %w", err)
	}

	pool, err := initPool(ctx, logger, dbURI)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize a connection pool: %w", err)
	}

	s := &Storage{
		logger: logger,
		pool:   pool,
	}

	return s, nil
}

//go:embed migrations/*.sql
var migrationsDir embed.FS

func runMigrations(dbURI string) error {
	d, err := iofs.New(migrationsDir, "migrations")
	if err != nil {
		return fmt.Errorf("failed to return an iofs driver: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, dbURI)
	if err != nil {
		return fmt.Errorf("failed to get a new migrate instance: %w", err)
	}
	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("failed to apply migrations to the DB: %w", err)
		}
	}
	return nil
}

func initPool(ctx context.Context, logger *zap.Logger, dbURI string) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(dbURI)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the DB URI: %w", err)
	}

	poolCfg.ConnConfig.Tracer = &queryTracer{logger: logger}
	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize a connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping the DB: %w", err)
	}

	return pool, nil
}

// Ping проверяет работоспособность БД.
func (s *Storage) Ping(ctx context.Context) error {
	if err := s.pool.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping DB: %w", err)
	}

	return nil
}

// Close закрывает соединение с БД.
func (s *Storage) Close() error {
	s.pool.Close()
	return nil
}

// AddUser добавляет нового пользователя.
func (s *Storage) AddUser(ctx context.Context, userLogin string, userPassword []byte) error {
	const addUserQuery = `INSERT INTO users (login, password) VALUES ($1, $2)`

	_, err := s.pool.Exec(ctx, addUserQuery, userLogin, userPassword)
	if err != nil {
		return fmt.Errorf("failed to execute add user query: %w", err)
	}

	return nil
}

// GetUserByLogin получить пользователя по логину.
func (s *Storage) GetUserByLogin(ctx context.Context, userLogin string) (models.User, error) {
	const query = `SELECT id, login, password FROM users WHERE login = $1 LIMIT 1`

	row := s.pool.QueryRow(ctx, query, userLogin)

	var u models.User
	err := row.Scan(&u.ID, &u.Login, &u.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, fmt.Errorf("%w with ID: %s", ErrUserNotFound, userLogin)
		}

		return models.User{}, fmt.Errorf(failedScanStr, err)
	}

	return u, nil
}

// GetUserByID получить пользователя по его ID.
func (s *Storage) GetUserByID(ctx context.Context, userID int) (models.User, error) {
	const query = `SELECT id, login, password FROM users WHERE id = $1 LIMIT 1`

	row := s.pool.QueryRow(ctx, query, userID)

	var u models.User
	err := row.Scan(&u.ID, &u.Login, &u.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, fmt.Errorf("%w with ID: %d", ErrUserNotFound, userID)
		}

		return models.User{}, fmt.Errorf(failedScanStr, err)
	}

	return u, nil
}

// FetchUserData получить базовую информацию о данных пользователя.
func (s *Storage) FetchUserData(ctx context.Context) ([]models.UserData, error) {
	const query = `SELECT id, type, mark, description FROM user_data WHERE user_id = $1`

	data := []models.UserData{}

	rows, err := s.pool.Query(ctx, query, ctx.Value(constants.KeyUserID))
	if err != nil {
		return []models.UserData{}, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var d models.UserData
		err = rows.Scan(&d.ID, &d.Type, &d.Mark, &d.Description)
		if err != nil {
			return []models.UserData{}, fmt.Errorf("failed to scan query: %w", err)
		}

		data = append(data, d)
	}

	rowsErr := rows.Err()
	if rowsErr != nil {
		return []models.UserData{}, fmt.Errorf("failed to read query: %w", err)
	}

	return data, nil
}

// AddUserData добавить данные пользователя.
func (s *Storage) AddUserData(
	ctx context.Context,
	encData []byte,
	mark string,
	description string,
	dataType string) (int, error) {
	const stmt = `
		INSERT INTO user_data (user_id, data, mark, description, type) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id
	`
	row := s.pool.QueryRow(ctx, stmt, ctx.Value(constants.KeyUserID), encData, mark, description, dataType)

	var id int

	err := row.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to scan a response row: %w", err)
	}

	return id, nil
}

// AddUserData получить данные пользователя.
func (s *Storage) GetUserData(ctx context.Context, id int, dataType string) ([]byte, string, string, error) {
	const query = `
		SELECT data, mark, description FROM user_data 
		WHERE user_id = $1 AND id = $2 AND type = $3 LIMIT 1
	`

	row := s.pool.QueryRow(ctx, query, ctx.Value(constants.KeyUserID), id, dataType)

	var data []byte
	var mark string
	var description string

	err := row.Scan(&data, &mark, &description)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, "", "", ErrUserDataNotFound
		}

		return nil, "", "", fmt.Errorf("failed to scan a response row: %w", err)
	}

	return data, mark, description, nil
}
