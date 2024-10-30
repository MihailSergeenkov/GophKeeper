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

func TestRegisterCmd(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	s := mocks.NewMockServicer(mockCtrl)

	req := models.RegisterUserRequest{
		Login:    "qwe",
		Password: "123",
	}

	type registerUser struct {
		err error
	}
	tests := []struct {
		name         string
		args         []string
		registerUser registerUser
		output       string
	}{
		{
			name: "register success",
			args: []string{"register", "-l", "qwe", "-p", "123"},
			registerUser: registerUser{
				err: nil,
			},
			output: "Register OK\n",
		},
		{
			name: "register failed with failed call login service",
			args: []string{"register", "-l", "qwe", "-p", "123"},
			registerUser: registerUser{
				err: errors.New("some error"),
			},
			output: "Failed: some error",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s.EXPECT().RegisterUser(req).Times(1).Return(test.registerUser.err)

			RootCmd.SetArgs(test.args)

			var outBuf bytes.Buffer
			RootCmd.SetOutput(&outBuf)

			Execute(s)

			assert.Equal(t, test.output, outBuf.String())
		})
	}
}
