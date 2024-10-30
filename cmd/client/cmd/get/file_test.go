package get

import (
	"bytes"
	"errors"
	"testing"

	"github.com/MihailSergeenkov/GophKeeper/cmd/client/cmd"
	"github.com/MihailSergeenkov/GophKeeper/cmd/client/cmd/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetFileCmd(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	s := mocks.NewMockServicer(mockCtrl)

	fileID := "1"
	dir := "."

	type getFile struct {
		err error
	}
	tests := []struct {
		name    string
		args    []string
		getFile getFile
		output  string
	}{
		{
			name: "get file success",
			args: []string{"get", "file", fileID},
			getFile: getFile{
				err: nil,
			},
			output: "File load in .",
		},
		{
			name: "get file failed",
			args: []string{"get", "file", fileID},
			getFile: getFile{
				err: errors.New("some error"),
			},
			output: "Failed: some error",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s.EXPECT().GetFile(fileID, dir).Times(1).Return(test.getFile.err)

			cmd.RootCmd.SetArgs(test.args)

			var outBuf bytes.Buffer
			cmd.RootCmd.SetOutput(&outBuf)

			cmd.Execute(s)

			assert.Equal(t, test.output, outBuf.String())
		})
	}
}
