package models

import (
	"time"
)

type Attendance struct {
	ID               uint      `json:"id"`
	StudentID        int       `json:"student_id"`
	CheckInDateTime  time.Time `json:"check_in_date_time"`
	CheckInLong      float64   `json:"check_in_long"`
	CheckInLat       float64   `json:"check_in_lat"`
	CheckOutDateTime time.Time `json:"check_out_date_time"`
	CheckOutLong     float64   `json:"check_out_long"`
	CheckOutLat      float64   `json:"check_out_lat"`
}
