package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"

	"github.com/shenikar/question-service/internal/models"
	"github.com/shenikar/question-service/internal/service"
)

// Handler обрабатывает HTTP-запросы.
type Handler struct {
	service service.Service
	logger  *logrus.Logger
}

// NewHandler создает новый экземпляр обработчика.
func NewHandler(s service.Service, logger *logrus.Logger) *Handler {
	return &Handler{service: s, logger: logger}
}

// CreateQuestion создает новый вопрос.
func (h *Handler) CreateQuestion(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Received request to create question")
	var question models.Question
	if err := json.NewDecoder(r.Body).Decode(&question); err != nil {
		h.logger.Warnf("Failed to decode request body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(&question); err != nil {
		h.logger.Warnf("Validation failed for question: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.CreateQuestion(&question); err != nil {
		h.logger.Errorf("Failed to create question: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(question); err != nil {
		h.logger.Errorf("Failed to encode response for CreateQuestion: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	h.logger.Infof("Question created successfully with ID: %d", question.ID)
}

// GetQuestion получает вопрос по ID.
func (h *Handler) GetQuestion(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	h.logger.Infof("Received request to get question with ID: %s", idStr)
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.logger.Warnf("Invalid question ID: %s, error: %v", idStr, err)
		http.Error(w, "Invalid question ID", http.StatusBadRequest)
		return
	}

	question, err := h.service.GetQuestion(uint(id))
	if err != nil {
		h.logger.Errorf("Failed to get question with ID %d: %v", id, err)
		http.Error(w, "Question not found", http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(question); err != nil {
		h.logger.Errorf("Failed to encode response for GetQuestion: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	h.logger.Infof("Question with ID %d retrieved successfully", id)
}

// GetQuestions получает все вопросы.
func (h *Handler) GetQuestions(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Received request to get all questions")
	questions, err := h.service.GetAllQuestions()
	if err != nil {
		h.logger.Errorf("Failed to get all questions: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(questions); err != nil {
		h.logger.Errorf("Failed to encode response for GetQuestions: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	h.logger.Info("All questions retrieved successfully")
}

// DeleteQuestion удаляет вопрос по ID.
func (h *Handler) DeleteQuestion(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	h.logger.Infof("Received request to delete question with ID: %s", idStr)
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.logger.Warnf("Invalid question ID for deletion: %s, error: %v", idStr, err)
		http.Error(w, "Invalid question ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteQuestion(uint(id)); err != nil {
		h.logger.Errorf("Failed to delete question with ID %d: %v", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	h.logger.Infof("Question with ID %d deleted successfully", id)
}

// CreateAnswer создает ответ на вопрос.
func (h *Handler) CreateAnswer(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	h.logger.Infof("Received request to create answer for question ID: %s", idStr)
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.logger.Warnf("Invalid question ID for answer creation: %s, error: %v", idStr, err)
		http.Error(w, "Invalid question ID", http.StatusBadRequest)
		return
	}

	var answer models.Answer
	if err := json.NewDecoder(r.Body).Decode(&answer); err != nil {
		h.logger.Warnf("Failed to decode answer request body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(&answer); err != nil {
		h.logger.Warnf("Validation failed for answer: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.CreateAnswer(uint(id), &answer); err != nil {
		h.logger.Errorf("Failed to create answer for question ID %d: %v", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(answer); err != nil {
		h.logger.Errorf("Failed to encode response for CreateAnswer: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	h.logger.Infof("Answer created successfully for question ID %d", id)
}

// GetAnswer получает ответ по ID.
func (h *Handler) GetAnswer(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	h.logger.Infof("Received request to get answer with ID: %s", idStr)
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.logger.Warnf("Invalid answer ID: %s, error: %v", idStr, err)
		http.Error(w, "Invalid answer ID", http.StatusBadRequest)
		return
	}

	answer, err := h.service.GetAnswer(uint(id))
	if err != nil {
		h.logger.Errorf("Failed to get answer with ID %d: %v", id, err)
		http.Error(w, "Answer not found", http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(answer); err != nil {
		h.logger.Errorf("Failed to encode response for GetAnswer: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	h.logger.Infof("Answer with ID %d retrieved successfully", id)
}

// DeleteAnswer удаляет ответ по ID.
func (h *Handler) DeleteAnswer(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	h.logger.Infof("Received request to delete answer with ID: %s", idStr)
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.logger.Warnf("Invalid answer ID for deletion: %s, error: %v", idStr, err)
		http.Error(w, "Invalid answer ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteAnswer(uint(id)); err != nil {
		h.logger.Errorf("Failed to delete answer with ID %d: %v", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	h.logger.Infof("Answer with ID %d deleted successfully", id)
}
