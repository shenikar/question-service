package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/shenikar/question-service/internal/handler"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(h *handler.Handler) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Swagger
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// Маршруты для вопросов
	r.Get("/questions", h.GetQuestions)
	r.Post("/questions", h.CreateQuestion)
	r.Get("/questions/{id}", h.GetQuestion)
	r.Delete("/questions/{id}", h.DeleteQuestion)

	// Маршруты для ответов
	r.Post("/questions/{id}/answers", h.CreateAnswer)
	r.Get("/answers/{id}", h.GetAnswer)
	r.Delete("/answers/{id}", h.DeleteAnswer)

	return r
}
