package cmd

import (
	"bytes"
	"errors"
	"testing"

	"github.com/MihailSergeenkov/GophKeeper/cmd/client/cmd/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestLogoutCmd(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	s := mocks.NewMockServicer(mockCtrl)

	type logoutUser struct {
		err error
	}
	tests := []struct {
		name       string
		args       []string
		logoutUser logoutUser
		output     string
	}{
		{
			name: "logout success",
			args: []string{"logout"},
			logoutUser: logoutUser{
				err: nil,
			},
			output: "Logout OK\n",
		},
		{
			name: "logout failed",
			args: []string{"logout"},
			logoutUser: logoutUser{
				err: errors.New("some error"),
			},
			output: "Failed: some error",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s.EXPECT().LogoutUser().Times(1).Return(test.logoutUser.err)

			RootCmd.SetArgs(test.args)

			var outBuf bytes.Buffer
			RootCmd.SetOutput(&outBuf)

			Execute(s)

			assert.Equal(t, test.output, outBuf.String())
		})
	}
}
