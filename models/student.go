package models

import (
	"time"
)

type Student struct {
	ID                    int     `json:"id" gorm:"primaryKey"`
	FirstName             string  `json:"first_name" gorm:"size:100"`
	LastName              string  `json:"last_name" gorm:"size:100"`
	DOB                   Date    `json:"dob" gorm:"type:date"` // Custom Date type
	Gender                string  `json:"gender" gorm:"size:10"`
	AddressLine1          string  `json:"address_line1" gorm:"size:255"`
	AddressLine2          string  `json:"address_line2" gorm:"size:255"`
	City                  string  `json:"city" gorm:"size:100"`
	ContactNumber         string  `json:"contact_number" gorm:"size:15"`
	ContactNumberGuardian string  `json:"contact_number_guardian" gorm:"size:15"`
	SupervisorID          int     `json:"supervisor_id"`
	EmployerID            int     `json:"employer_id"`
	Remarks               string  `json:"remarks" gorm:"size:500"`
	Photo                 []byte  `json:"photo"`
	HomeLong              float64 `json:"home_long" gorm:"default:0"`
	HomeLat               float64 `json:"home_lat" gorm:"default:0"`
}

func (Student) TableName() string {
	return "student"
}

// Custom Date type to handle "YYYY-MM-DD" format
type Date struct {
	time.Time
}

const dateFormat = "2006-01-02"

// UnmarshalJSON for custom Date type
func (d *Date) UnmarshalJSON(b []byte) error {
	str := string(b)
	parsedTime, err := time.Parse(`"`+dateFormat+`"`, str)
	if err != nil {
		return err
	}
	d.Time = parsedTime
	return nil
}

// MarshalJSON for custom Date type
func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.Time.Format(dateFormat) + `"`), nil
}

// Scan for database compatibility
func (d *Date) Scan(value interface{}) error {
	if dateStr, ok := value.(string); ok {
		parsedTime, err := time.Parse(dateFormat, dateStr)
		if err != nil {
			return err
		}
		d.Time = parsedTime
		return nil
	}
	return nil
}

// Value for database compatibility
func (d Date) Value() (interface{}, error) {
	return d.Time.Format(dateFormat), nil
}

/*
package models

import (
	"time"
)

type Student struct {
	ID                    int       `json:"id" gorm:"primaryKey"`
	FirstName             string    `json:"first_name" gorm:"size:100"`
	LastName              string    `json:"last_name" gorm:"size:100"`
	DOB                   time.Time `json:"dob" gorm:"type:date"`
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
