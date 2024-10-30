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
		}

		cfg := GetConfig()
		err := cfg.UpdateData(data)

		require.NoError(t, err)
	})
}

func TestAddData(t *testing.T) {
	t.Run("add data", func(t *testing.T) {
		cfgFile := ""
		init := Initializer(&cfgFile)
		init()

		data := models.UserData{
			ID:          1,
			Type:        "text",
			Mark:        "test",
			Description: "test",
		}

		cfg := GetConfig()
		err := cfg.AddData(data)

		require.NoError(t, err)
	})
}
