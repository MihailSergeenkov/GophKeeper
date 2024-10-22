package services

import (
	"context"
	"errors"
	"testing"

	"github.com/MihailSergeenkov/GophKeeper/internal/server/config"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/models"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/services/mocks"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/storage"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestRegisterUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	store := mocks.NewMockStorager(mockCtrl)
	fs := mocks.NewMockFileStorager(mockCtrl)
	crypter := mocks.NewMockCrypter(mockCtrl)
	settings := config.Settings{}
	s := NewServices(store, fs, crypter, &settings)

	ctx := context.Background()
	errSome := errors.New("some error")

	type arg struct {
		req models.RegisterUserRequest
	}

	type sResponse struct {
		err error
	}

	type want struct {
		err error
	}

	tests := []struct {
		name      string
		arg       arg
		sResponse sResponse
		want      want
	}{
		{
			name: "register user success",
			arg: arg{
				req: models.RegisterUserRequest{
					Login:    "test",
					Password: "test",
				},
			},
			sResponse: sResponse{
				err: nil,
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "register user failed",
			arg: arg{
				req: models.RegisterUserRequest{
					Login:    "test",
					Password: "test",
				},
			},
			sResponse: sResponse{
				err: errSome,
			},
			want: want{
				err: errSome,
			},
		},
		{
			name: "register user failed with constraint error",
			arg: arg{
				req: models.RegisterUserRequest{
					Login:    "test",
					Password: "test",
				},
			},
			sResponse: sResponse{
				err: &pgconn.PgError{Code: pgerrcode.UniqueViolation},
			},
			want: want{
				err: ErrUserLoginExist,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store.EXPECT().AddUser(ctx, test.arg.req.Login, gomock.Any()).Times(1).Return(test.sResponse.err)

			err := s.RegisterUser(ctx, test.arg.req)

			if test.sResponse.err == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.ErrorContains(t, err, test.want.err.Error())
			}
		})
	}
}

func TestValidationFailedRegisterUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	store := mocks.NewMockStorager(mockCtrl)
	fs := mocks.NewMockFileStorager(mockCtrl)
	crypter := mocks.NewMockCrypter(mockCtrl)
	settings := config.Settings{}
	s := NewServices(store, fs, crypter, &settings)

	ctx := context.Background()

	type arg struct {
		req models.RegisterUserRequest
	}

	type want struct {
		err error
	}

	tests := []struct {
		name string
		arg  arg
		want want
	}{
		{
			name: "when login empty",
			arg: arg{
				req: models.RegisterUserRequest{
					Login:    "",
					Password: "test",
				},
			},
			want: want{
				err: ErrUserValidationFields,
			},
		},
		{
			name: "when password empty",
			arg: arg{
				req: models.RegisterUserRequest{
					Login:    "test",
					Password: "",
				},
			},
			want: want{
				err: ErrUserValidationFields,
			},
		},
		{
			name: "when password very big",
			arg: arg{
				req: models.RegisterUserRequest{
					Login:    "test",
					Password: "testtesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttest",
				},
			},
			want: want{
				err: bcrypt.ErrPasswordTooLong,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store.EXPECT().AddUser(ctx, gomock.Any(), gomock.Any()).Times(0)

			err := s.RegisterUser(ctx, test.arg.req)

			require.Error(t, err)
			assert.ErrorContains(t, err, test.want.err.Error())
		})
	}
}

func TestCreateUserToken(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	store := mocks.NewMockStorager(mockCtrl)
	fs := mocks.NewMockFileStorager(mockCtrl)
	crypter := mocks.NewMockCrypter(mockCtrl)
	settings := config.Settings{}
	s := NewServices(store, fs, crypter, &settings)
	ctx := context.Background()
	errSome := errors.New("some error")

	type arg struct {
		req models.CreateUserTokenRequest
	}

	type sResponse struct {
		user models.User
		err  error
	}

	type want struct {
		res models.CreateUserTokenResponse
		err error
	}

	tests := []struct {
		name      string
		arg       arg
		sResponse sResponse
		wantErr   bool
		want      want
	}{
		{
			name: "create user token success",
			arg: arg{
				req: models.CreateUserTokenRequest{
					Login:    "test",
					Password: "test",
				},
			},
			sResponse: sResponse{
				user: models.User{
					ID:       1,
					Login:    "test",
					Password: []byte("$2a$10$eqoHdZljD4bk/zPKKGAPre6Mmq2mj8XxSrjF4SpavRy.pT/uxijYa"),
				},
				err: nil,
			},
			wantErr: false,
			want: want{
				res: models.CreateUserTokenResponse{
					AuthToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOjF9.06lbO3Sb1wyS45SCYsUxwrUyon5u6l1bnCbzwp83wbI",
				},
				err: nil,
			},
		},
		{
			name: "create user token failed",
			arg: arg{
				req: models.CreateUserTokenRequest{
					Login:    "test",
					Password: "test",
				},
			},
			sResponse: sResponse{
				user: models.User{},
				err:  errSome,
			},
			wantErr: true,
			want: want{
				res: models.CreateUserTokenResponse{},
				err: errSome,
			},
		},
		{
			name: "when user not found",
			arg: arg{
				req: models.CreateUserTokenRequest{
					Login:    "test",
					Password: "test",
				},
			},
			sResponse: sResponse{
				user: models.User{},
				err:  storage.ErrUserNotFound,
			},
			wantErr: true,
			want: want{
				res: models.CreateUserTokenResponse{},
				err: ErrUserLoginCreds,
			},
		},
		{
			name: "when incorrect password",
			arg: arg{
				req: models.CreateUserTokenRequest{
					Login:    "test",
					Password: "test2",
				},
			},
			sResponse: sResponse{
				user: models.User{
					ID:       1,
					Login:    "test",
					Password: []byte("$2a$10$eqoHdZljD4bk/zPKKGAPre6Mmq2mj8XxSrjF4SpavRy.pT/uxijYa"),
				},
				err: nil,
			},
			wantErr: true,
			want: want{
				res: models.CreateUserTokenResponse{},
				err: ErrUserLoginCreds,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store.EXPECT().GetUserByLogin(ctx, test.arg.req.Login).Times(1).Return(test.sResponse.user, test.sResponse.err)

			result, err := s.CreateUserToken(ctx, test.arg.req)

			assert.Equal(t, test.want.res, result)

			if test.wantErr {
				require.Error(t, err)
				assert.ErrorContains(t, err, test.want.err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
