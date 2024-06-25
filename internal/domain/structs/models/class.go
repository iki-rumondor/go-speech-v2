package models

import (
	"github.com/google/uuid"
	"github.com/iki-rumondor/go-speech/internal/domain/structs/response"
	"gorm.io/gorm"
)

type Class struct {
	ID            uint   `gorm:"primaryKey"`
	Uuid          string `gorm:"not_null;unique;size:64"`
	TeacherID     uint   `gorm:"not_null"`
	Name          string `gorm:"not_null;size:32"`
	Code          string `gorm:"not_null;unique;size:8"`
	CreatedAt     int64  `gorm:"autoCreateTime:milli"`
	UpdatedAt     int64  `gorm:"autoCreateTime:milli;autoUpdateTime:milli"`
	Teacher       *Teacher
	ClassRequests *[]ClassRequest
}

func (m *Class) BeforeCreate(tx *gorm.DB) error {
	m.Uuid = uuid.NewString()
	return nil
}

func (m *Class) BeforeSave(tx *gorm.DB) error {
	if result := tx.First(&Class{}, "code = ? AND id != ?", m.Code, m.ID).RowsAffected; result > 0 {
		return response.BADREQ_ERR("Kode Yang Dipakai Sudah Terdaftar Pada Kelas Lain")
	}
	return nil
}
