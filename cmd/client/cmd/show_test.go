package cmd

import (
	"bytes"
	"errors"
	"testing"

	"github.com/MihailSergeenkov/GophKeeper/cmd/client/cmd/mocks"
	"github.com/MihailSergeenkov/GophKeeper/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestShowCmd(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	s := mocks.NewMockServicer(mockCtrl)

	data := []models.UserData{
		{
			ID:          1,
			Type:        "text",
			Mark:        "test",
			Description: "test",
		},
	}
	expectedOutput := `[
  {
    "mark": "test",
    "description": "test",
    "type": "text",
    "id": 1
  }
]
`

	type syncData struct {
		count int
		err   error
	}
	type getData struct {
		count int
		resp  []models.UserData
	}
	tests := []struct {
		name     string
		args     []string
		syncData syncData
		getData  getData
		output   string
	}{
		{
			name: "show success without sync",
			args: []string{"show"},
			syncData: syncData{
				count: 0,
				err:   nil,
			},
			getData: getData{
				count: 1,
				resp:  data,
			},
			output: expectedOutput,
		},
		{
			name: "show success with sync",
			args: []string{"show", "-s"},
			syncData: syncData{
				count: 1,
				err:   nil,
			},
			getData: getData{
				count: 1,
				resp:  data,
			},
			output: expectedOutput,
		},
		{
			name: "show failed with failed sync",
			args: []string{"show", "-s"},
			syncData: syncData{
				count: 1,
				err:   errors.New("some error"),
			},
			getData: getData{
				count: 0,
				resp:  []models.UserData{},
			},
			output: "Failed: some error",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s.EXPECT().SyncData().Times(test.syncData.count).Return(test.syncData.err)
			s.EXPECT().GetData().Times(test.getData.count).Return(test.getData.resp)

			RootCmd.SetArgs(test.args)

			var outBuf bytes.Buffer
			RootCmd.SetOutput(&outBuf)

			Execute(s)

			assert.Equal(t, test.output, outBuf.String())
		})
	}
}
