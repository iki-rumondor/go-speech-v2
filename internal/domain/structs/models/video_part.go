package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VideoPart struct {
	ID           uint   `gorm:"primaryKey"`
	Uuid         string `gorm:"not_null;unique;size:64"`
	VideoName    string `gorm:"not_null"`
	SubtitleName string `gorm:"not_null"`
	CreatedAt    int64  `gorm:"autoCreateTime:milli"`
	UpdatedAt    int64  `gorm:"autoCreateTime:milli;autoUpdateTime:milli"`
}

func (m *VideoPart) BeforeCreate(tx *gorm.DB) error {
	m.Uuid = uuid.NewString()
	return nil
}
