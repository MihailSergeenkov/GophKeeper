package services

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/MihailSergeenkov/GophKeeper/internal/models"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/config"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/services/mocks"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/storage"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddFile(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	store := mocks.NewMockStorager(mockCtrl)
	fs := mocks.NewMockFileStorager(mockCtrl)
	crypter := mocks.NewMockCrypter(mockCtrl)
	settings := config.Settings{}
	s := NewServices(store, fs, crypter, &settings)

	req := models.AddFileRequest{
		File:        strings.NewReader("test"),
		FileName:    "test",
		FileSize:    int64(100),
		Mark:        "test",
		Description: "test",
	}

	ctx := context.Background()
	dataType := "file"
	encData := []byte("some data")

	type sResponse struct {
		id  int
		err error
	}
	tests := []struct {
		name      string
		sResponse sResponse
		wantErr   bool
	}{
		{
			name: "add file success",
			sResponse: sResponse{
				id:  1,
				err: nil,
			},
			wantErr: false,
		},
		{
			name: "add file failed",
			sResponse: sResponse{
				id:  0,
				err: errors.New("some error"),
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fs.EXPECT().AddFile(ctx, req.File, req.FileName, req.FileSize).Times(1).Return(nil)
			crypter.EXPECT().EncryptData(gomock.Any()).Times(1).Return(encData)
			store.EXPECT().
				AddUserData(ctx, encData, req.Mark, req.Description, dataType).
				Times(1).Return(test.sResponse.id, test.sResponse.err)

			id, err := s.AddFile(ctx, req)

			if test.wantErr {
				require.Error(t, err)
				assert.ErrorContains(t, err, "failed to add user data", "some error")
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.sResponse.id, id)
			}
		})
	}
}

func TestAddFileFileStorageFailed(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	store := mocks.NewMockStorager(mockCtrl)
	fs := mocks.NewMockFileStorager(mockCtrl)
	crypter := mocks.NewMockCrypter(mockCtrl)
	settings := config.Settings{}
	s := NewServices(store, fs, crypter, &settings)

	ctx := context.Background()
	req := models.AddFileRequest{
		File:        strings.NewReader("test"),
		FileName:    "test",
		FileSize:    int64(100),
		Mark:        "test",
		Description: "test",
	}

	someErr := errors.New("some error")

	t.Run("file storage failed", func(t *testing.T) {
		fs.EXPECT().AddFile(ctx, req.File, req.FileName, req.FileSize).Times(1).Return(someErr)
		crypter.EXPECT().EncryptData(gomock.Any()).Times(0)
		store.EXPECT().AddUserData(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		_, err := s.AddFile(ctx, req)

		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to add file to filestorage", someErr.Error())
	})
}

func TestAddFileValidationFailed(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	store := mocks.NewMockStorager(mockCtrl)
	fs := mocks.NewMockFileStorager(mockCtrl)
	crypter := mocks.NewMockCrypter(mockCtrl)
	settings := config.Settings{}
	s := NewServices(store, fs, crypter, &settings)

	ctx := context.Background()

	type arg struct {
		req models.AddFileRequest
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
			name: "when user mark very big",
			arg: arg{
				req: models.AddFileRequest{
					File:        strings.NewReader("test"),
					FileName:    "test",
					FileSize:    int64(100),
					Mark:        "testtesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttest",
					Description: "test",
				},
			},
			want: want{
				err: ErrUserMarkIsTooBig,
			},
		},
		{
			name: "when user description very big",
			arg: arg{
				req: models.AddFileRequest{
					File:        strings.NewReader("test"),
					FileName:    "test",
					FileSize:    int64(100),
					Mark:        "test",
					Description: "testtesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttest",
				},
			},
			want: want{
				err: ErrUserDescriptionIsTooBig,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fs.EXPECT().AddFile(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			crypter.EXPECT().EncryptData(gomock.Any()).Times(0)
			store.EXPECT().AddUserData(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

			_, err := s.AddFile(ctx, test.arg.req)

			require.Error(t, err)
			assert.ErrorContains(t, err, test.want.err.Error())
		})
	}
}

func TestGetFile(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	store := mocks.NewMockStorager(mockCtrl)
	fs := mocks.NewMockFileStorager(mockCtrl)
	crypter := mocks.NewMockCrypter(mockCtrl)
	settings := config.Settings{}
	s := NewServices(store, fs, crypter, &settings)

	ctx := context.Background()
	userDataID := 1
	decData := []byte("some data")
	mark := "test"
	description := "test"
	dataType := "file"

	jsonData := []byte(`{"file_name":"test"}`)

	type fsResponse struct {
		file io.ReadCloser
		err  error
	}
	tests := []struct {
		name       string
		fsResponse fsResponse
		wantErr    bool
	}{
		{
			name: "get user data success",
			fsResponse: fsResponse{
				file: io.NopCloser(strings.NewReader("some data")),
				err:  nil,
			},
			wantErr: false,
		},
		{
			name: "when get file from fileserver failed",
			fsResponse: fsResponse{
				file: nil,
				err:  errors.New("some error"),
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store.EXPECT().GetUserData(ctx, userDataID, dataType).
				Times(1).Return(decData, mark, description, nil)

			crypter.EXPECT().DecryptData(decData).Times(1).Return(jsonData, nil)
			fs.EXPECT().GetFile(ctx, "test").Times(1).Return(test.fsResponse.file, test.fsResponse.err)

			resp, err := s.GetFile(ctx, userDataID)

			if test.wantErr {
				require.Error(t, err)
				assert.ErrorContains(t, err, "failed to get file from filestorage")
			} else {
				require.NoError(t, err)
				assert.Equal(t, models.File{File: test.fsResponse.file}, resp)
			}
		})
	}
}

func TestGetFileDecryptDataFailed(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	store := mocks.NewMockStorager(mockCtrl)
	fs := mocks.NewMockFileStorager(mockCtrl)
	crypter := mocks.NewMockCrypter(mockCtrl)
	settings := config.Settings{}
	s := NewServices(store, fs, crypter, &settings)

	ctx := context.Background()
	userDataID := 1
	dataType := "file"
	decData := []byte("some data")

	type cResponse struct {
		jsonData []byte
		err      error
	}
	tests := []struct {
		name      string
		cResponse cResponse
		wantErr   bool
		errText   string
	}{
		{
			name: "when decrypt data failed",
			cResponse: cResponse{
				jsonData: nil,
				err:      errors.New("some error"),
			},
			wantErr: true,
			errText: "failed to decrypt data",
		},
		{
			name: "when generate user data failed",
			cResponse: cResponse{
				jsonData: []byte(`test`),
				err:      nil,
			},
			wantErr: true,
			errText: "failed to generate data",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store.EXPECT().GetUserData(ctx, userDataID, dataType).Times(1).Return(decData, "", "", nil)
			crypter.EXPECT().DecryptData(gomock.Any()).Times(1).Return(test.cResponse.jsonData, test.cResponse.err)
			fs.EXPECT().GetFile(ctx, gomock.Any()).Times(0)

			_, err := s.GetFile(ctx, userDataID)

			require.Error(t, err)
			assert.ErrorContains(t, err, test.errText)
		})
	}
}

func TestGetFileFailedStorage(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	store := mocks.NewMockStorager(mockCtrl)
	fs := mocks.NewMockFileStorager(mockCtrl)
	crypter := mocks.NewMockCrypter(mockCtrl)
	settings := config.Settings{}
	s := NewServices(store, fs, crypter, &settings)

	ctx := context.Background()
	userDataID := 1
	dataType := "file"

	type sResponse struct {
		err error
	}
	tests := []struct {
		name      string
		sResponse sResponse
		errText   string
	}{
		{
			name: "failed to get user data",
			sResponse: sResponse{
				err: errors.New("some error"),
			},
			errText: "failed to get user data",
		},
		{
			name: "when user data not found",
			sResponse: sResponse{
				err: storage.ErrUserDataNotFound,
			},
			errText: "requested data no found",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store.EXPECT().GetUserData(ctx, userDataID, dataType).Times(1).Return([]byte{}, "", "", test.sResponse.err)
			crypter.EXPECT().DecryptData(gomock.Any()).Times(0)
			fs.EXPECT().GetFile(ctx, gomock.Any()).Times(0)

			_, err := s.GetFile(ctx, userDataID)

			require.Error(t, err)
			assert.ErrorContains(t, err, test.errText, test.sResponse.err.Error())
		})
	}
}
