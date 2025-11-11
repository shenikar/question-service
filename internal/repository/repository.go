package repository

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/shenikar/question-service/internal/models"
)

// Repository определяет интерфейс для работы с хранилищем данных.
type Repository interface {
	CreateQuestion(question *models.Question) error
	GetQuestion(id uint) (*models.Question, error)
	GetAllQuestions() ([]models.Question, error)
	DeleteQuestion(id uint) error
	CreateAnswer(answer *models.Answer) error
	GetAnswer(id uint) (*models.Answer, error)
	DeleteAnswer(id uint) error
}

// dbRepository - реализация Repository для работы с базой данных.
type dbRepository struct {
	db     *gorm.DB
	logger *logrus.Logger
}

// NewRepository создает новый экземпляр репозитория.
func NewRepository(db *gorm.DB, logger *logrus.Logger) Repository {
	return &dbRepository{db: db, logger: logger}
}

// CreateQuestion создает новый вопрос в базе данных.
func (r *dbRepository) CreateQuestion(question *models.Question) error {
	r.logger.Debugf("Creating question: %+v", question)
	return r.db.Create(question).Error
}

// GetQuestion получает вопрос из базы данных по его ID.
func (r *dbRepository) GetQuestion(id uint) (*models.Question, error) {
	r.logger.Debugf("Getting question with ID: %d", id)
	var question models.Question
	err := r.db.Preload("Answers").First(&question, id).Error
	return &question, err
}

// CreateAnswer создает новый ответ в базе данных.
func (r *dbRepository) CreateAnswer(answer *models.Answer) error {
	r.logger.Debugf("Creating answer: %+v", answer)
	return r.db.Create(answer).Error
}

// GetAllQuestions получает все вопросы из базы данных.
func (r *dbRepository) GetAllQuestions() ([]models.Question, error) {
	r.logger.Debug("Getting all questions")
	var questions []models.Question
	err := r.db.Preload("Answers").Find(&questions).Error
	return questions, err
}

// DeleteQuestion удаляет вопрос из базы данных по его ID.
func (r *dbRepository) DeleteQuestion(id uint) error {
	r.logger.Debugf("Deleting question with ID: %d", id)
	return r.db.Delete(&models.Question{}, id).Error
}

// GetAnswer получает ответ из базы данных по его ID.
func (r *dbRepository) GetAnswer(id uint) (*models.Answer, error) {
	r.logger.Debugf("Getting answer with ID: %d", id)
	var answer models.Answer
	err := r.db.First(&answer, id).Error
	return &answer, err
}

// DeleteAnswer удаляет ответ из базы данных по его ID.
func (r *dbRepository) DeleteAnswer(id uint) error {
	r.logger.Debugf("Deleting answer with ID: %d", id)
	return r.db.Delete(&models.Answer{}, id).Error
}
