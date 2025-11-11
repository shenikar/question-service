package repository

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/shenikar/question-service/internal/models"
)

func newMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	sqldb, mock, err := sqlmock.New()
	assert.NoError(t, err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqldb,
	}), &gorm.Config{})
	assert.NoError(t, err)

	return gormDB, mock
}

func TestCreateQuestion(t *testing.T) {
	gormDB, mock := newMockDB(t)
	repo := NewRepository(gormDB, logrus.New())

	question := &models.Question{
		Text: "Test Question",
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "questions"`).
		WithArgs(question.Text, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(1, time.Now()))
	mock.ExpectCommit()

	err := repo.CreateQuestion(question)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), question.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetQuestion(t *testing.T) {
	gormDB, mock := newMockDB(t)
	repo := NewRepository(gormDB, logrus.New())

	expectedQuestion := &models.Question{
		ID:        1,
		Text:      "Test Question",
		CreatedAt: time.Now(),
	}

	mock.ExpectQuery(
		`SELECT \* FROM "questions" WHERE "questions"."id" = \$1 ORDER BY "questions"."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "text", "created_at"}).
			AddRow(expectedQuestion.ID, expectedQuestion.Text, expectedQuestion.CreatedAt))

	mock.ExpectQuery(
		`SELECT \* FROM "answers" WHERE "answers"."question_id" = \$1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "question_id", "user_id", "text", "created_at"})) // пустой результат

	question, err := repo.GetQuestion(1)
	assert.NoError(t, err)
	assert.NotNil(t, question)
	assert.Equal(t, expectedQuestion.ID, question.ID)
	assert.Equal(t, expectedQuestion.Text, question.Text)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetQuestionNotFound(t *testing.T) {
	gormDB, mock := newMockDB(t)
	repo := NewRepository(gormDB, logrus.New())

	mock.ExpectQuery(
		`SELECT \* FROM "questions" WHERE "questions"."id" = \$1 ORDER BY "questions"."id" LIMIT \$2`). // <-- Изменено
		WithArgs(999, 1).                                                                               // <-- Добавлен аргумент для LIMIT
		WillReturnError(gorm.ErrRecordNotFound)                                                         // Возвращаем ошибку GORM

	question, err := repo.GetQuestion(999)
	assert.Error(t, err)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	assert.NotNil(t, question)            // GORM возвращает пустой объект, не nil
	assert.Equal(t, uint(0), question.ID) // Проверяем, что ID дефолтный
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllQuestions(t *testing.T) {
	gormDB, mock := newMockDB(t)
	repo := NewRepository(gormDB, logrus.New())

	q1 := models.Question{ID: 1, Text: "Q1", CreatedAt: time.Now()}
	q2 := models.Question{ID: 2, Text: "Q2", CreatedAt: time.Now()}

	mock.ExpectQuery(
		`SELECT \* FROM "questions"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "text", "created_at"}).
			AddRow(q1.ID, q1.Text, q1.CreatedAt).
			AddRow(q2.ID, q2.Text, q2.CreatedAt))

	mock.ExpectQuery(
		`SELECT \* FROM "answers" WHERE "answers"."question_id" IN \(\$1,\$2\)`).
		WithArgs(1, 2).
		WillReturnRows(sqlmock.NewRows([]string{"id", "question_id", "user_id", "text", "created_at"})) // пустой результат

	questions, err := repo.GetAllQuestions()
	assert.NoError(t, err)
	assert.Len(t, questions, 2)
	assert.Equal(t, q1.Text, questions[0].Text)
	assert.Equal(t, q2.Text, questions[1].Text)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteQuestion(t *testing.T) {
	gormDB, mock := newMockDB(t)
	repo := NewRepository(gormDB, logrus.New())

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "questions" WHERE "questions"."id" = \$1`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.DeleteQuestion(1)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateAnswer(t *testing.T) {
	gormDB, mock := newMockDB(t)
	repo := NewRepository(gormDB, logrus.New())

	answer := &models.Answer{
		QuestionID: 1,
		UserID:     uuid.New(),
		Text:       "Test Answer",
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "answers"`).
		WithArgs(answer.QuestionID, sqlmock.AnyArg(), answer.Text, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(1, time.Now()))
	mock.ExpectCommit()

	err := repo.CreateAnswer(answer)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), answer.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAnswer(t *testing.T) {
	gormDB, mock := newMockDB(t)
	repo := NewRepository(gormDB, logrus.New())

	expectedAnswer := &models.Answer{
		ID:         1,
		QuestionID: 1,
		UserID:     uuid.New(),
		Text:       "Test Answer",
		CreatedAt:  time.Now(),
	}

	mock.ExpectQuery(
		`SELECT \* FROM "answers" WHERE "answers"."id" = \$1 ORDER BY "answers"."id" LIMIT \$2`). // <-- Изменено
		WithArgs(1, 1).                                                                           // <-- Добавлен аргумент для LIMIT
		WillReturnRows(sqlmock.NewRows([]string{"id", "question_id", "user_id", "text", "created_at"}).
			AddRow(expectedAnswer.ID, expectedAnswer.QuestionID, expectedAnswer.UserID, expectedAnswer.Text, expectedAnswer.CreatedAt))

	answer, err := repo.GetAnswer(1)
	assert.NoError(t, err)
	assert.NotNil(t, answer)
	assert.Equal(t, expectedAnswer.ID, answer.ID)
	assert.Equal(t, expectedAnswer.Text, answer.Text)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteAnswer(t *testing.T) {
	gormDB, mock := newMockDB(t)
	repo := NewRepository(gormDB, logrus.New())

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "answers" WHERE "answers"."id" = \$1`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.DeleteAnswer(1)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
