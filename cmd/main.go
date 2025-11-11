package main

import (
	"log"
	"github.com/shenikar/question-service/internal/config"
	"github.com/shenikar/question-service/internal/db"
	"github.com/shenikar/question-service/internal/handler"
	"github.com/shenikar/question-service/internal/repository"
	"github.com/shenikar/question-service/internal/router"
	"github.com/shenikar/question-service/internal/server"
	"github.com/shenikar/question-service/internal/service"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Подключение к базе данных
	gormDB, sqlDB, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer func() {
		if err := sqlDB.Close(); err != nil { // Корректно закрываем sqlDB
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	// Инициализация репозитория
	repo := repository.NewRepository(gormDB)

	// Инициализация сервисов
	s := service.NewService(repo)

	// Инициализация обработчиков
	h := handler.NewHandler(s)

	// Настройка роутера
	r := router.NewRouter(h)

	// Инициализация и запуск сервера
	srv := server.NewServer(r)
	if err := srv.Run(); err != nil {
		log.Fatalf("Server stopped with error: %v", err)
	}
}