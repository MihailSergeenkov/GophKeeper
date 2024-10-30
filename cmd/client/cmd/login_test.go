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

func TestLoginCmd(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	s := mocks.NewMockServicer(mockCtrl)

	req := models.CreateUserTokenRequest{
		Login:    "qwe",
		Password: "123",
	}

	type loginUser struct {
		err error
	}
	type syncData struct {
		count int
		err   error
	}
	tests := []struct {
		name      string
		args      []string
		loginUser loginUser
		syncData  syncData
		output    string
	}{
		{
			name: "login success",
			args: []string{"login", "-l", "qwe", "-p", "123"},
			loginUser: loginUser{
				err: nil,
			},
			syncData: syncData{
				count: 1,
				err:   nil,
			},
			output: "Login OK\n",
		},
		{
			name: "login failed with failed call login service",
			args: []string{"login", "-l", "qwe", "-p", "123"},
			loginUser: loginUser{
				err: errors.New("some error"),
			},
			syncData: syncData{
				count: 0,
				err:   nil,
			},
			output: "Failed: some error",
		},
		{
			name: "login failed with failed sync data",
			args: []string{"login", "-l", "qwe", "-p", "123"},
			loginUser: loginUser{
				err: nil,
			},
			syncData: syncData{
				count: 1,
				err:   errors.New("some error"),
			},
			output: "Failed: some error",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s.EXPECT().LoginUser(req).Times(1).Return(test.loginUser.err)
			s.EXPECT().SyncData().Times(test.syncData.count).Return(test.syncData.err)

			RootCmd.SetArgs(test.args)

			var outBuf bytes.Buffer
			RootCmd.SetOutput(&outBuf)

			Execute(s)

			assert.Equal(t, test.output, outBuf.String())
		})
	}
}
