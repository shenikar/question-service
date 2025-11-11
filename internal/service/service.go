package service

import (
	"fmt"

	"github.com/google/uuid"
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
	repo repository.Repository
}

// NewService создает новый экземпляр сервиса.
func NewService(repo repository.Repository) Service {
	return &questionAnswerService{repo: repo}
}

// CreateQuestion создает новый вопрос.
func (s *questionAnswerService) CreateQuestion(question *models.Question) error {
	return s.repo.CreateQuestion(question)
}

// GetQuestion получает вопрос по ID.
func (s *questionAnswerService) GetQuestion(id uint) (*models.Question, error) {
	return s.repo.GetQuestion(id)
}

// GetAllQuestions получает все вопросы.
func (s *questionAnswerService) GetAllQuestions() ([]models.Question, error) {
	return s.repo.GetAllQuestions()
}

// DeleteQuestion удаляет вопрос по ID.
func (s *questionAnswerService) DeleteQuestion(id uint) error {
	return s.repo.DeleteQuestion(id)
}

// CreateAnswer создает новый ответ.
func (s *questionAnswerService) CreateAnswer(questionID uint, answer *models.Answer) error {
	// Бизнес-логика: Нельзя создать ответ к несуществующему вопросу.
	_, err := s.repo.GetQuestion(questionID)
	if err != nil {
		return fmt.Errorf("question with ID %d not found: %w", questionID, err)
	}

	answer.QuestionID = questionID
	answer.UserID = uuid.New() // Бизнес-логика: ID пользователя генерируется здесь
	return s.repo.CreateAnswer(answer)
}

// GetAnswer получает ответ по ID.
func (s *questionAnswerService) GetAnswer(id uint) (*models.Answer, error) {
	return s.repo.GetAnswer(id)
}

// DeleteAnswer удаляет ответ по ID.
func (s *questionAnswerService) DeleteAnswer(id uint) error {
	return s.repo.DeleteAnswer(id)
}