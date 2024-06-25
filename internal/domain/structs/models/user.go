package models

import (
	"github.com/google/uuid"
	"github.com/iki-rumondor/go-speech/internal/domain/structs/response"
	"github.com/iki-rumondor/go-speech/internal/utils"
	"gorm.io/gorm"
)

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Uuid      string `gorm:"not_null;unique;size:64"`
	Name      string `gorm:"not_null;size:32"`
	Email     string `gorm:"not_null;unique;size:32"`
	Password  string `gorm:"not_null;size:64"`
	Active    bool   `gorm:"not_null"`
	RoleID    uint   `gorm:"not_null"`
	CreatedAt int64  `gorm:"autoCreateTime:milli"`
	UpdatedAt int64  `gorm:"autoCreateTime:milli;autoUpdateTime:milli"`
	Role      *Role
	Teacher   *Teacher
	Student   *Student
}

func (m *User) BeforeSave(tx *gorm.DB) error {
	if result := tx.First(&User{}, "email = ? AND id != ?", m.Email, m.ID).RowsAffected; result > 0 {
		return response.BADREQ_ERR("Email Yang Dipakai Sudah Terdaftar Pada Akun Lain")
	}
	return nil
}

func (m *User) BeforeCreate(tx *gorm.DB) error {
	hashPass, err := utils.HashPassword(m.Password)
	if err != nil {
		return err
	}
	m.Password = hashPass
	m.Uuid = uuid.NewString()
	return nil
}
