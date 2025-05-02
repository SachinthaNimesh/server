package models

import "time"

// EmployeeTable represents the model for the query result
type EmployeeModel struct {
	StudentID       int        `json:"student_id" gorm:"column:student_id"`
	StudentName     string     `json:"student_name" gorm:"column:student_name"`
	StudentContact  string     `json:"student_contact" gorm:"column:student_contact"`
	EmployerID      *int       `json:"employer_id,omitempty" gorm:"column:employer_id"`
	EmployerName    *string    `json:"employer_name,omitempty" gorm:"column:employer_name"`
	EmployerContact *string    `json:"employer_contact,omitempty" gorm:"column:employer_contact"`
	EmployerAddress *string    `json:"employer_address,omitempty" gorm:"column:employer_address"`
	SupervisorID    *int       `json:"supervisor_id,omitempty" gorm:"column:supervisor_id"`
	SupervisorName  *string    `json:"supervisor_name,omitempty" gorm:"column:supervisor_name"`
	LatestOTPCode   *string    `json:"latest_otp_code,omitempty" gorm:"column:latest_otp_code"`
	ExpiresAt       *time.Time `json:"expires_at,omitempty" gorm:"column:expires_at"`
}
