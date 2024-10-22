package storage

import (
	"context"
	"errors"
	"testing"

	"github.com/MihailSergeenkov/GophKeeper/internal/server/constants"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/storage/mocks"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestPing(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	pool := mocks.NewMockDBPooler(mockCtrl)
	logger := zap.NewNop()
	storage := Storage{
		pool:   pool,
		logger: logger,
	}
	ctx := context.Background()

	tests := []struct {
		name    string
		wantErr bool
		err     error
	}{
		{
			name:    "success ping",
			wantErr: false,
			err:     nil,
		},
		{
			name:    "failed ping",
			wantErr: true,
			err:     errors.New("some error"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pool.EXPECT().Ping(ctx).Times(1).Return(test.err)

			err := storage.Ping(ctx)

			if test.wantErr {
				require.Error(t, err)
				assert.ErrorContains(t, err, "failed to ping DB")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestClose(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	pool := mocks.NewMockDBPooler(mockCtrl)
	logger := zap.NewNop()
	storage := Storage{
		pool:   pool,
		logger: logger,
	}

	t.Run("success close", func(t *testing.T) {
		pool.EXPECT().Close().Times(1)

		err := storage.Close()
		require.NoError(t, err)
	})
}

func TestAddUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	pool := mocks.NewMockDBPooler(mockCtrl)
	logger := zap.NewNop()
	storage := Storage{
		pool:   pool,
		logger: logger,
	}
	ctx := context.Background()
	stmt := `INSERT INTO users (login, password) VALUES ($1, $2)`

	login := "login"
	password := []byte("password")

	tests := []struct {
		name    string
		wantErr bool
		err     error
	}{
		{
			name:    "success add",
			wantErr: false,
			err:     nil,
		},
		{
			name:    "failed add",
			wantErr: true,
			err:     errors.New("some error"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pool.EXPECT().Exec(ctx, stmt, login, password).Times(1).Return(pgconn.CommandTag{}, test.err)

			err := storage.AddUser(ctx, login, password)

			if test.wantErr {
				require.Error(t, err)
				assert.ErrorContains(t, err, "failed to execute add user query")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetUserByLogin(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	pool := mocks.NewMockDBPooler(mockCtrl)
	logger := zap.NewNop()
	storage := Storage{
		pool:   pool,
		logger: logger,
	}
	ctx := context.Background()
	stmt := `SELECT id, login, password FROM users WHERE login = $1 LIMIT 1`

	row := mocks.NewMockRow(mockCtrl)
	login := "login"

	tests := []struct {
		name    string
		wantErr bool
		errText string
		rowErr  error
	}{
		{
			name:    "success get",
			wantErr: false,
			errText: "",
			rowErr:  nil,
		},
		{
			name:    "failed read row",
			wantErr: true,
			errText: "failed to scan a response row",
			rowErr:  errors.New("some error"),
		},
		{
			name:    "not found row",
			wantErr: true,
			errText: "user not found",
			rowErr:  pgx.ErrNoRows,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pool.EXPECT().QueryRow(ctx, stmt, login).Times(1).Return(row)

			row.EXPECT().Scan(gomock.Any()).Times(1).Return(test.rowErr)

			_, err := storage.GetUserByLogin(ctx, login)

			if test.wantErr {
				require.Error(t, err)
				assert.ErrorContains(t, err, test.errText)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetUserByID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	pool := mocks.NewMockDBPooler(mockCtrl)
	logger := zap.NewNop()
	storage := Storage{
		pool:   pool,
		logger: logger,
	}
	ctx := context.Background()
	stmt := `SELECT id, login, password FROM users WHERE id = $1 LIMIT 1`

	row := mocks.NewMockRow(mockCtrl)
	userID := 1

	tests := []struct {
		name    string
		wantErr bool
		errText string
		rowErr  error
	}{
		{
			name:    "success get",
			wantErr: false,
			errText: "",
			rowErr:  nil,
		},
		{
			name:    "failed read row",
			wantErr: true,
			errText: "failed to scan a response row",
			rowErr:  errors.New("some error"),
		},
		{
			name:    "not found row",
			wantErr: true,
			errText: "user not found",
			rowErr:  pgx.ErrNoRows,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pool.EXPECT().QueryRow(ctx, stmt, userID).Times(1).Return(row)

			row.EXPECT().Scan(gomock.Any()).Times(1).Return(test.rowErr)

			_, err := storage.GetUserByID(ctx, userID)

			if test.wantErr {
				require.Error(t, err)
				assert.ErrorContains(t, err, test.errText)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestFetchUserData(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	pool := mocks.NewMockDBPooler(mockCtrl)
	logger := zap.NewNop()
	storage := Storage{
		pool:   pool,
		logger: logger,
	}
	currentUserID := "some_id"
	ctx := context.WithValue(context.Background(), constants.KeyUserID, currentUserID)
	stmt := `SELECT id, type, mark, description FROM user_data WHERE user_id = $1`

	rows := mocks.NewMockRows(mockCtrl)

	tests := []struct {
		name    string
		wantErr bool
		errText string
		rowsErr error
	}{
		{
			name:    "success fetch",
			wantErr: false,
			errText: "",
			rowsErr: nil,
		},
		{
			name:    "failed read rows",
			wantErr: true,
			errText: "failed to read query",
			rowsErr: errors.New("some error"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pool.EXPECT().Query(ctx, stmt, currentUserID).Times(1).Return(rows, nil)

			rows.EXPECT().Close().Times(1)
			rows.EXPECT().Next().Times(1).Return(false)
			rows.EXPECT().Err().Times(1).Return(test.rowsErr)

			_, err := storage.FetchUserData(ctx)

			if test.wantErr {
				require.Error(t, err)
				assert.ErrorContains(t, err, test.errText)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestFetchUserData_Failed(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	pool := mocks.NewMockDBPooler(mockCtrl)
	logger := zap.NewNop()
	storage := Storage{
		pool:   pool,
		logger: logger,
	}
	currentUserID := "some_id"
	ctx := context.WithValue(context.Background(), constants.KeyUserID, currentUserID)
	stmt := `SELECT id, type, mark, description FROM user_data WHERE user_id = $1`

	rows := mocks.NewMockRows(mockCtrl)
	someErr := errors.New("some error")

	t.Run("failed fetch", func(t *testing.T) {
		pool.EXPECT().Query(ctx, stmt, currentUserID).Times(1).Return(rows, someErr)

		_, err := storage.FetchUserData(ctx)

		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to execute query")
	})
}

func TestAddUserData(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	pool := mocks.NewMockDBPooler(mockCtrl)
	logger := zap.NewNop()
	storage := Storage{
		pool:   pool,
		logger: logger,
	}
	currentUserID := "some_id"
	ctx := context.WithValue(context.Background(), constants.KeyUserID, currentUserID)
	const stmt = `
		INSERT INTO user_data (user_id, data, mark, description, type) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id
	`

	row := mocks.NewMockRow(mockCtrl)

	encData := []byte("some data")
	mark := "test"
	description := "test"
	dataType := "files"

	tests := []struct {
		name    string
		wantErr bool
		rowErr  error
	}{
		{
			name:    "success add user data",
			wantErr: false,
			rowErr:  nil,
		},
		{
			name:    "failed read row",
			wantErr: true,
			rowErr:  errors.New("some error"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pool.EXPECT().QueryRow(ctx, stmt, currentUserID, encData, mark, description, dataType).Times(1).Return(row)

			row.EXPECT().Scan(gomock.Any()).Times(1).Return(test.rowErr)

			_, err := storage.AddUserData(ctx, encData, mark, description, dataType)

			if test.wantErr {
				require.Error(t, err)
				assert.ErrorContains(t, err, "failed to scan a response row", "some error")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetUserData(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	pool := mocks.NewMockDBPooler(mockCtrl)
	logger := zap.NewNop()
	storage := Storage{
		pool:   pool,
		logger: logger,
	}
	currentUserID := "some_id"
	ctx := context.WithValue(context.Background(), constants.KeyUserID, currentUserID)
	stmt := `
		SELECT data, mark, description FROM user_data 
		WHERE user_id = $1 AND id = $2 AND type = $3 LIMIT 1
	`

	row := mocks.NewMockRow(mockCtrl)
	userDataID := 1
	dataType := "files"

	tests := []struct {
		name    string
		wantErr bool
		errText string
		rowErr  error
	}{
		{
			name:    "success get",
			wantErr: false,
			errText: "",
			rowErr:  nil,
		},
		{
			name:    "failed read row",
			wantErr: true,
			errText: "failed to scan a response row",
			rowErr:  errors.New("some error"),
		},
		{
			name:    "not found row",
			wantErr: true,
			errText: "user data not found",
			rowErr:  pgx.ErrNoRows,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pool.EXPECT().QueryRow(ctx, stmt, currentUserID, userDataID, dataType).Times(1).Return(row)

			row.EXPECT().Scan(gomock.Any()).Times(1).Return(test.rowErr)

			_, _, _, err := storage.GetUserData(ctx, userDataID, dataType)

			if test.wantErr {
				require.Error(t, err)
				assert.ErrorContains(t, err, test.errText)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
