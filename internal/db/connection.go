package db

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/shenikar/question-service/internal/config"
)

// Connect устанавливает соединение с базой данных с помощью GORM
func Connect(cfg *config.Config) (*gorm.DB, error) {
	connStr := cfg.GetDatabaseURL()

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Connected to the database successfully")
	return db, nil
}
