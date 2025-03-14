package models

import (
	"time"
)

type Student struct {
	ID                    int       `json:"id" gorm:"primaryKey"`
	FirstName             string    `json:"first_name" gorm:"size:100"`
	LastName              string    `json:"last_name" gorm:"size:100"`
	DOB                   time.Time `json:"dob"`
	Gender                string    `json:"gender" gorm:"size:10"`
	AddressLine1          string    `json:"address_line1" gorm:"size:255"`
	AddressLine2          string    `json:"address_line2" gorm:"size:255"`
	City                  string    `json:"city" gorm:"size:100"`
	ContactNumber         string    `json:"contact_number" gorm:"size:15"`
	ContactNumberGuardian string    `json:"contact_number_guardian" gorm:"size:15"`
	SupervisorID          int       `json:"supervisor_id"`
	EmployerID            int       `json:"employer_id"`
	Remarks               string    `json:"remarks" gorm:"size:500"`
	Photo                 []byte    `json:"photo"`
	HomeCoordinates       string    `json:"home_coordinates" gorm:"size:100"`
}

func (Student) TableName() string {
	return "student"
}
