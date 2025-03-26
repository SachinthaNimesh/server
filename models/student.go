package models

import (
	"encoding/json"
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
