package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Note struct {
	ID        uint   `gorm:"primaryKey"`
	Uuid      string `gorm:"not_null;unique;size:64"`
	ClassID   uint   `gorm:"not_null"`
	Title     string `gorm:"not_null"`
	Body      string `gorm:"not_null"`
	CreatedAt int64  `gorm:"autoCreateTime:milli"`
	UpdatedAt int64  `gorm:"autoCreateTime:milli;autoUpdateTime:milli"`
	Class     *Class
}

func (m *Note) BeforeCreate(tx *gorm.DB) error {
	m.Uuid = uuid.NewString()
	return nil
}
