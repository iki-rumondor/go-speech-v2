package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Answer struct {
	ID           uint   `gorm:"primaryKey"`
	Uuid         string `gorm:"not_null;unique;size:64"`
	AssignmentID uint   `gorm:"not_null"`
	UserID       uint   `gorm:"not_null"`
	Filename     string `gorm:"not_null"`
	Grade        int    `gorm:"not_null"`
	Ontime       bool   `gorm:"not_null"`
	CreatedAt    int64  `gorm:"autoCreateTime:milli"`
	UpdatedAt    int64  `gorm:"autoCreateTime:milli;autoUpdateTime:milli"`
	Assignment   *Assignment
	User         *User
}

func (m *Answer) BeforeCreate(tx *gorm.DB) error {
	m.Uuid = uuid.NewString()
	return nil
}
