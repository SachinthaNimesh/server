package models

import "time"

type StudentDetailedResponse struct {
	StudentID        int64     `json:"student_id"`
	FirstName        string    `json:"first_name"`
	LastName         string    `json:"last_name"`
	EmployerName     string    `json:"employer_name"`
	CheckInDateTime  time.Time `json:"check_in_date_time"`
	CheckOutDateTime time.Time `json:"check_out_date_time"`
	Emotion          string    `json:"emotion"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
