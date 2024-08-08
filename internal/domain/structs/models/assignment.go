package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Assignment struct {
	ID          uint   `gorm:"primaryKey"`
	Uuid        string `gorm:"not_null;unique;size:64"`
	ClassID     uint   `gorm:"not_null"`
	Title       string `gorm:"not_null"`
	Description string `gorm:"not_null"`
	Deadline    int64  `gorm:"not_null"`
	CreatedAt   int64  `gorm:"autoCreateTime:milli"`
	UpdatedAt   int64  `gorm:"autoCreateTime:milli;autoUpdateTime:milli"`
	Answers     *[]Answer
	Class       *Class
}

func (m *Assignment) BeforeCreate(tx *gorm.DB) error {
	m.Uuid = uuid.NewString()
	return nil
}
