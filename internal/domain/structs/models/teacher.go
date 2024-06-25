package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Teacher struct {
	ID           uint   `gorm:"primaryKey"`
	Uuid         string `gorm:"not_null;unique;size:64"`
	UserID       uint   `gorm:"not_null"`
	DepartmentID uint   `gorm:"not_null"`
	Nip          string `gorm:"not_null;unique;size:16"`
	CreatedAt    int64  `gorm:"autoCreateTime:milli"`
	UpdatedAt    int64  `gorm:"autoCreateTime:milli;autoUpdateTime:milli"`
	User         *User
	Department   *Department
}

func (m *Teacher) BeforeCreate(tx *gorm.DB) error {
	m.Uuid = uuid.NewString()
	return nil
}
