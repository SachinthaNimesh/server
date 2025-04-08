package models

import (
	"encoding/json"
	"time"
)

type Student struct {
	ID                    int64     `gorm:"primaryKey;column:id" json:"id"`
	FirstName             string    `gorm:"column:first_name" json:"first_name"`
	LastName              string    `gorm:"column:last_name" json:"last_name"`
	DOB                   time.Time `gorm:"column:dob" json:"dob"`
	Gender                string    `gorm:"column:gender" json:"gender"`
	AddressLine1          string    `gorm:"column:address_line1" json:"address_line1"`
	AddressLine2          string    `gorm:"column:address_line2" json:"address_line2"`
	City                  string    `gorm:"column:city" json:"city"`
	ContactNumber         string    `gorm:"column:contact_number" json:"contact_number"`
	ContactNumberGuardian string    `gorm:"column:contact_number_guardian" json:"contact_number_guardian"`
	SupervisorID          int64     `gorm:"column:supervisor_id" json:"supervisor_id"`
	EmployerID            int64     `gorm:"column:employer_id" json:"employer_id"`
	Remarks               string    `gorm:"column:remarks" json:"remarks"`
	Photo                 []byte    `gorm:"column:photo" json:"photo,omitempty"`
	HomeLong              float64   `gorm:"column:home_long" json:"home_long"`
	HomeLat               float64   `gorm:"column:home_lat" json:"home_lat"`
	IMEI                  string    `gorm:"column:imei" json:"imei"`
}

/*
	type Student struct {
		ID                    int       `json:"id" gorm:"primaryKey"`
		FirstName             string    `json:"first_name" gorm:"size:100"`
		LastName              string    `json:"last_name" gorm:"size:100"`
		DOB                   time.Time `json:"dob" gorm:"type:timestamp"`
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
		HomeLong              float64   `json:"home_long" gorm:"default:0"`
		HomeLat               float64   `json:"home_lat" gorm:"default:0"`
	}
*/
func (Student) TableName() string {
	return "student"
}

// UnmarshalJSON provides flexible date parsing for the DOB field
func (s *Student) UnmarshalJSON(data []byte) error {
	// Temporary struct to avoid recursive unmarshal
	type Alias Student
	aux := &struct {
		DOB string `json:"dob"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}

	// First, unmarshal the entire JSON
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Try parsing the DOB with multiple formats
	var err error
	dateFormats := []string{
		"2006-01-02",                // Simple date
		"2006-01-02T15:04:05Z07:00", // Full timestamp
		time.RFC3339,                // Another timestamp format
		"2006-01-02 15:04:05",       // Timestamp without timezone
	}

	for _, format := range dateFormats {
		s.DOB, err = time.Parse(format, aux.DOB)
		if err == nil {
			return nil
		}
	}

	// If no format works, return the last error
	return err
}

/*package models

import (
	"time"
)

type Student struct {
	ID                    int       `json:"id" gorm:"primaryKey"`
	FirstName             string    `json:"first_name" gorm:"size:100"`
	LastName              string    `json:"last_name" gorm:"size:100"`
	DOB                   time.Time `json:"dob" gorm:"type:timestamp"`
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
	HomeLong              float64   `json:"home_long" gorm:"default:0"`
	HomeLat               float64   `json:"home_lat" gorm:"default:0"`
}

func (Student) TableName() string {
	return "student"
}
*/
