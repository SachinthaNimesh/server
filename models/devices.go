package models

import "time"

type AuthorizedDevice struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	StudentID  int       `json:"student_id" gorm:"not null;column:student_id"`
	SecretCode string    `json:"secret_code" gorm:"not null;column:secret_code"`
	CreatedAt  time.Time `json:"created_at" gorm:"not null;column:created_at"`
}

func (AuthorizedDevice) TableName() string {
	return "authorized_devices"
}

type AuthRequest struct {
	StudentID  int    `json:"student_id"`
	SecretCode string `json:"secret_code"`
}
