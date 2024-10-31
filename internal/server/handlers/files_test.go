package handlers

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MihailSergeenkov/GophKeeper/internal/models"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/config"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/handlers/mocks"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/routes"
	rMocks "github.com/MihailSergeenkov/GophKeeper/internal/server/routes/mocks"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/services"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestAddFile(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	s := mocks.NewMockServicer(mockCtrl)
	l := mocks.NewMockLogger(mockCtrl)
	handlers := NewHandlers(s, l)

	path := "./testdata/test.txt"
	file, err := os.Open(path)
	require.NoError(t, err)
	defer func() {
		if err := file.Close(); err != nil {
			t.Logf("failed to close file %v", err)
		}
	}()

	type serviceResponse struct {
		id  int
		err error
	}
	type want struct {
		code          int
		body          string
		errorLogTimes int
		log           string
	}
	tests := []struct {
		name            string
		serviceResponse serviceResponse
		want            want
	}{
		{
			name: "add file success",
			serviceResponse: serviceResponse{
				id:  1,
				err: nil,
			},
			want: want{
				code:          http.StatusCreated,
				body:          "{\"id\":1}\n",
				errorLogTimes: 0,
				log:           "",
			},
		},
		{
			name: "add file failed",
			serviceResponse: serviceResponse{
				id:  0,
				err: errors.New("some error"),
			},
			want: want{
				code:          http.StatusInternalServerError,
				body:          "",
				errorLogTimes: 1,
				log:           "failed to add file",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s.EXPECT().
				AddFile(gomock.Any(), gomock.Any()).
				Times(1).
				Return(test.serviceResponse.id, test.serviceResponse.err)

			l.EXPECT().Error(test.want.log, zap.Error(test.serviceResponse.err)).Times(test.want.errorLogTimes)

			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			defer func() {
				if err := writer.Close(); err != nil {
					t.Logf("failed to close writer %v", err)
				}
			}()
			part, err := writer.CreateFormFile("file", filepath.Base(path))

			require.NoError(t, err)
			_, err = io.Copy(part, file)
			require.NoError(t, err)

			err = writer.WriteField("mark", "test")
			require.NoError(t, err)
			err = writer.WriteField("description", "test")
			require.NoError(t, err)

			err = writer.Close()
			require.NoError(t, err)

			request := httptest.NewRequest(http.MethodPost, "/api/user/files", body)

			request.Header.Add("Content-Type", writer.FormDataContentType())

			w := httptest.NewRecorder()
			handlers.AddFile()(w, request)

			res := w.Result()
			defer closeBody(t, res)

			assert.Equal(t, test.want.code, res.StatusCode)

			if http.StatusCreated == res.StatusCode {
				resBody, err := io.ReadAll(res.Body)
				require.NoError(t, err)
				assert.Equal(t, test.want.body, string(resBody))
			}
		})
	}
}

func TestAddFileFailedReadContent(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	s := mocks.NewMockServicer(mockCtrl)
	l := mocks.NewMockLogger(mockCtrl)
	handlers := NewHandlers(s, l)

	path := "./testdata/test.txt"
	file, err := os.Open(path)
	require.NoError(t, err)
	defer func() {
		if err := file.Close(); err != nil {
			t.Logf("failed to close file %v", err)
		}
	}()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	defer func() {
		if err := writer.Close(); err != nil {
			t.Logf("failed to close writer %v", err)
		}
	}()

	type want struct {
		code        int
		contentType string
		log         string
	}
	tests := []struct {
		name    string
		prepare func()
		want    want
	}{
		{
			name:    "when another content type",
			prepare: func() {},
			want: want{
				code:        http.StatusInternalServerError,
				contentType: JSONContentType,
				log:         "failed to parse file",
			},
		},
		{
			name: "when file field absent",
			prepare: func() {
				part, err := writer.CreateFormFile("file2s", filepath.Base(path))

				require.NoError(t, err)
				_, err = io.Copy(part, file)
				require.NoError(t, err)

				err = writer.WriteField("mark", "test")
				require.NoError(t, err)
				err = writer.WriteField("description", "test")
				require.NoError(t, err)

				err = writer.Close()
				require.NoError(t, err)
			},
			want: want{
				code:        http.StatusBadRequest,
				contentType: writer.FormDataContentType(),
				log:         "failed to retrieving the file",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s.EXPECT().AddFile(gomock.Any(), gomock.Any()).Times(0)
			l.EXPECT().Error(test.want.log, gomock.Any()).Times(1)

			test.prepare()

			request := httptest.NewRequest(http.MethodPost, "/api/user/files", body)

			request.Header.Add("Content-Type", test.want.contentType)

			w := httptest.NewRecorder()
			handlers.AddFile()(w, request)

			res := w.Result()
			defer closeBody(t, res)

			assert.Equal(t, test.want.code, res.StatusCode)
		})
	}
}

func TestGetFile(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	s := mocks.NewMockServicer(mockCtrl)
	l := mocks.NewMockLogger(mockCtrl)
	handlers := NewHandlers(s, l)

	settings, err := config.Setup(false)
	require.NoError(t, err)
	storage := rMocks.NewMockStorager(mockCtrl)

	fileMark := "test"

	r := routes.NewRouter(handlers, settings, zap.NewNop(), storage)
	ts := httptest.NewServer(r)
	defer ts.Close()

	type serviceResponse struct {
		res models.File
		err error
	}
	type want struct {
		code          int
		errorLogTimes int
		log           string
	}
	tests := []struct {
		name            string
		serviceResponse serviceResponse
		want            want
	}{
		{
			name: "get file success",
			serviceResponse: serviceResponse{
				res: models.File{
					File: io.NopCloser(strings.NewReader("some data")),
				},
				err: nil,
			},
			want: want{
				code:          http.StatusOK,
				errorLogTimes: 0,
				log:           "",
			},
		},
		{
			name: "file no found",
			serviceResponse: serviceResponse{
				res: models.File{},
				err: services.ErrNotFound,
			},
			want: want{
				code:          http.StatusNotFound,
				errorLogTimes: 0,
				log:           "",
			},
		},
		{
			name: "get file failed with some error",
			serviceResponse: serviceResponse{
				res: models.File{},
				err: errors.New("some error"),
			},
			want: want{
				code:          http.StatusInternalServerError,
				errorLogTimes: 1,
				log:           "failed to get file",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s.EXPECT().GetFile(gomock.Any(), fileMark).Times(1).
				Return(test.serviceResponse.res, test.serviceResponse.err)

			l.EXPECT().Error(test.want.log, zap.Error(test.serviceResponse.err)).Times(test.want.errorLogTimes)
			storage.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Times(1)

			res, _ := testGetRequest(t, ts, "/api/user/files/test")
			closeBody(t, res)

			assert.Equal(t, test.want.code, res.StatusCode)
		})
	}
}
