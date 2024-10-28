package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MihailSergeenkov/GophKeeper/internal/server/handlers/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func testGetRequest(t *testing.T, ts *httptest.Server, path string) (*http.Response, string) {
	t.Helper()

	req, err := http.NewRequest(http.MethodGet, ts.URL+path, http.NoBody)
	require.NoError(t, err)

	req.Header.Add(
		"X-Auth-Token",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOjF9.Glkq00DHorJbUkJxeD4TDfA9zqxFLWtfYiqCVfVUe9U",
	)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	closeBody(t, resp)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}
