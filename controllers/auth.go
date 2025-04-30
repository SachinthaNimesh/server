package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"server/database"
	"server/models"
	"time"

	"gorm.io/gorm"
)

// AuthService handles authentication-related operations
type AuthService struct {
	db *gorm.DB
}

// NewAuthService creates a new auth service
func NewAuthService() *AuthService {
	return &AuthService{
		db: database.DB, // Use the GORM DB instance
	}
}

// GenerateOTP creates a new OTP for a student
func (s *AuthService) GenerateOTP(studentID int) (*models.OTPResponse, error) {
	// Check if student exists
	var count int64
	err := s.db.Model(&models.Student{}).
		Where("id = ?", studentID). // Removed is_active check
		Count(&count).Error
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if count == 0 {
		return nil, errors.New("student not found")
	}

	// Generate a random 4-digit OTP
	otp, err := s.generateRandomOTP(4)
	if err != nil {
		return nil, fmt.Errorf("failed to generate OTP: %w", err)
	}

	// Set expiration time (30 minutes from now)
	expiresAt := time.Now().Add(30 * time.Minute)

	// Invalidate any existing unused OTPs for this student
	err = s.db.Model(&models.OTP{}).
		Where("student_id = ? AND is_used = false", studentID).
		Update("is_used", true).Error
	if err != nil {
		return nil, fmt.Errorf("failed to invalidate existing OTPs: %w", err)
	}

	// Insert new OTP
	newOTP := models.OTP{
		StudentID: studentID,
		OTPCode:   otp,
		ExpiresAt: expiresAt,
		IsUsed:    false,
	}
	err = s.db.Create(&newOTP).Error
	if err != nil {
		return nil, fmt.Errorf("failed to store OTP: %w", err)
	}

	return &models.OTPResponse{
		StudentID: studentID,
		OTPCode:   otp,
		ExpiresAt: expiresAt,
	}, nil
}

// ValidateOTP checks if an OTP is valid and returns student_id and a new secret code
func (s *AuthService) ValidateOTP(otpCode string) (*models.OTPValidationResponse, error) {
	var otp models.OTP
	err := s.db.Where("otp_code = ?", otpCode).First(&otp).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &models.OTPValidationResponse{
			Success: false,
			Message: "Invalid OTP",
		}, nil
	} else if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Check if OTP is already used
	if otp.IsUsed {
		return &models.OTPValidationResponse{
			Success: false,
			Message: "OTP has already been used",
		}, nil
	}

	// Check if OTP is expired
	if time.Now().After(otp.ExpiresAt) {
		otp.IsUsed = true
		s.db.Save(&otp)
		return &models.OTPValidationResponse{
			Success: false,
			Message: "OTP has expired",
		}, nil
	}

	// Generate a secret code for the device
	secretCode, err := s.generateSecretCode()
	if err != nil {
		return nil, fmt.Errorf("failed to generate secret code: %w", err)
	}

	// Mark OTP as used
	otp.IsUsed = true
	err = s.db.Save(&otp).Error
	if err != nil {
		return nil, fmt.Errorf("failed to mark OTP as used: %w", err)
	}

	// Store the secret code associated with the student_id
	authDevice := models.AuthorizedDevice{
		StudentID:  otp.StudentID,
		SecretCode: secretCode,
	}
	err = s.db.Create(&authDevice).Error
	if err != nil {
		return nil, fmt.Errorf("failed to store device authorization: %w", err)
	}

	return &models.OTPValidationResponse{
		Success:    true,
		StudentID:  otp.StudentID,
		SecretCode: secretCode,
		Message:    "Authentication successful",
	}, nil
}

// VerifyDeviceAuth verifies if a device is authorized using student_id and secret_code
func (s *AuthService) VerifyDeviceAuth(studentID int, secretCode string) (bool, error) {
	var authDevice models.AuthorizedDevice

	// Check if the device exists
	err := s.db.Where("student_id = ? AND secret_code = ?", studentID, secretCode).
		First(&authDevice).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("database error: %w", err)
	}

	// Removed is_active logic

	return true, nil
}

// Helper function to generate a random 4-digit OTP
func (s *AuthService) generateRandomOTP(digits int) (string, error) {
	maxNum := big.NewInt(0).Exp(big.NewInt(10), big.NewInt(int64(digits)), nil)
	n, err := rand.Int(rand.Reader, maxNum)
	if err != nil {
		return "", err
	}

	// Format with leading zeros
	return fmt.Sprintf("%0*d", digits, n), nil
}

// Helper function to generate a secure random secret code
func (s *AuthService) generateSecretCode() (string, error) {
	bytes := make([]byte, 32) // 256 bits of entropy
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
