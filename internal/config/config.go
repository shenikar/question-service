package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config хранит все конфигурации приложения.
type Config struct {
	DatabaseURL string
}

// Load считывает конфигурацию из .env файла или переменных окружения.
func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("Info: .env file not found, loading from environment variables")
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
