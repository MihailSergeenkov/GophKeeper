package requests

import (
	"testing"

	"github.com/MihailSergeenkov/GophKeeper/internal/client/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRequests(t *testing.T) {
	t.Run("init requests", func(t *testing.T) {
		cfg := config.GetConfig()

		result := NewRequests(cfg)

		assert.NotEmpty(t, result.r)
	})
}

func TestGet(t *testing.T) {
	t.Run("get request", func(t *testing.T) {
		cfg := config.GetConfig()
		r := NewRequests(cfg)

		_, err := r.Get("http://localhost/api")

		require.Error(t, err)
	})
}

func TestPost(t *testing.T) {
	t.Run("post request", func(t *testing.T) {
		cfg := config.GetConfig()
		r := NewRequests(cfg)

		_, err := r.Post("http://localhost/api")

		require.Error(t, err)
	})
}
