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

func TestGetTextCmd(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	s := mocks.NewMockServicer(mockCtrl)

	textID := "1"
	data := models.Text{
		ID:          1,
		Data:        "test",
		Mark:        "test",
		Description: "test",
	}
	expectedOutput := `{
  "data": "test",
  "mark": "test",
  "description": "test",
  "id": 1
}
`

	type getText struct {
		resp models.Text
		err  error
	}
	tests := []struct {
		name    string
		args    []string
		getText getText
		output  string
	}{
		{
			name: "get text success",
			args: []string{"get", "text", textID},
			getText: getText{
				resp: data,
				err:  nil,
			},
			output: expectedOutput,
		},
		{
			name: "get text failed",
			args: []string{"get", "text", textID},
			getText: getText{
				resp: models.Text{},
				err:  errors.New("some error"),
			},
			output: "Failed: some error",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s.EXPECT().GetText(textID).Times(1).Return(test.getText.resp, test.getText.err)

			cmd.RootCmd.SetArgs(test.args)

			var outBuf bytes.Buffer
			cmd.RootCmd.SetOutput(&outBuf)

			cmd.Execute(s)

			assert.Equal(t, test.output, outBuf.String())
		})
	}
}
