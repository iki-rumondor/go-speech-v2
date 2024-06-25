package models

import (
	"github.com/google/uuid"
	"github.com/iki-rumondor/go-speech/internal/domain/structs/response"
	"gorm.io/gorm"
)

// Status : 1 = Pending, 2 = Accept, 3 = Unaccept

type ClassRequest struct {
	ID        uint   `gorm:"primaryKey"`
	Uuid      string `gorm:"not_null;unique;size:64"`
	StudentID uint   `gorm:"not_null"`
	ClassID   uint   `gorm:"not_null"`
	Status    uint   `gorm:"not_null;size:1"`
	CreatedAt int64  `gorm:"autoCreateTime:milli"`
	UpdatedAt int64  `gorm:"autoCreateTime:milli;autoUpdateTime:milli"`
	Class     *Class
	Student   *Student
}

func (m *ClassRequest) BeforeCreate(tx *gorm.DB) error {
	if result := tx.First(&ClassRequest{}, "student_id = ? AND class_id = ? AND id != ?", m.ClassID, m.StudentID, m.ID).RowsAffected; result > 0 {
		return response.BADREQ_ERR("Anda Sudah Mendaftar Untuk Kelas Ini")
	}

	m.Uuid = uuid.NewString()
	return nil
}
