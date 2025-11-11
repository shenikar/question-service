package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// Config хранит все конфигурации приложения.
type Config struct {
	DatabaseURL string
}

// Load считывает конфигурацию из .env файла или переменных окружения.
func Load(log *logrus.Logger) (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Infoln("Info: .env file not found, loading from environment variables") // <-- Используем logrus
	}

	config := &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}

	return config, nil
}

// GetDatabaseURL возвращает строку подключения к базе данных.
func (c *Config) GetDatabaseURL() string {
	return c.DatabaseURL
}
