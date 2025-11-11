package service

import (
	"errors"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/shenikar/question-service/internal/models"
)

// MockRepository - мок для интерфейса repository.Repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateQuestion(question *models.Question) error {
	args := m.Called(question)
	return args.Error(0)
}

func (m *MockRepository) GetQuestion(id uint) (*models.Question, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Question), args.Error(1)
}

func (m *MockRepository) GetAllQuestions() ([]models.Question, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Question), args.Error(1)
}

func (m *MockRepository) DeleteQuestion(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRepository) CreateAnswer(answer *models.Answer) error {
	args := m.Called(answer)
	return args.Error(0)
}

func (m *MockRepository) GetAnswer(id uint) (*models.Answer, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Answer), args.Error(1)
}

func (m *MockRepository) DeleteAnswer(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreateQuestionService(t *testing.T) {
	mockRepo := new(MockRepository)
	logger := logrus.New()
	service := NewService(mockRepo, logger)

	question := &models.Question{
		Text: "Test Question",
	}

	mockRepo.On("CreateQuestion", question).Return(nil)

	err := service.CreateQuestion(question)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetQuestionService(t *testing.T) {
	mockRepo := new(MockRepository)
	logger := logrus.New()
	service := NewService(mockRepo, logger)

	expectedQuestion := &models.Question{
		ID:        1,
		Text:      "Test Question",
		CreatedAt: time.Now(),
	}

	mockRepo.On("GetQuestion", uint(1)).Return(expectedQuestion, nil)

	question, err := service.GetQuestion(1)
	assert.NoError(t, err)
	assert.NotNil(t, question)
	assert.Equal(t, expectedQuestion.ID, question.ID)
	mockRepo.AssertExpectations(t)
}

func TestCreateAnswerService(t *testing.T) {
	mockRepo := new(MockRepository)
	logger := logrus.New()
	service := NewService(mockRepo, logger)

	questionID := uint(1)
	answer := &models.Answer{
		Text: "Test Answer",
	}
	expectedQuestion := &models.Question{
		ID:        questionID,
		Text:      "Existing Question",
		CreatedAt: time.Now(),
	}

	// Ожидаем, что сервис сначала проверит существование вопроса
	mockRepo.On("GetQuestion", questionID).Return(expectedQuestion, nil)
	// Затем ожидаем создание ответа
	mockRepo.On("CreateAnswer", mock.AnythingOfType("*models.Answer")).Return(nil)

	err := service.CreateAnswer(questionID, answer)
	assert.NoError(t, err)
	assert.Equal(t, questionID, answer.QuestionID)
	assert.NotNil(t, answer.UserID) // Проверяем, что UserID был сгенерирован
	mockRepo.AssertExpectations(t)
}

func TestCreateAnswerServiceQuestionNotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	logger := logrus.New()
	service := NewService(mockRepo, logger)

	questionID := uint(1)
	answer := &models.Answer{
		Text: "Test Answer",
	}

	// Ожидаем, что сервис проверит существование вопроса и вернет ошибку
	mockRepo.On("GetQuestion", questionID).Return(nil, errors.New("not found"))

	err := service.CreateAnswer(questionID, answer)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "question with ID 1 not found")
	mockRepo.AssertExpectations(t)
}

func TestGetAllQuestionsService(t *testing.T) {
	mockRepo := new(MockRepository)
	logger := logrus.New()
	service := NewService(mockRepo, logger)

	expectedQuestions := []models.Question{
		{ID: 1, Text: "Q1"},
		{ID: 2, Text: "Q2"},
	}

	mockRepo.On("GetAllQuestions").Return(expectedQuestions, nil)

	questions, err := service.GetAllQuestions()
	assert.NoError(t, err)
	assert.Len(t, questions, 2)
	assert.Equal(t, expectedQuestions[0].Text, questions[0].Text)
	mockRepo.AssertExpectations(t)
}

func TestDeleteQuestionService(t *testing.T) {
	mockRepo := new(MockRepository)
	logger := logrus.New()
	service := NewService(mockRepo, logger)

	mockRepo.On("DeleteQuestion", uint(1)).Return(nil)

	err := service.DeleteQuestion(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetAnswerService(t *testing.T) {
	mockRepo := new(MockRepository)
	logger := logrus.New()
	service := NewService(mockRepo, logger)

	expectedAnswer := &models.Answer{ID: 1, Text: "A1"}

	mockRepo.On("GetAnswer", uint(1)).Return(expectedAnswer, nil)

	answer, err := service.GetAnswer(1)
	assert.NoError(t, err)
	assert.NotNil(t, answer)
	assert.Equal(t, expectedAnswer.ID, answer.ID)
	mockRepo.AssertExpectations(t)
}

func TestDeleteAnswerService(t *testing.T) {
	mockRepo := new(MockRepository)
	logger := logrus.New()
	service := NewService(mockRepo, logger)

	mockRepo.On("DeleteAnswer", uint(1)).Return(nil)

	err := service.DeleteAnswer(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
