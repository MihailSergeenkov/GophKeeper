package routes

import (
	"net/http"
	"testing"

	"github.com/MihailSergeenkov/GophKeeper/internal/server/config"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/routes/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNewRouter(t *testing.T) {
	t.Run("init router", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		handlres := mocks.NewMockHandlerer(mockCtrl)
		settings, err := config.Setup(false)
		require.NoError(t, err)

		logger := zap.NewNop()
		storage := mocks.NewMockStorager(mockCtrl)

		handlres.EXPECT().Ping().Times(1)
		handlres.EXPECT().RegisterUser().Times(1)
		handlres.EXPECT().CreateUserToken().Times(1)
		handlres.EXPECT().FetchUserData().Times(1)
		handlres.EXPECT().GetPassword().Times(1)
		handlres.EXPECT().AddPassword().Times(1)
		handlres.EXPECT().GetCard().Times(1)
		handlres.EXPECT().AddCard().Times(1)
		handlres.EXPECT().GetText().Times(1)
		handlres.EXPECT().AddText().Times(1)
		handlres.EXPECT().GetFile().Times(1)
		handlres.EXPECT().AddFile().Times(1)

		r := NewRouter(handlres, settings, logger, storage)
		assert.Implements(t, (*chi.Router)(nil), r)
	})
}

func closeBody(t *testing.T, r *http.Response) {
	t.Helper()
	err := r.Body.Close()

	if err != nil {
		t.Log(err)
	}
}
