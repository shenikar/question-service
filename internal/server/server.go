package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

// Server представляет HTTP-сервер.
type Server struct {
	httpServer *http.Server
	logger     *logrus.Logger
}

// NewServer создает новый экземпляр сервера.
func NewServer(handler http.Handler, logger *logrus.Logger) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:              ":8080",
			Handler:           handler,
			ReadHeaderTimeout: 5 * time.Second,
		},
		logger: logger,
	}
}

// Run запускает сервер и настраивает graceful shutdown.
func (s *Server) Run() error {
	// Канал для получения сигналов ОС
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	// Запускаем сервер в отдельной горутине
	go func() {
		s.logger.Infof("Starting server on %s", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Ожидаем сигнал завершения
	<-stop
	s.logger.Info("Shutting down server...")

	// Создаем контекст с таймаутом для graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Корректно завершаем работу сервера
	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Fatalf("Server shutdown failed: %v", err)
	}

	s.logger.Info("Server gracefully stopped.")
	return nil
}
