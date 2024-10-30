package get

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

func TestGetPasswordCmd(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	s := mocks.NewMockServicer(mockCtrl)

	passwordID := "1"
	data := models.Password{
		ID:          1,
		Login:       "test",
		Password:    "test",
		Mark:        "test",
		Description: "test",
	}
	expectedOutput := `{
  "login": "test",
  "password": "test",
  "mark": "test",
  "description": "test",
  "id": 1
}
`

	type getPassword struct {
		resp models.Password
		err  error
	}
	tests := []struct {
		name        string
		args        []string
		getPassword getPassword
		output      string
	}{
		{
			name: "get password success",
			args: []string{"get", "password", passwordID},
			getPassword: getPassword{
				resp: data,
				err:  nil,
			},
			output: expectedOutput,
		},
		{
			name: "get password failed",
			args: []string{"get", "password", passwordID},
			getPassword: getPassword{
				resp: models.Password{},
				err:  errors.New("some error"),
			},
			output: "Failed: some error",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s.EXPECT().GetPassword(passwordID).Times(1).Return(test.getPassword.resp, test.getPassword.err)

			cmd.RootCmd.SetArgs(test.args)

			var outBuf bytes.Buffer
			cmd.RootCmd.SetOutput(&outBuf)

			cmd.Execute(s)

			assert.Equal(t, test.output, outBuf.String())
		})
	}
}
