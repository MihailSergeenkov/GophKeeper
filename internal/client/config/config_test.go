package config

import (
	"testing"

	"github.com/MihailSergeenkov/GophKeeper/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetConfig(t *testing.T) {
	t.Run("get config", func(t *testing.T) {
		cfgFile := ""
		init := Initializer(&cfgFile)
		init()

		cfg := GetConfig()

		assert.NotEmpty(t, cfg)
	})
}

func TestGetServerAPI(t *testing.T) {
	t.Run("get server api", func(t *testing.T) {
		cfgFile := ""
		init := Initializer(&cfgFile)
		init()

		cfg := GetConfig()
		result := cfg.GetServerAPI()

		assert.NotEmpty(t, result)
	})
}

func TestGetRequestRetry(t *testing.T) {
	t.Run("get request retry", func(t *testing.T) {
		cfgFile := ""
		init := Initializer(&cfgFile)
		init()

		cfg := GetConfig()
		result := cfg.GetRequestRetry()

		assert.NotEmpty(t, result)
	})
}

func TestGetRequestTimeout(t *testing.T) {
	t.Run("get request timeout", func(t *testing.T) {
		cfgFile := ""
		init := Initializer(&cfgFile)
		init()

		cfg := GetConfig()
		result := cfg.GetRequestTimeout()

		assert.NotEmpty(t, result)
	})
}

func TestUpdateToken(t *testing.T) {
	t.Run("update token", func(t *testing.T) {
		cfgFile := ""
		init := Initializer(&cfgFile)
		init()

		token := "test"

		cfg := GetConfig()
		err := cfg.UpdateToken(token)

		require.NoError(t, err)
	})
}

func TestUpdateData(t *testing.T) {
	t.Run("update data", func(t *testing.T) {
		cfgFile := ""
		init := Initializer(&cfgFile)
		init()

		data := []models.UserData{
			{
				ID:          1,
				Type:        "text",
				Mark:        "test",
				Description: "test",
			},
			{
				ID:          1,
				Type:        "file",
				Mark:        "filetest",
				Description: "test",
			},
		}

		cfg := GetConfig()
		err := cfg.UpdateData(data)

		require.NoError(t, err)
	})
}

func TestAddData(t *testing.T) {
	cfgFile := ""
	init := Initializer(&cfgFile)
	init()

	tests := []struct {
		name string
		data models.UserData
	}{
		{
			name: "add text",
			data: models.UserData{
				ID:          1,
				Type:        "text",
				Mark:        "test",
				Description: "test",
			},
		},
		{
			name: "add file",
			data: models.UserData{
				ID:          1,
				Type:        "file",
				Mark:        "filetest",
				Description: "test",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg := GetConfig()
			err := cfg.AddData(test.data)

			require.NoError(t, err)
		})
	}
}
