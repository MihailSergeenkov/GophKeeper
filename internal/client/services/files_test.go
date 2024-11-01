package services

import (
	"errors"
	"net/http"
	"testing"

	"github.com/MihailSergeenkov/GophKeeper/internal/client/services/mocks"
	"github.com/MihailSergeenkov/GophKeeper/internal/models"
	"github.com/go-resty/resty/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddFile(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cfg := mocks.NewMockConfigurer(mockCtrl)
	r := mocks.NewMockRequester(mockCtrl)
	s := Init(cfg, r)

	url := "http://some/api"
	filePath := "test.txt"
	mark := "test"
	description := "test"

	type postResponse struct {
		resp *resty.Response
		err  error
	}
	type addData struct {
		count int
		err   error
	}
	tests := []struct {
		name         string
		postResponse postResponse
		addData      addData
		wantErr      bool
		errText      string
	}{
		{
			name: "add file success",
			postResponse: postResponse{
				resp: &resty.Response{
					RawResponse: &http.Response{StatusCode: http.StatusCreated},
				},
				err: nil,
			},
			addData: addData{
				count: 1,
				err:   nil,
			},
			wantErr: false,
			errText: "",
		},
		{
			name: "add file failed",
			postResponse: postResponse{
				resp: &resty.Response{
					RawResponse: &http.Response{StatusCode: http.StatusCreated},
				},
				err: nil,
			},
			addData: addData{
				count: 1,
				err:   errors.New("some error"),
			},
			wantErr: true,
			errText: "failed to dump data",
		},
		{
			name: "add file failed when response status not 201",
			postResponse: postResponse{
				resp: &resty.Response{
					RawResponse: &http.Response{StatusCode: http.StatusForbidden},
				},
				err: nil,
			},
			addData: addData{
				count: 0,
				err:   nil,
			},
			wantErr: true,
			errText: "response status",
		},
		{
			name: "add file failed when request failed",
			postResponse: postResponse{
				resp: nil,
				err:  errors.New("some error"),
			},
			addData: addData{
				count: 0,
				err:   nil,
			},
			wantErr: true,
			errText: "failed request",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg.EXPECT().GetToken().Times(1).Return("token")
			cfg.EXPECT().GetServerAPI().Times(1).Return(url)

			r.EXPECT().Post(url+"/user/files", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Times(1).Return(test.postResponse.resp, test.postResponse.err)

			cfg.EXPECT().AddData(gomock.Any()).Times(test.addData.count).Return(test.addData.err)

			err := s.AddFile(filePath, mark, description)

			if test.wantErr {
				require.Error(t, err)
				assert.ErrorContains(t, err, test.errText)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetFile(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cfg := mocks.NewMockConfigurer(mockCtrl)
	r := mocks.NewMockRequester(mockCtrl)
	s := Init(cfg, r)

	url := "http://some/api"
	fileMark := "test"
	dir := "."

	type getResponse struct {
		count int
		resp  *resty.Response
		err   error
	}
	tests := []struct {
		name        string
		data        map[string]models.UserData
		getResponse getResponse
		wantErr     bool
		errText     string
	}{
		{
			name: "file text success",
			data: map[string]models.UserData{
				fileMark: {
					ID: 1,
				},
			},
			getResponse: getResponse{
				count: 1,
				resp: &resty.Response{
					RawResponse: &http.Response{StatusCode: http.StatusOK},
				},
				err: nil,
			},
			wantErr: false,
			errText: "",
		},
		{
			name: "file not found",
			data: map[string]models.UserData{
				"test2": {
					ID: 2,
				},
			},
			getResponse: getResponse{
				count: 0,
				resp:  nil,
				err:   nil,
			},
			wantErr: true,
			errText: "file mark not found",
		},
		{
			name: "get file failed when response status not 200",
			data: map[string]models.UserData{
				fileMark: {
					ID: 1,
				},
			},
			getResponse: getResponse{
				count: 1,
				resp: &resty.Response{
					RawResponse: &http.Response{StatusCode: http.StatusForbidden},
				},
				err: nil,
			},
			wantErr: true,
			errText: "response status",
		},
		{
			name: "get file failed when request failed",
			data: map[string]models.UserData{
				fileMark: {
					ID: 1,
				},
			},
			getResponse: getResponse{
				count: 1,
				resp:  nil,
				err:   errors.New("some error"),
			},
			wantErr: true,
			errText: "failed request",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg.EXPECT().GetData().Times(1).Return(test.data)
			cfg.EXPECT().GetToken().Times(test.getResponse.count).Return("token")
			cfg.EXPECT().GetServerAPI().Times(test.getResponse.count).Return(url)

			r.EXPECT().Get(url+"/user/files/{fileMark}", gomock.Any(), gomock.Any(), gomock.Any()).
				Times(test.getResponse.count).Return(test.getResponse.resp, test.getResponse.err)

			err := s.GetFile(fileMark, dir)

			if test.wantErr {
				require.Error(t, err)
				assert.ErrorContains(t, err, test.errText)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
