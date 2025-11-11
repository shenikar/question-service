package models

import (
	"time"

	"github.com/google/uuid"
)

// Question представляет модель вопроса
type Question struct {
	ID        uint      `gorm:"primaryKey"`
	Text      string    `gorm:"not null" validate:"required,min=3,max=500"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	Answers   []Answer  `gorm:"foreignKey:QuestionID;constraint:OnDelete:CASCADE;"`
}

// Answer представляет модель ответа
type Answer struct {
	ID         uint      `gorm:"primaryKey"`
	QuestionID uint      `gorm:"not null"`
	UserID     uuid.UUID `gorm:"not null"`
	Text       string    `gorm:"not null" validate:"required,min=3,max=500"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}
