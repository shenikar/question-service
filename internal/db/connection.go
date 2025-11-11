package db

import (
	"database/sql" // <-- Новый импорт
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/shenikar/question-service/internal/config"
)

// Connect устанавливает соединение с базой данных с помощью GORM
// Возвращает *gorm.DB и базовый *sql.DB для закрытия соединений.
func Connect(cfg *config.Config) (*gorm.DB, *sql.DB, error) {
	connStr := cfg.GetDatabaseURL()

	gormDB, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	log.Println("Connected to the database successfully")
	return gormDB, sqlDB, nil
}
