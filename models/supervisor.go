package models

type Supervisor struct {
	SupervisorID  int    `json:"supervisor_id" gorm:"primaryKey"`
	StudentID     int    `json:"student_id" gorm:"index"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	EmailAddress  string `json:"email_address" gorm:"unique"`
	ContactNumber string `json:"contact_number"`
}

func (Supervisor) TableName() string {
	return "supervisor"
}
