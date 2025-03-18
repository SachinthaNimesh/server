package models

import (
	"time"
)

type Attendance struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	StudentID         int       `json:"student_id" gorm:"not null"`
	CheckInDateTime   time.Time `json:"check_in_date_time" gorm:"not null"`
	CheckInLongitude  float64   `json:"check_in_long" gorm:"type:double precision"`
	CheckInLatitude   float64   `json:"check_in_lat" gorm:"type:double precision"`
	CheckOutDateTime  time.Time `json:"check_out_date_time"`
	CheckOutLongitude float64   `json:"check_out_long" gorm:"type:double precision"`
	CheckOutLatitude  float64   `json:"check_out_lat" gorm:"type:double precision"`
}

func (Attendance) TableName() string {
	return "attendance"
}
