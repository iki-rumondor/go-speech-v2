package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Status : 1 = Pending, 2 = Accept, 3 = Unaccept

type StudentClasses struct {
	ID        uint   `gorm:"primaryKey"`
	Uuid      string `gorm:"not_null;unique;size:64"`
	StudentID uint   `gorm:"not_null"`
	ClassID   uint   `gorm:"not_null"`
	CreatedAt int64  `gorm:"autoCreateTime:milli"`
	UpdatedAt int64  `gorm:"autoCreateTime:milli;autoUpdateTime:milli"`
	Class     *Class
	Student   *Student
}

func (m *StudentClasses) BeforeCreate(tx *gorm.DB) error {

	m.Uuid = uuid.NewString()
	return nil
}
