package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"

	"github.com/go-chi/chi/v5"
	"github.com/shenikar/question-service/internal/models"
	"github.com/shenikar/question-service/internal/service" // <-- Новый импорт
)

// Handler обрабатывает HTTP-запросы.
type Handler struct {
	service service.Service // <-- Зависим от слоя сервисов
}

// NewHandler создает новый экземпляр обработчика.
func NewHandler(s service.Service) *Handler {
	return &Handler{service: s}
}

// CreateQuestion создает новый вопрос.
func (h *Handler) CreateQuestion(w http.ResponseWriter, r *http.Request) {
	var question models.Question
	if err := json.NewDecoder(r.Body).Decode(&question); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(&question); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.CreateQuestion(&question); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(question)
}

// GetQuestion получает вопрос по ID.
func (h *Handler) GetQuestion(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid question ID", http.StatusBadRequest)
		return
	}

	question, err := h.service.GetQuestion(uint(id))
	if err != nil {
		http.Error(w, "Question not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(question)
}

// GetQuestions получает все вопросы.
func (h *Handler) GetQuestions(w http.ResponseWriter, r *http.Request) {
	questions, err := h.service.GetAllQuestions()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(questions)
}

// DeleteQuestion удаляет вопрос по ID.
func (h *Handler) DeleteQuestion(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid question ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteQuestion(uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// CreateAnswer создает ответ на вопрос.
func (h *Handler) CreateAnswer(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid question ID", http.StatusBadRequest)
		return
	}

	var answer models.Answer
	if err := json.NewDecoder(r.Body).Decode(&answer); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(&answer); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.CreateAnswer(uint(id), &answer); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(answer)
}

// GetAnswer получает ответ по ID.
func (h *Handler) GetAnswer(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid answer ID", http.StatusBadRequest)
		return
	}

	answer, err := h.service.GetAnswer(uint(id))
	if err != nil {
		http.Error(w, "Answer not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(answer)
}

// DeleteAnswer удаляет ответ по ID.
func (h *Handler) DeleteAnswer(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid answer ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteAnswer(uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
