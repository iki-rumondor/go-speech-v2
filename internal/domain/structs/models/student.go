package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Student struct {
	ID            uint   `gorm:"primaryKey"`
	Uuid          string `gorm:"not_null;unique;size:64"`
	UserID        uint   `gorm:"not_null"`
	Nim           string `gorm:"not_null;unique;size:16"`
	CreatedAt     int64  `gorm:"autoCreateTime:milli"`
	UpdatedAt     int64  `gorm:"autoCreateTime:milli;autoUpdateTime:milli"`
	User          *User
	ClassRequests *[]ClassRequest
}

func (m *Student) BeforeCreate(tx *gorm.DB) error {
	m.Uuid = uuid.NewString()
	return nil
}
