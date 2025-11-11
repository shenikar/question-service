package main

import (
	"github.com/shenikar/question-service/internal/config"
	"github.com/shenikar/question-service/internal/db"
	"github.com/shenikar/question-service/internal/handler"
	"github.com/shenikar/question-service/internal/logger"
	"github.com/shenikar/question-service/internal/repository"
	"github.com/shenikar/question-service/internal/router"
	"github.com/shenikar/question-service/internal/server"
	"github.com/shenikar/question-service/internal/service"
)

func main() {
	// Инициализация логгера
	appLogger := logger.NewLogger()

	// Загрузка конфигурации
	cfg, err := config.Load(appLogger)
	if err != nil {
		appLogger.Fatalf("Error loading .env file: %v", err)
	}

	// Подключение к базе данных
	gormDB, sqlDB, err := db.Connect(cfg, appLogger)
	if err != nil {
		appLogger.Fatalf("failed to connect database: %v", err)
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			appLogger.Errorf("Error closing database connection: %v", err)
		}
	}()

	// Инициализация репозитория
	repo := repository.NewRepository(gormDB, appLogger)

	// Инициализация сервисов
	s := service.NewService(repo, appLogger)

	// Инициализация обработчиков
	h := handler.NewHandler(s, appLogger)

	// Настройка роутера
	r := router.NewRouter(h)

	// Инициализация и запуск сервера
	srv := server.NewServer(r, appLogger)
	if err := srv.Run(); err != nil {
		appLogger.Fatalf("Server stopped with error: %v", err)
	}
}
