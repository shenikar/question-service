package repository

import (
	"gorm.io/gorm"
	"github.com/shenikar/question-service/internal/models"
)

// Repository определяет интерфейс для работы с хранилищем данных.
type Repository interface {
	CreateQuestion(question *models.Question) error
	GetQuestion(id uint) (*models.Question, error)
	GetQuestions() ([]models.Question, error)
	DeleteQuestion(id uint) error
	CreateAnswer(answer *models.Answer) error
	GetAnswer(id uint) (*models.Answer, error)
	DeleteAnswer(id uint) error
}

// dbRepository - реализация Repository для работы с базой данных.
type dbRepository struct {
	db *gorm.DB
}

// NewRepository создает новый экземпляр репозитория.
func NewRepository(db *gorm.DB) Repository {
	return &dbRepository{db: db}
}

// CreateQuestion создает новый вопрос в базе данных.
func (r *dbRepository) CreateQuestion(question *models.Question) error {
	return r.db.Create(question).Error
}

// GetQuestion получает вопрос из базы данных по его ID.
func (r *dbRepository) GetQuestion(id uint) (*models.Question, error) {
	var question models.Question
	err := r.db.Preload("Answers").First(&question, id).Error
	return &question, err
}

// CreateAnswer создает новый ответ в базе данных.
func (r *dbRepository) CreateAnswer(answer *models.Answer) error {
	// Проверяем, существует ли вопрос
	var question models.Question
	if err := r.db.First(&question, answer.QuestionID).Error; err != nil {
		return err // Возвращаем ошибку, если вопрос не найден
	}
	return r.db.Create(answer).Error
}

// GetQuestions получает все вопросы из базы данных.
func (r *dbRepository) GetQuestions() ([]models.Question, error) {
	var questions []models.Question
	err := r.db.Preload("Answers").Find(&questions).Error
	return questions, err
}

// DeleteQuestion удаляет вопрос из базы данных по его ID.
func (r *dbRepository) DeleteQuestion(id uint) error {
	return r.db.Delete(&models.Question{}, id).Error
}

// GetAnswer получает ответ из базы данных по его ID.
func (r *dbRepository) GetAnswer(id uint) (*models.Answer, error) {
	var answer models.Answer
	err := r.db.First(&answer, id).Error
	return &answer, err
}

// DeleteAnswer удаляет ответ из базы данных по его ID.
func (r *dbRepository) DeleteAnswer(id uint) error {
	return r.db.Delete(&models.Answer{}, id).Error
}