package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestSetup(t *testing.T) {
	runAddr := "localhost:8081"

	tests := []struct {
		name    string
		setEnv  func()
		wantErr bool
		errText string
	}{
		{
			name: "success setup",
			setEnv: func() {
				require.NoError(t, os.Setenv("SERVER_ADDRESS", runAddr))
			},
			wantErr: false,
			errText: "",
		},
		{
			name: "failed to get config",
			setEnv: func() {
				require.NoError(t, os.Setenv("CONFIG", "set.json"))
			},
			wantErr: true,
			errText: "failed to get config",
		},
		{
			name: "failed to parse config",
			setEnv: func() {
				require.NoError(t, os.Setenv("CONFIG", "testdata/bad_settings.json"))
			},
			wantErr: true,
			errText: "failed to parse config data",
		},
		{
			name: "failed to parse envs",
			setEnv: func() {
				require.NoError(t, os.Setenv("SERVER_ADDRESS", "localhost:8081"))
				require.NoError(t, os.Setenv("LOG_LEVEL", "some string"))
			},
			wantErr: true,
			errText: "failed to parse envs",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			os.Clearenv()
			test.setEnv()

			config, err := Setup(false)

			if test.wantErr {
				require.Error(t, err)
				require.ErrorContains(t, err, test.errText)
			} else {
				require.NoError(t, err)
				assert.Equal(t, runAddr, config.RunAddr)
			}
		})
	}
}

func TestGetConfigData(t *testing.T) {
	tests := []struct {
		setEnv      func()
		name        string
		wantErr     bool
		presentData bool
	}{
		{
			name: "config present",
			setEnv: func() {
				require.NoError(t, os.Setenv("CONFIG", "testdata/settings.json"))
			},
			wantErr:     false,
			presentData: true,
		},
		{
			name: "config absent",
			setEnv: func() {
				require.NoError(t, os.Setenv("CONFIG", ""))
			},
			wantErr:     false,
			presentData: false,
		},
		{
			name: "config failed",
			setEnv: func() {
				require.NoError(t, os.Setenv("CONFIG", "set.json"))
			},
			wantErr:     true,
			presentData: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			os.Clearenv()
			test.setEnv()

			_, presentData, err := getConfigData()

			if test.wantErr {
				require.Error(t, err)
				require.ErrorContains(t, err, "failed to read config file")
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.presentData, presentData)
			}
		})
	}
}

func TestParseConfigData(t *testing.T) {
	config := Settings{}

	tests := []struct {
		name       string
		configFile string
		wantErr    bool
	}{
		{
			name:       "success parsed config",
			wantErr:    false,
			configFile: "testdata/settings.json",
		},
		{
			name:       "bad config",
			wantErr:    true,
			configFile: "testdata/bad_settings.json",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data, err := os.ReadFile(test.configFile)
			require.NoError(t, err)

			err = config.parseConfigData(data)

			if test.wantErr {
				require.Error(t, err)
				require.ErrorContains(t, err, "failed to unmarshal json")
			} else {
				require.NoError(t, err)
				assert.Equal(t, "localhost:8081", config.RunAddr)
				assert.Equal(t, "postgresql://localhost:5432/goph_keeper", config.DatabaseURI)
				assert.Equal(t, "12345", config.SecretKey)
				assert.Equal(t, zapcore.InfoLevel, config.LogLevel)
				assert.Equal(t, true, config.EnableHTTPS)
			}
		})
	}
}

func TestParseEnv(t *testing.T) {
	config := Settings{}

	tests := []struct {
		setEnv  func()
		name    string
		wantErr bool
	}{
		{
			name: "valid config",
			setEnv: func() {
				require.NoError(t, os.Setenv("SERVER_ADDRESS", "localhost:8081"))
				require.NoError(t, os.Setenv("LOG_LEVEL", "INFO"))
			},
			wantErr: false,
		},
		{
			name: "invalid config",
			setEnv: func() {
				require.NoError(t, os.Setenv("SERVER_ADDRESS", "localhost:8081"))
				require.NoError(t, os.Setenv("LOG_LEVEL", "some string"))
			},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			os.Clearenv()
			test.setEnv()

			err := config.parseEnv()

			if test.wantErr {
				require.Error(t, err)
				require.ErrorContains(t, err, "env error")
			} else {
				require.NoError(t, err)
			}
		})
	}
}
