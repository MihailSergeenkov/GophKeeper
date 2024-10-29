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

func TestSuncData(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cfg := mocks.NewMockConfigurer(mockCtrl)
	r := mocks.NewMockRequester(mockCtrl)
	s := Init(cfg, r)

	url := "http://some/api"

	type getResponse struct {
		resp *resty.Response
		err  error
	}
	type updateData struct {
		count int
		err   error
	}
	tests := []struct {
		name        string
		getResponse getResponse
		updateData  updateData
		wantErr     bool
		errText     string
	}{
		{
			name: "sync data success",
			getResponse: getResponse{
				resp: &resty.Response{
					RawResponse: &http.Response{StatusCode: http.StatusOK},
				},
				err: nil,
			},
			updateData: updateData{
				count: 1,
				err:   nil,
			},
			wantErr: false,
			errText: "",
		},
		{
			name: "sync data failed",
			getResponse: getResponse{
				resp: &resty.Response{
					RawResponse: &http.Response{StatusCode: http.StatusOK},
				},
				err: nil,
			},
			updateData: updateData{
				count: 1,
				err:   errors.New("some error"),
			},
			wantErr: true,
			errText: "failed to update data",
		},
		{
			name: "sync data failed when response status not 200 or 204",
			getResponse: getResponse{
				resp: &resty.Response{
					RawResponse: &http.Response{StatusCode: http.StatusForbidden},
				},
				err: nil,
			},
			updateData: updateData{
				count: 0,
				err:   nil,
			},
			wantErr: true,
			errText: "response status",
		},
		{
			name: "sync data failed when request failed",
			getResponse: getResponse{
				resp: nil,
				err:  errors.New("some error"),
			},
			updateData: updateData{
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

			r.EXPECT().Get(url+"/user/data", gomock.Any(), gomock.Any(), gomock.Any()).
				Times(1).Return(test.getResponse.resp, test.getResponse.err)

			cfg.EXPECT().UpdateData(gomock.Any()).Times(test.updateData.count).Return(test.updateData.err)

			err := s.SyncData()

			if test.wantErr {
				require.Error(t, err)
				assert.ErrorContains(t, err, test.errText)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetData(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cfg := mocks.NewMockConfigurer(mockCtrl)
	r := mocks.NewMockRequester(mockCtrl)

	s := Init(cfg, r)

	data := map[string]models.UserData{
		"1": {
			ID: 1,
		},
	}

	cfg.EXPECT().GetData().Times(1).Return(data)

	t.Run("get data", func(t *testing.T) {
		result := s.GetData()

		assert.Equal(t, data["1"], result[0])
	})
}
