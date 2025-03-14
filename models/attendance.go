package models

import (
	"time"

	"github.com/twpayne/go-geom"
)

// Attendance represents the attendance record for a student
type Attendance struct {
	ID               int        `json:"id" gorm:"primary_key"`
	StudentID        int        `json:"student_id"`
	CheckInDateTime  time.Time  `json:"check_in_date_time"`
	CheckOutDateTime time.Time  `json:"check_out_date_time"`
	CheckInLocation  geom.Point `json:"check_in_location" gorm:"type:geometry(POINT,4326)"`
	CheckOutLocation geom.Point `json:"check_out_location" gorm:"type:geometry(POINT,4326)"`
}

func (Attendance) TableName() string {
	return "attendance"
}
