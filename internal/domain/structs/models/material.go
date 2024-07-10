package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Material struct {
	ID          uint   `gorm:"primaryKey"`
	Uuid        string `gorm:"not_null;unique;size:64"`
	ClassID     uint   `gorm:"not_null"`
	VideoID     uint   `gorm:"not_null"`
	BookID      uint   `gorm:"not_null"`
	Title       string `gorm:"not_null"`
	Description string `gorm:"not_null"`
	CreatedAt   int64  `gorm:"autoCreateTime:milli"`
	UpdatedAt   int64  `gorm:"autoCreateTime:milli;autoUpdateTime:milli"`
	Class       *Class
	Book        *BookPart
	Video       *VideoPart
}

func (m *Material) BeforeCreate(tx *gorm.DB) error {
	m.Uuid = uuid.NewString()
	return nil
}
