package models

type ReadNotification struct {
	ID             uint `gorm:"primaryKey"`
	NotificationID uint `gorm:"not_null"`
	StudentID      uint `gorm:"not_null"`
}
