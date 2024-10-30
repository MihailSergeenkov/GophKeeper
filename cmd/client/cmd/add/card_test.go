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

func TestAddCardCmd(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	s := mocks.NewMockServicer(mockCtrl)

	req := &models.AddCardRequest{
		Number:      "1234123412341234",
		Owner:       "test",
		ExpiryDate:  "11/2300",
		CVV2:        "777",
		Mark:        "test",
		Description: "test",
	}

	type addCard struct {
		err error
	}
	tests := []struct {
		name    string
		args    []string
		addCard addCard
		output  string
	}{
		{
			name: "add card success",
			args: []string{
				"add", "card", "-n", "1234123412341234",
				"-o", "test", "-e", "11/2300", "--cvv2", "777", "-m", "test", "-d", "test",
			},
			addCard: addCard{
				err: nil,
			},
			output: "Add card OK\n",
		},
		{
			name: "add card failed",
			args: []string{
				"add", "card", "-n", "1234123412341234",
				"-o", "test", "-e", "11/2300", "--cvv2", "777", "-m", "test", "-d", "test",
			},
			addCard: addCard{
				err: errors.New("some error"),
			},
			output: "Failed: some error",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s.EXPECT().AddCard(req).Times(1).Return(test.addCard.err)

			cmd.RootCmd.SetArgs(test.args)

			var outBuf bytes.Buffer
			cmd.RootCmd.SetOutput(&outBuf)

			cmd.Execute(s)

			assert.Equal(t, test.output, outBuf.String())
		})
	}
}
