package service

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/shenikar/question-service/internal/models"
	"github.com/shenikar/question-service/internal/repository"
)

// Service определяет интерфейс для бизнес-логики приложения.
type Service interface {
	CreateQuestion(question *models.Question) error
	GetQuestion(id uint) (*models.Question, error)
	GetAllQuestions() ([]models.Question, error)
	DeleteQuestion(id uint) error
	CreateAnswer(questionID uint, answer *models.Answer) error
	GetAnswer(id uint) (*models.Answer, error)
	DeleteAnswer(id uint) error
}

// questionAnswerService - реализация Service.
type questionAnswerService struct {
	repo   repository.Repository
	logger *logrus.Logger
}

// NewService создает новый экземпляр сервиса.
func NewService(repo repository.Repository, logger *logrus.Logger) Service {
	return &questionAnswerService{repo: repo, logger: logger}
}

// CreateQuestion создает новый вопрос.
func (s *questionAnswerService) CreateQuestion(question *models.Question) error {
	s.logger.Debugf("Creating question: %+v", question)
	return s.repo.CreateQuestion(question)
}

// GetQuestion получает вопрос по ID.
func (s *questionAnswerService) GetQuestion(id uint) (*models.Question, error) {
	s.logger.Debugf("Getting question with ID: %d", id)
	return s.repo.GetQuestion(id)
}

// GetAllQuestions получает все вопросы.
func (s *questionAnswerService) GetAllQuestions() ([]models.Question, error) {
	s.logger.Debug("Getting all questions")
	return s.repo.GetAllQuestions()
}

// DeleteQuestion удаляет вопрос по ID.
func (s *questionAnswerService) DeleteQuestion(id uint) error {
	s.logger.Debugf("Deleting question with ID: %d", id)
	return s.repo.DeleteQuestion(id)
}

// CreateAnswer создает новый ответ.
func (s *questionAnswerService) CreateAnswer(questionID uint, answer *models.Answer) error {
	s.logger.Debugf("Creating answer for question ID %d: %+v", questionID, answer)
	// Бизнес-логика: Нельзя создать ответ к несуществующему вопросу.
	_, err := s.repo.GetQuestion(questionID)
	if err != nil {
		s.logger.Warnf("Attempted to create answer for non-existent question ID %d", questionID)
		return fmt.Errorf("question with ID %d not found: %w", questionID, err)
	}

	answer.QuestionID = questionID
	answer.UserID = uuid.New() // Бизнес-логика: ID пользователя генерируется здесь
	return s.repo.CreateAnswer(answer)
}

// GetAnswer получает ответ по ID.
func (s *questionAnswerService) GetAnswer(id uint) (*models.Answer, error) {
	s.logger.Debugf("Getting answer with ID: %d", id)
	return s.repo.GetAnswer(id)
}

// DeleteAnswer удаляет ответ по ID.
func (s *questionAnswerService) DeleteAnswer(id uint) error {
	s.logger.Debugf("Deleting answer with ID: %d", id)
	return s.repo.DeleteAnswer(id)
}