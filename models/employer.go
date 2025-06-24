package models

type Employer struct {
	ID            uint    `json:"id"`
	Name          string  `json:"name"`
	StudentID     int     `json:"student_id"`
	ContactNumber string  `json:"contact_number"`
	AddressLine1  string  `json:"address_line1,omitempty"`
	AddressLine2  string  `json:"address_line2,omitempty"`
	AddressLine3  string  `json:"address_line3,omitempty"`
	Longitude     float64 `json:"addr_long"`
	Latitude      float64 `json:"addr_lat"`
}
