package models

import (
	"time"
)

type Attendance struct {
	ID                uint      `gorm:"primaryKey"`
	StudentID         int       `gorm:"not null"`
	CheckInDateTime   time.Time `gorm:"not null" json:"check_in_date_time"`
	CheckInLongitude  float64   `gorm:"type:double precision" json:"check_in_long"`
	CheckInLatitude   float64   `gorm:"type:double precision" json:"check_in_lat"`
	CheckOutDateTime  time.Time `json:"check_out_date_time"`
	CheckOutLongitude float64   `gorm:"type:double precision" json:"check_out_long"`
	CheckOutLatitude  float64   `gorm:"type:double precision" json:"check_out_lat"`
}

func (Attendance) TableName() string {
	return "attendance"
}
