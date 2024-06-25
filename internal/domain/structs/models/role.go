package models

type Role struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"not_null;size:16"`
	User *[]User
}
