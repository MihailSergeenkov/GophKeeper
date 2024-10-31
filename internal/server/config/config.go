package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"dario.cat/mergo"
	"github.com/caarlos0/env/v11"
	"go.uber.org/zap/zapcore"
)

// Settings структура для конфигурирования сервиса.
type Settings struct {
	RunAddr     string        `json:"server_address" env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	DatabaseURI string        `json:"db_uri" env:"DATABASE_URI" envDefault:"postgresql://localhost:5432/test"`
	SecretKey   string        `json:"secret_key" env:"SECRET_KEY" envDefault:"1234567890"`
	S3          S3Settings    `json:"s3"`
	LogLevel    zapcore.Level `json:"log_level" env:"LOG_LEVEL" envDefault:"ERROR"`
	EnableHTTPS bool          `json:"enable_https" env:"ENABLE_HTTPS" envDefault:"false"`
}

type S3Settings struct {
	Endpoint        string `json:"endpoint" env:"S3_ENDPOINT" envDefault:"localhost:9000"`
	AccessKeyID     string `json:"access_key_id" env:"S3_ACCESS_KEY_ID" envDefault:"test_id"`
	SecretAccessKey string `json:"secret_access_key" env:"S3_SECRET_ACCESS_KEY" envDefault:"test_secret"`
	Region          string `json:"region" env:"S3_REGION" envDefault:"us-east-1"`
	SecretPassword  string `json:"secret_password" env:"S3_SECRET_PASSWORD" envDefault:"12345678"`
	UseSSL          bool   `json:"use_ssl" env:"S3_USE_SSL" envDefault:"false"`
	SecureFiles     bool   `json:"secure_files" env:"S3_SECURE_FILES" envDefault:"false"`
}

// Setup функция считывания и применения пользовательских настроек сервиса.
func Setup(withFlags bool) (*Settings, error) {
	s := Settings{LogLevel: zapcore.ErrorLevel}

	configData, presentData, err := getConfigData()
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}
	if presentData {
		if err := s.parseConfigData(configData); err != nil {
			return nil, fmt.Errorf("failed to parse config data: %w", err)
		}
	}

	if err := s.parseEnv(); err != nil {
		return nil, fmt.Errorf("failed to parse envs: %w", err)
	}

	if withFlags {
		s.parseFlags()
	}

	return &s, nil
}

func getConfigData() ([]byte, bool, error) {
	configFile := os.Getenv("CONFIG")

	for i, arg := range os.Args {
		if arg == "-c" || arg == "-config" {
			configFile = os.Args[i+1]
			break
		}
	}

	if configFile == "" {
		return []byte{}, false, nil
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return []byte{}, false, fmt.Errorf("failed to read config file: %w", err)
	}

	return data, true, nil
}

func (s *Settings) parseConfigData(data []byte) error {
	config := Settings{}

	err := json.Unmarshal(data, &config)
	if err != nil {
		return fmt.Errorf("failed to unmarshal json: %w", err)
	}

	if err := mergo.Merge(s, config); err != nil {
		return fmt.Errorf("failed to map config to settings: %w", err)
	}

	return nil
}

func (s *Settings) parseEnv() error {
	if err := env.Parse(s); err != nil {
		return fmt.Errorf("env error: %w", err)
	}

	return nil
}

func (s *Settings) parseFlags() {
	flag.StringVar(&s.RunAddr, "a", s.RunAddr, "address and port to run server")
	flag.Func("l", `level for logger (default "ERROR")`, func(v string) error {
		lev, err := zapcore.ParseLevel(v)

		if err != nil {
			return fmt.Errorf("parse log level env error: %w", err)
		}

		s.LogLevel = lev
		return nil
	})

	flag.StringVar(&s.DatabaseURI, "d", s.DatabaseURI, "database URI")
	flag.StringVar(&s.SecretKey, "sk", s.SecretKey, "secret key for generate cookie token")
	flag.BoolVar(&s.EnableHTTPS, "s", s.EnableHTTPS, "enable HTTPS")

	flag.StringVar(&s.S3.Endpoint, "se", s.S3.Endpoint, "address and port for s3")
	flag.StringVar(&s.S3.AccessKeyID, "sa", s.S3.AccessKeyID, "access key id for s3")
	flag.StringVar(&s.S3.SecretAccessKey, "ss", s.S3.SecretAccessKey, "secret access key for s3")
	flag.BoolVar(&s.S3.UseSSL, "su", s.S3.UseSSL, "enable SSL for S3")
	flag.StringVar(&s.S3.Region, "sr", s.S3.Region, "region for s3")
	flag.StringVar(&s.S3.SecretPassword, "sp", s.S3.SecretPassword, "secret password for s3")
	flag.BoolVar(&s.S3.SecureFiles, "sf", s.S3.SecureFiles, "secure files in S3")

	flag.String("c", "", "config file path (shorthand)")
	flag.String("config", "", "config file path")

	flag.Parse()
}
