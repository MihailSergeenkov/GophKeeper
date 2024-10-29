package crypt

import (
	"testing"

	"github.com/MihailSergeenkov/GophKeeper/internal/server/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCrypt(t *testing.T) {
	t.Run("init crypt", func(t *testing.T) {
		settings, err := config.Setup(false)
		require.NoError(t, err)

		c, err := NewCrypt(settings)

		require.NoError(t, err)
		assert.Equal(t, settings, c.settings)
	})
}

func TestEncryptData(t *testing.T) {
	t.Run("encrypt data", func(t *testing.T) {
		settings, err := config.Setup(false)
		require.NoError(t, err)
		c, err := NewCrypt(settings)
		require.NoError(t, err)

		someData := []byte("some data")
		result := c.EncryptData(someData)

		assert.NotEmpty(t, result)
	})
}

func TestDecryptData(t *testing.T) {
	settings, err := config.Setup(false)
	require.NoError(t, err)
	c, err := NewCrypt(settings)
	require.NoError(t, err)

	someData := []byte("some data")

	tests := []struct {
		name     string
		someData []byte
		wantErr  bool
	}{
		{
			name:     "decrypt success",
			someData: c.EncryptData(someData),
			wantErr:  false,
		},
		{
			name:     "decrypt failed",
			someData: []byte("some encrypt data"),
			wantErr:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := c.DecryptData(test.someData)

			if test.wantErr {
				require.Error(t, err)
				assert.ErrorContains(t, err, "failed to decrypt data")
			} else {
				require.NoError(t, err)
				assert.Equal(t, result, someData)
			}
		})
	}
}
