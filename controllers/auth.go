package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"server/database"
	"server/models"
	"time"

	"github.com/gorilla/mux"
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

// HandleGenerateOTP godoc
// @Summary Generate OTP for a student
// @Description Generate a new OTP for a student
// @Tags authentication
// @Accept json
// @Produce json
// @Param request body struct { StudentID int `json:"student_id"` } true "Student ID"
// @Success 200 {object} models.OTPResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /generate-otp [post]
func (s *AuthService) HandleGenerateOTP(w http.ResponseWriter, r *http.Request) {
	var req struct {
		StudentID int `json:"student_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	resp, err := s.GenerateOTP(req.StudentID)
	if err != nil {
		log.Printf("Error generating OTP: %v", err)
		http.Error(w, "Failed to generate OTP", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// HandleValidateOTP godoc
// @Summary Validate OTP
// @Description Validate an OTP and generate a secret code
// @Tags authentication
// @Accept json
// @Produce json
// @Param request body struct { OTPCode string `json:"otp_code"` } true "OTP Code"
// @Success 200 {object} models.OTPValidationResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /validate-otp [post]
func (s *AuthService) HandleValidateOTP(w http.ResponseWriter, r *http.Request) {
	var req struct {
		OTPCode string `json:"otp_code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	resp, err := s.ValidateOTP(req.OTPCode)
	if err != nil {
		log.Printf("Error validating OTP: %v", err)
		http.Error(w, "Failed to validate OTP", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// HandleVerifyDeviceAuth godoc
// @Summary Verify device authorization
// @Description Verify if a device is authorized using student ID and secret code
// @Tags authentication
// @Accept json
// @Produce json
// @Param request body struct { StudentID int `json:"student_id"` SecretCode string `json:"secret_code"` } true "Authorization details"
// @Success 200 {object} map[string]bool
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /verify-device-auth [post]
func (s *AuthService) HandleVerifyDeviceAuth(w http.ResponseWriter, r *http.Request) {
	var req struct {
		StudentID  int    `json:"student_id"`
		SecretCode string `json:"secret_code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	isAuthorized, err := s.VerifyDeviceAuth(req.StudentID, req.SecretCode)
	if err != nil {
		log.Printf("Error verifying device authorization: %v", err)
		http.Error(w, "Failed to verify device authorization", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]bool{"authorized": isAuthorized}); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
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
		log.Printf("Database error while checking student existence: %v", err)
		return nil, fmt.Errorf("database error: %w", err)
	}
	if count == 0 {
		return nil, errors.New("student not found")
	}

	// Generate a random 4-digit OTP
	otp, err := s.generateRandomOTP(4)
	if err != nil {
		log.Printf("Error generating random OTP: %v", err)
		return nil, fmt.Errorf("failed to generate OTP: %w", err)
	}

	// Set expiration time (30 minutes from now)
	expiresAt := time.Now().Add(30 * time.Minute)

	// Invalidate any existing unused OTPs for this student
	err = s.db.Model(&models.OTP{}).
		Where("student_id = ? AND is_used = false", studentID).
		Update("is_used", true).Error
	if err != nil {
		log.Printf("Error invalidating existing OTPs: %v", err)
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
		log.Printf("Error storing new OTP: %v", err)
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
		log.Printf("OTP not found for code: %s", otpCode) // Improved logging
		return &models.OTPValidationResponse{
			Success: false,
			Message: "Invalid OTP",
		}, nil
	} else if err != nil {
		log.Printf("Database error while fetching OTP for code %s: %v", otpCode, err) // Improved logging
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Check if OTP is already used
	if otp.IsUsed {
		log.Printf("OTP already used for code: %s", otpCode) // Improved logging
		return &models.OTPValidationResponse{
			Success: false,
			Message: "OTP has already been used",
		}, nil
	}

	// Check if OTP is expired
	if time.Now().After(otp.ExpiresAt) {
		log.Printf("OTP expired for code: %s", otpCode) // Improved logging
		otp.IsUsed = true
		if err := s.db.Save(&otp).Error; err != nil {
			log.Printf("Error marking expired OTP as used for code %s: %v", otpCode, err) // Improved logging
		}
		return &models.OTPValidationResponse{
			Success: false,
			Message: "OTP has expired",
		}, nil
	}

	// Generate a secret code for the device
	secretCode, err := s.generateSecretCode()
	if err != nil {
		log.Printf("Error generating secret code for OTP code %s: %v", otpCode, err) // Improved logging
		return nil, fmt.Errorf("failed to generate secret code: %w", err)
	}

	// Mark OTP as used
	otp.IsUsed = true
	if err := s.db.Save(&otp).Error; err != nil {
		log.Printf("Error marking OTP as used for code %s: %v", otpCode, err) // Improved logging
	}

	// Store the secret code associated with the student_id
	authDevice := models.AuthorizedDevice{
		StudentID:  otp.StudentID,
		SecretCode: secretCode,
	}
	if err := s.db.Create(&authDevice).Error; err != nil {
		log.Printf("Error storing device authorization for student ID %d: %v", otp.StudentID, err) // Improved logging
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
		log.Printf("Database error while verifying device authorization: %v", err)
		return false, fmt.Errorf("database error: %w", err)
	}

	return true, nil
}

// RegisterRoutes registers the routes for AuthService
func (s *AuthService) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/generate-otp", s.HandleGenerateOTP).Methods("POST")
	router.HandleFunc("/validate-otp", s.HandleValidateOTP).Methods("POST")
	router.HandleFunc("/verify-device-auth", s.HandleVerifyDeviceAuth).Methods("POST")
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
