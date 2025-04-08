package models

type Employer struct {
	ID            uint    `json:"id" gorm:"primary_key;auto_increment"`
	Name          string  `json:"name" gorm:"type:varchar(100);not null"`
	StudentID     int     `json:"student_id" gorm:"index"`
	ContactNumber string  `json:"contact_number" gorm:"type:varchar(20)"`
	AddressLine1  string  `json:"address_line1,omitempty" gorm:"type:varchar(255);null"` //omitempty and null added
	AddressLine2  string  `json:"address_line2,omitempty" gorm:"type:varchar(255);null"`
	AddressLine3  string  `json:"address_line3,omitempty" gorm:"type:varchar(255);null"`
	Longitude     float64 `json:"addr_long" gorm:"column:addr_long"`
	Latitude      float64 `json:"addr_lat" gorm:"column:addr_lat"`
}

func (Employer) TableName() string {
	return "employer"
}
