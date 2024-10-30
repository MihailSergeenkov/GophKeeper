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

func TestAddTextCmd(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	s := mocks.NewMockServicer(mockCtrl)

	req := models.AddTextRequest{
		Data:        "test",
		Mark:        "test",
		Description: "test",
	}

	type addText struct {
		err error
	}
	tests := []struct {
		name    string
		args    []string
		addText addText
		output  string
	}{
		{
			name: "add text success",
			args: []string{"add", "text", "-t", "test", "-m", "test", "-d", "test"},
			addText: addText{
				err: nil,
			},
			output: "Add text OK\n",
		},
		{
			name: "add text failed",
			args: []string{"add", "text", "-t", "test", "-m", "test", "-d", "test"},
			addText: addText{
				err: errors.New("some error"),
			},
			output: "Failed: some error",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s.EXPECT().AddText(req).Times(1).Return(test.addText.err)

			cmd.RootCmd.SetArgs(test.args)

			var outBuf bytes.Buffer
			cmd.RootCmd.SetOutput(&outBuf)

			cmd.Execute(s)

			assert.Equal(t, test.output, outBuf.String())
		})
	}
}
