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

func TestAddText(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cfg := mocks.NewMockConfigurer(mockCtrl)
	r := mocks.NewMockRequester(mockCtrl)
	s := Init(cfg, r)

	url := "http://some/api"
	req := models.AddTextRequest{}

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
			name: "add text success",
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
			name: "add text failed",
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
			name: "add text failed when response status not 201",
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
			name: "add text failed when request failed",
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

			r.EXPECT().Post(url+"/user/texts", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Times(1).Return(test.postResponse.resp, test.postResponse.err)

			cfg.EXPECT().AddData(gomock.Any()).Times(test.addData.count).Return(test.addData.err)

			err := s.AddText(req)

			if test.wantErr {
				require.Error(t, err)
				assert.ErrorContains(t, err, test.errText)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetText(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cfg := mocks.NewMockConfigurer(mockCtrl)
	r := mocks.NewMockRequester(mockCtrl)
	s := Init(cfg, r)

	url := "http://some/api"
	textID := "1"

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
			name: "get text success",
			data: map[string]models.UserData{
				"1": {
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
			name: "text not found",
			data: map[string]models.UserData{
				"2": {
					ID: 2,
				},
			},
			getResponse: getResponse{
				count: 0,
				resp:  nil,
				err:   nil,
			},
			wantErr: true,
			errText: "text id not found",
		},
		{
			name: "get text failed when response status not 200",
			data: map[string]models.UserData{
				"1": {
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
			name: "get text failed when request failed",
			data: map[string]models.UserData{
				"1": {
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

			r.EXPECT().Get(url+"/user/texts/{id}", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Times(test.getResponse.count).Return(test.getResponse.resp, test.getResponse.err)

			_, err := s.GetText(textID)

			if test.wantErr {
				require.Error(t, err)
				assert.ErrorContains(t, err, test.errText)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
