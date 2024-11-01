package add

import (
	"bytes"
	"errors"
	"testing"

	"github.com/MihailSergeenkov/GophKeeper/cmd/client/cmd"
	"github.com/MihailSergeenkov/GophKeeper/cmd/client/cmd/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAddFileCmd(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	s := mocks.NewMockServicer(mockCtrl)

	file := "./test.txt"
	mark := "test"
	description := "test"

	type addFile struct {
		err error
	}
	tests := []struct {
		name    string
		args    []string
		addFile addFile
		output  string
	}{
		{
			name: "add file success",
			args: []string{"add", "file", "-f", file, "-m", mark, "-d", description},
			addFile: addFile{
				err: nil,
			},
			output: "Add file OK\n",
		},
		{
			name: "add file failed",
			args: []string{"add", "file", "-f", file, "-m", mark, "-d", description},
			addFile: addFile{
				err: errors.New("some error"),
			},
			output: "Failed: some error",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s.EXPECT().AddFile(file, mark, description).Times(1).Return(test.addFile.err)

			cmd.RootCmd.SetArgs(test.args)

			var outBuf bytes.Buffer
			cmd.RootCmd.SetOutput(&outBuf)

			cmd.Execute(s)

			assert.Equal(t, test.output, outBuf.String())
		})
	}
}
