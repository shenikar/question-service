package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/shenikar/question-service/internal/models"
)

// MockService - мок для интерфейса service.Service
type MockService struct {
	mock.Mock
}

func (m *MockService) CreateQuestion(question *models.Question) error {
	args := m.Called(question)
	return args.Error(0)
}

func (m *MockService) GetQuestion(id uint) (*models.Question, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Question), args.Error(1)
}

func (m *MockService) GetAllQuestions() ([]models.Question, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Question), args.Error(1)
}

func (m *MockService) DeleteQuestion(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockService) CreateAnswer(questionID uint, answer *models.Answer) error {
	args := m.Called(questionID, answer)
	return args.Error(0)
}

func (m *MockService) GetAnswer(id uint) (*models.Answer, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Answer), args.Error(1)
}

func (m *MockService) DeleteAnswer(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreateQuestionHandler(t *testing.T) {
	mockService := new(MockService)
	logger := logrus.New()
	handler := NewHandler(mockService, logger)

	question := &models.Question{Text: "Test Question"}
	questionJSON, _ := json.Marshal(question)

	req := httptest.NewRequest(http.MethodPost, "/questions", bytes.NewBuffer(questionJSON))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	mockService.On("CreateQuestion", mock.AnythingOfType("*models.Question")).Return(nil)

	handler.CreateQuestion(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	mockService.AssertExpectations(t)
}

func TestCreateQuestionHandlerServiceError(t *testing.T) {
	mockService := new(MockService)
	logger := logrus.New()
	handler := NewHandler(mockService, logger)

	question := &models.Question{Text: "Test Question"}
	questionJSON, _ := json.Marshal(question)

	req := httptest.NewRequest(http.MethodPost, "/questions", bytes.NewBuffer(questionJSON))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	mockService.On("CreateQuestion", mock.AnythingOfType("*models.Question")).Return(errors.New("service error"))

	handler.CreateQuestion(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	mockService.AssertExpectations(t)
}

func TestCreateQuestionHandlerInvalidInput(t *testing.T) {
	mockService := new(MockService)
	logger := logrus.New()
	handler := NewHandler(mockService, logger)

	questionJSON := []byte(`{"text": ""}`) // Пустой текст
	req := httptest.NewRequest(http.MethodPost, "/questions", bytes.NewBuffer(questionJSON))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.CreateQuestion(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockService.AssertNotCalled(t, "CreateQuestion", mock.Anything)
}

func TestGetQuestionHandler(t *testing.T) {
	mockService := new(MockService)
	logger := logrus.New()
	handler := NewHandler(mockService, logger)

	expectedQuestion := &models.Question{ID: 1, Text: "Test Question"}

	mockService.On("GetQuestion", uint(1)).Return(expectedQuestion, nil)

	req := httptest.NewRequest(http.MethodGet, "/questions/1", nil)
	rr := httptest.NewRecorder()

	// Используем chi.NewRouter для обработки URL-параметров
	r := chi.NewRouter()
	r.Get("/questions/{id}", handler.GetQuestion)
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var responseQuestion models.Question
	assert.NoError(t, json.NewDecoder(rr.Body).Decode(&responseQuestion))
	assert.Equal(t, expectedQuestion.ID, responseQuestion.ID)
	mockService.AssertExpectations(t)
}

func TestGetAllQuestionsHandler(t *testing.T) {
	mockService := new(MockService)
	logger := logrus.New()
	handler := NewHandler(mockService, logger)

	expectedQuestions := []models.Question{
		{ID: 1, Text: "Question 1"},
		{ID: 2, Text: "Question 2"},
	}

	mockService.On("GetAllQuestions").Return(expectedQuestions, nil)

	req := httptest.NewRequest(http.MethodGet, "/questions", nil)
	rr := httptest.NewRecorder()

	handler.GetQuestions(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var responseQuestions []models.Question
	assert.NoError(t, json.NewDecoder(rr.Body).Decode(&responseQuestions))
	assert.Len(t, responseQuestions, 2)
	assert.Equal(t, expectedQuestions[0].Text, responseQuestions[0].Text)
	mockService.AssertExpectations(t)
}

func TestGetQuestionHandlerInvalidID(t *testing.T) {
	mockService := new(MockService)
	logger := logrus.New()
	handler := NewHandler(mockService, logger)

	req := httptest.NewRequest(http.MethodGet, "/questions/abc", nil) // Некорректный ID
	rr := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/questions/{id}", handler.GetQuestion)
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockService.AssertNotCalled(t, "GetQuestion", mock.Anything)
}

func TestGetAllQuestionsHandlerError(t *testing.T) {
	mockService := new(MockService)
	logger := logrus.New()
	handler := NewHandler(mockService, logger)

	mockService.On("GetAllQuestions").Return(nil, errors.New("database error"))

	req := httptest.NewRequest(http.MethodGet, "/questions", nil)
	rr := httptest.NewRecorder()

	handler.GetQuestions(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	mockService.AssertExpectations(t)
}

func TestDeleteQuestionHandler(t *testing.T) {
	mockService := new(MockService)
	logger := logrus.New()
	handler := NewHandler(mockService, logger)

	mockService.On("DeleteQuestion", uint(1)).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/questions/1", nil)
	rr := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Delete("/questions/{id}", handler.DeleteQuestion)
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
	mockService.AssertExpectations(t)
}

func TestGetAllQuestionsHandlerEmpty(t *testing.T) {
	mockService := new(MockService)
	logger := logrus.New()
	handler := NewHandler(mockService, logger)

	mockService.On("GetAllQuestions").Return([]models.Question{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/questions", nil)
	rr := httptest.NewRecorder()

	handler.GetQuestions(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var responseQuestions []models.Question
	assert.NoError(t, json.NewDecoder(rr.Body).Decode(&responseQuestions))
	assert.Len(t, responseQuestions, 0)
	mockService.AssertExpectations(t)
}

func TestDeleteQuestionHandlerNotFound(t *testing.T) {
	mockService := new(MockService)
	logger := logrus.New()
	handler := NewHandler(mockService, logger)

	mockService.On("DeleteQuestion", uint(999)).Return(errors.New("not found"))

	req := httptest.NewRequest(http.MethodDelete, "/questions/999", nil)
	rr := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Delete("/questions/{id}", handler.DeleteQuestion)
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code) // Ошибка сервиса, возвращаем 500
	mockService.AssertExpectations(t)
}

func TestCreateAnswerHandler(t *testing.T) {
	mockService := new(MockService)
	logger := logrus.New()
	handler := NewHandler(mockService, logger)

	questionID := uint(1)
	answer := &models.Answer{Text: "Test Answer"}
	answerJSON, _ := json.Marshal(answer)

	req := httptest.NewRequest(http.MethodPost, "/questions/1/answers", bytes.NewBuffer(answerJSON))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	mockService.On("CreateAnswer", questionID, mock.AnythingOfType("*models.Answer")).Return(nil)

	r := chi.NewRouter()
	r.Post("/questions/{id}/answers", handler.CreateAnswer)
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	mockService.AssertExpectations(t)
}

func TestDeleteQuestionHandlerInvalidID(t *testing.T) {
	mockService := new(MockService)
	logger := logrus.New()
	handler := NewHandler(mockService, logger)

	req := httptest.NewRequest(http.MethodDelete, "/questions/abc", nil) // Некорректный ID
	rr := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Delete("/questions/{id}", handler.DeleteQuestion)
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockService.AssertNotCalled(t, "DeleteQuestion", mock.Anything)
}

func TestCreateAnswerHandlerInvalidInput(t *testing.T) {
	mockService := new(MockService)
	logger := logrus.New()
	handler := NewHandler(mockService, logger)

	answerJSON := []byte(`{"text": ""}`) // Пустой текст
	req := httptest.NewRequest(http.MethodPost, "/questions/1/answers", bytes.NewBuffer(answerJSON))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Post("/questions/{id}/answers", handler.CreateAnswer)
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockService.AssertNotCalled(t, "CreateAnswer", mock.Anything, mock.Anything)
}

func TestCreateAnswerHandlerServiceError(t *testing.T) {
	mockService := new(MockService)
	logger := logrus.New()
	handler := NewHandler(mockService, logger)

	questionID := uint(1)
	answer := &models.Answer{Text: "Test Answer"}
	answerJSON, _ := json.Marshal(answer)

	mockService.On("CreateAnswer", questionID, mock.AnythingOfType("*models.Answer")).
		Return(errors.New("service error"))

	req := httptest.NewRequest(http.MethodPost, "/questions/1/answers", bytes.NewBuffer(answerJSON))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Post("/questions/{id}/answers", handler.CreateAnswer)
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	mockService.AssertExpectations(t)
}

func TestGetAnswerHandler(t *testing.T) {
	mockService := new(MockService)
	logger := logrus.New()
	handler := NewHandler(mockService, logger)

	expectedAnswer := &models.Answer{ID: 1, QuestionID: 1, Text: "Test Answer"}

	mockService.On("GetAnswer", uint(1)).Return(expectedAnswer, nil)

	req := httptest.NewRequest(http.MethodGet, "/answers/1", nil)
	rr := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/answers/{id}", handler.GetAnswer)
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var responseAnswer models.Answer
	assert.NoError(t, json.NewDecoder(rr.Body).Decode(&responseAnswer))
	assert.Equal(t, expectedAnswer.ID, responseAnswer.ID)
	mockService.AssertExpectations(t)
}

func TestCreateAnswerHandlerInvalidQuestionID(t *testing.T) {
	mockService := new(MockService)
	logger := logrus.New()
	handler := NewHandler(mockService, logger)

	answer := &models.Answer{Text: "Test Answer"}
	answerJSON, _ := json.Marshal(answer)

	req := httptest.NewRequest(http.MethodPost, "/questions/abc/answers", // Некорректный ID вопроса
		bytes.NewBuffer(answerJSON))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Post("/questions/{id}/answers", handler.CreateAnswer)
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockService.AssertNotCalled(t, "CreateAnswer", mock.Anything, mock.Anything)
}

func TestGetAnswerHandlerNotFound(t *testing.T) {
	mockService := new(MockService)
	logger := logrus.New()
	handler := NewHandler(mockService, logger)

	mockService.On("GetAnswer", uint(999)).Return(nil, errors.New("not found"))

	req := httptest.NewRequest(http.MethodGet, "/answers/999", nil)
	rr := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/answers/{id}", handler.GetAnswer)
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	mockService.AssertExpectations(t)
}

func TestGetQuestionHandlerNotFound(t *testing.T) {
	mockService := new(MockService)
	logger := logrus.New()
	handler := NewHandler(mockService, logger)

	mockService.On("GetQuestion", uint(999)).Return(nil, errors.New("not found"))

	req := httptest.NewRequest(http.MethodGet, "/questions/999", nil)
	rr := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/questions/{id}", handler.GetQuestion)
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	mockService.AssertExpectations(t)
}

func TestDeleteAnswerHandler(t *testing.T) {
	mockService := new(MockService)
	logger := logrus.New()
	handler := NewHandler(mockService, logger)

	mockService.On("DeleteAnswer", uint(1)).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/answers/1", nil)
	rr := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Delete("/answers/{id}", handler.DeleteAnswer)
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
	mockService.AssertExpectations(t)
}

func TestGetAnswerHandlerInvalidID(t *testing.T) {
	mockService := new(MockService)
	logger := logrus.New()
	handler := NewHandler(mockService, logger)

	req := httptest.NewRequest(http.MethodGet, "/answers/abc", nil) // Некорректный ID
	rr := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/answers/{id}", handler.GetAnswer)
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockService.AssertNotCalled(t, "GetAnswer", mock.Anything)
}

func TestDeleteAnswerHandlerNotFound(t *testing.T) {
	mockService := new(MockService)
	logger := logrus.New()
	handler := NewHandler(mockService, logger)

	mockService.On("DeleteAnswer", uint(999)).Return(errors.New("not found"))

	req := httptest.NewRequest(http.MethodDelete, "/answers/999", nil)
	rr := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Delete("/answers/{id}", handler.DeleteAnswer)
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code) // Ошибка сервиса, возвращаем 500
	mockService.AssertExpectations(t)
}

func TestDeleteAnswerHandlerInvalidID(t *testing.T) {
	mockService := new(MockService)
	logger := logrus.New()
	handler := NewHandler(mockService, logger)

	req := httptest.NewRequest(http.MethodDelete, "/answers/abc", nil) // Некорректный ID
	rr := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Delete("/answers/{id}", handler.DeleteAnswer)
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockService.AssertNotCalled(t, "DeleteAnswer", mock.Anything)
}
