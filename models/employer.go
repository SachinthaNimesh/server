package models

import (
	"github.com/twpayne/go-geom"
)

type Employer struct {
	ID                 uint       `gorm:"primary_key;auto_increment"`
	Name               string     `gorm:"type:varchar(100);not null"`
	StudentID          int        `gorm:"index"`
	ContactNumber      string     `gorm:"type:varchar(20)"`
	AddressLine1       string     `gorm:"type:varchar(255)"`
	AddressLine2       string     `gorm:"type:varchar(255)"`
	AddressLine3       string     `gorm:"type:varchar(255)"`
	AddressCoordinates geom.Point `gorm:"type:geometry;srid:4326"`
}

func (Employer) TableName() string {
	return "employer"
}
