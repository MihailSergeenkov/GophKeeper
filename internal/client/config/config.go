package config

import (
	"fmt"
	"maps"
	"os"
	"strconv"

	"github.com/MihailSergeenkov/GophKeeper/internal/models"
	"github.com/spf13/viper"
)

const (
	requestRetryDefault   = 3
	requestTimeoutDefault = 5
)

var cfg Config

type Config struct {
	Data           map[string]models.UserData
	ServerAPI      string `mapstructure:"server_api"`
	Token          string
	RequestRetry   int `mapstructure:"request_retry"`
	RequestTimeout int `mapstructure:"request_timeout"`
}

func Initializer(cfgFile *string) func() {
	return func() {
		initConfigFile(cfgFile)

		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok { //nolint:errorlint // Не работает с AS
				fmt.Println("Can't read config:", err)
				os.Exit(1)
			}

			viper.SetDefault("server_api", "http://localhost:8080/api")
			viper.SetDefault("request_retry", requestRetryDefault)
			viper.SetDefault("request_timeout", requestTimeoutDefault)
			viper.SetDefault("token", "")

			if err := viper.SafeWriteConfig(); err != nil {
				fmt.Println("Error while creating config file:", err)
				os.Exit(1)
			}
		}

		if err := viper.Unmarshal(&cfg); err != nil {
			fmt.Println("Error unmarshal config:", err)
			os.Exit(1)
		}
	}
}

func (cfg Config) GetServerAPI() string {
	return cfg.ServerAPI
}

func (cfg Config) GetRequestRetry() int {
	return cfg.RequestRetry
}

func (cfg Config) GetRequestTimeout() int {
	return cfg.RequestTimeout
}

func (cfg Config) GetToken() string {
	return cfg.Token
}

func (cfg Config) GetData() map[string]models.UserData {
	return cfg.Data
}

func (cfg *Config) UpdateToken(token string) error {
	viper.Set("token", token)

	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("failed update config file: %w", err)
	}

	cfg.Token = token

	return nil
}

func (cfg *Config) UpdateData(data []models.UserData) error {
	updateData := make(map[string]models.UserData, len(data))
	for _, v := range data {
		if v.Type == "file" {
			updateData[v.Mark] = v
			continue
		}

		updateData[strconv.Itoa(v.ID)] = v
	}
	viper.Set("data", updateData)

	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("failed update config file: %w", err)
	}

	cfg.Data = updateData

	return nil
}

func (cfg *Config) AddData(data models.UserData) error {
	updateData := make(map[string]models.UserData, 0)
	maps.Copy(updateData, cfg.Data)
	key := getKey(data)
	updateData[key] = data
	viper.Set("data", updateData)

	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("failed update config file: %w", err)
	}

	cfg.Data = updateData

	return nil
}

func GetConfig() *Config {
	return &cfg
}

func initConfigFile(cfgFile *string) {
	if *cfgFile != "" {
		viper.SetConfigFile(*cfgFile)
		return
	}

	viper.SetConfigName(".goph-keeper")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME")
}

func getKey(data models.UserData) string {
	if data.Type == "file" {
		return data.Mark
	}

	return strconv.Itoa(data.ID)
}
