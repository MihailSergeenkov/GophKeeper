package handlers

import (
	"net/http"
	"testing"

	"github.com/MihailSergeenkov/GophKeeper/internal/server/handlers/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewHandlers(t *testing.T) {
	t.Run("init handlers", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		s := mocks.NewMockServicer(mockCtrl)
		l := mocks.NewMockLogger(mockCtrl)

		handlers := NewHandlers(s, l)

		assert.Equal(t, s, handlers.services)
		assert.Equal(t, l, handlers.logger)
	})
}

func closeBody(t *testing.T, r *http.Response) {
	t.Helper()
	err := r.Body.Close()

	if err != nil {
		t.Log(err)
	}
}
