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

func TestGetCardCmd(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	s := mocks.NewMockServicer(mockCtrl)

	cardID := "1"
	data := models.Card{
		ID:          1,
		Number:      "1234123412341234",
		Owner:       "test",
		ExpiryDate:  "11/2300",
		CVV2:        "777",
		Mark:        "test",
		Description: "test",
	}
	expectedOutput := `{
  "number": "1234123412341234",
  "owner": "test",
  "expiry_date": "11/2300",
  "cvv2": "777",
  "mark": "test",
  "description": "test",
  "id": 1
}
`

	type getCard struct {
		resp models.Card
		err  error
	}
	tests := []struct {
		name    string
		args    []string
		getCard getCard
		output  string
	}{
		{
			name: "get card success",
			args: []string{"get", "card", cardID},
			getCard: getCard{
				resp: data,
				err:  nil,
			},
			output: expectedOutput,
		},
		{
			name: "get card failed",
			args: []string{"get", "card", cardID},
			getCard: getCard{
				resp: models.Card{},
				err:  errors.New("some error"),
			},
			output: "Failed: some error",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s.EXPECT().GetCard(cardID).Times(1).Return(test.getCard.resp, test.getCard.err)

			cmd.RootCmd.SetArgs(test.args)

			var outBuf bytes.Buffer
			cmd.RootCmd.SetOutput(&outBuf)

			cmd.Execute(s)

			assert.Equal(t, test.output, outBuf.String())
		})
	}
}
