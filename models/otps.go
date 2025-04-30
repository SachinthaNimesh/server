package models

import (
	"time"

	"gorm.io/gorm"
)

// OTP represents an OTP code for authentication
type OTP struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	StudentID int       `json:"student_id" gorm:"not null;column:student_id"`
	OTPCode   string    `json:"otp_code" gorm:"not null;column:otp_code"`
	CreatedAt time.Time `json:"created_at" gorm:"not null;column:created_at"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null;column:expires_at"`
	IsUsed    bool      `json:"is_used" gorm:"not null;column:is_used"`
}

// OTPRequest is used for OTP generation <@WEB DASHBOARD>
type OTPRequest struct {
	StudentID int `json:"student_id"`
}

// OTPResponse is returned after OTP generation <@WEB DASHBOARD>
type OTPResponse struct {
	StudentID int       `json:"student_id"`
	OTPCode   string    `json:"otp_code"`
	ExpiresAt time.Time `json:"expires_at"`
}

// OTPValidationRequest is used when validating OTP from a mobile app
type OTPValidationRequest struct {
	OTPCode string `json:"otp_code"`
}

// OTPValidationResponse is returned after successful OTP validation
type OTPValidationResponse struct {
	Success    bool   `json:"success"`
	StudentID  int    `json:"student_id"`
	SecretCode string `json:"secret_code,omitempty"`
	Message    string `json:"message,omitempty"`
}

func (OTP) TableName() string {
	return "otps"
}

func (otp *OTP) BeforeCreate(tx *gorm.DB) (err error) {
	otp.CreatedAt = time.Now()
	otp.ExpiresAt = otp.CreatedAt.Add(30 * time.Minute) // Set expiration to 30 minutes from creation
	return
}
