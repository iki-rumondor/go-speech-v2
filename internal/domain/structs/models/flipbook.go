package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FlipBook struct {
	ID        uint   `gorm:"primaryKey"`
	Uuid      string `gorm:"not_null;unique;size:64"`
	BookID    uint   `gorm:"not_null"`
	URL       string `gorm:"not_null"`
	Thumbnail string `gorm:"not_null"`
	Pdf       string `gorm:"not_null"`
	CreatedAt int64  `gorm:"autoCreateTime:milli"`
	UpdatedAt int64  `gorm:"autoCreateTime:milli;autoUpdateTime:milli"`
	Book      *Book
}

func (m *FlipBook) BeforeCreate(tx *gorm.DB) error {
	m.Uuid = uuid.NewString()
	return nil
}
