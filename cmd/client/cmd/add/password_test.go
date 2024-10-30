package add

import (
	"bytes"
	"errors"
	"testing"

	"github.com/MihailSergeenkov/GophKeeper/cmd/client/cmd"
	"github.com/MihailSergeenkov/GophKeeper/cmd/client/cmd/mocks"
	"github.com/MihailSergeenkov/GophKeeper/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAddPasswordCmd(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	s := mocks.NewMockServicer(mockCtrl)

	req := models.AddPasswordRequest{
		Login:       "test",
		Password:    "test",
		Mark:        "test",
		Description: "test",
	}

	type addPassword struct {
		err error
	}
	tests := []struct {
		name        string
		args        []string
		addPassword addPassword
		output      string
	}{
		{
			name: "add password success",
			args: []string{"add", "password", "-l", "test", "-p", "test", "-m", "test", "-d", "test"},
			addPassword: addPassword{
				err: nil,
			},
			output: "Add password OK\n",
		},
		{
			name: "add password failed",
			args: []string{"add", "password", "-l", "test", "-p", "test", "-m", "test", "-d", "test"},
			addPassword: addPassword{
				err: errors.New("some error"),
			},
			output: "Failed: some error",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s.EXPECT().AddPassword(req).Times(1).Return(test.addPassword.err)

			cmd.RootCmd.SetArgs(test.args)

			var outBuf bytes.Buffer
			cmd.RootCmd.SetOutput(&outBuf)

			cmd.Execute(s)

			assert.Equal(t, test.output, outBuf.String())
		})
	}
}
