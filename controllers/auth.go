package controllers

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"server/database"
	"server/models"
	"strconv"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// AuthService handles authentication-related operations
type AuthService struct {
	db *sql.DB
}

// NewAuthService creates a new auth service
func NewAuthService() *AuthService {
	return &AuthService{
		db: database.DB, // Use the sql.DB instance
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
	StudentIDHeader := r.Header.Get("student-id")
	if StudentIDHeader == "" {
		log.Println("Missing student-id header")
		http.Error(w, "Missing student-id header", http.StatusBadRequest)
		return
	}

	studentID, err := strconv.Atoi(StudentIDHeader)
	if err != nil {
		log.Printf("Invalid student-id header: %v", err)
		http.Error(w, "Invalid student-id header", http.StatusBadRequest)
		return
	}

	resp, err := s.GenerateOTP(studentID)
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
	OTPCodeHeader := r.Header.Get("otp-code")
	if OTPCodeHeader == "" {
		log.Println("Missing otp-code header")
		http.Error(w, "Missing otp-code header", http.StatusBadRequest)
		return
	}

	resp, err := s.ValidateOTP(OTPCodeHeader)
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
	err := s.db.QueryRow("SELECT COUNT(*) FROM students WHERE id = ?", studentID).Scan(&count)
	if err != nil {
		log.Printf("Database error while checking student existence: %v", err)
		return nil, fmt.Errorf("database error: %w", err)
	}
	if count == 0 {
		return nil, errors.New("student not found")
	}

	// Check if there's an existing unused OTP that hasn't expired
	var existingOTP models.OTP
	err = s.db.QueryRow("SELECT otp_code, expires_at FROM otps WHERE student_id = ? AND is_used = false AND expires_at > CURRENT_TIMESTAMP", studentID).Scan(&existingOTP.OTPCode, &existingOTP.ExpiresAt)
	if err == nil {
		log.Printf("An unused OTP already exists for student ID %d", studentID)
		return &models.OTPResponse{
			StudentID: studentID,
			OTPCode:   existingOTP.OTPCode,
			ExpiresAt: existingOTP.ExpiresAt,
		}, nil
	} else if errors.Is(err, sql.ErrNoRows) {
		// Mark expired OTPs as used
		_, err = s.db.Exec("UPDATE otps SET is_used = true WHERE student_id = ? AND is_used = false AND expires_at <= CURRENT_TIMESTAMP", studentID)
		if err != nil {
			log.Printf("Error marking expired OTPs as used for student ID %d: %v", studentID, err)
			return nil, fmt.Errorf("failed to update expired OTPs: %w", err)
		}
	} else {
		log.Printf("Database error while checking existing OTP for student ID %d: %v", studentID, err)
		return nil, fmt.Errorf("database error: %w", err)
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
	_, err = s.db.Exec("UPDATE otps SET is_used = true WHERE student_id = ? AND is_used = false", studentID)
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
	err = s.db.QueryRow("INSERT INTO otps (student_id, otp_code, expires_at, is_used) VALUES (?, ?, ?, ?) RETURNING otp_code", studentID, otp, expiresAt, false).Scan(&newOTP.OTPCode)
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
	err := s.db.QueryRow("SELECT student_id, is_used, expires_at FROM otps WHERE otp_code = ?", otpCode).Scan(&otp.StudentID, &otp.IsUsed, &otp.ExpiresAt)
	if errors.Is(err, sql.ErrNoRows) {
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
		_, err := s.db.Exec("UPDATE otps SET is_used = true WHERE otp_code = ?", otpCode)
		if err != nil {
			log.Printf("Error marking expired OTP as used for code %s: %v", otpCode, err) // Improved logging
		}
		return &models.OTPValidationResponse{
			Success: false,
			Message: "OTP has expired",
		}, nil
	}

	// Mark OTP as used
	otp.IsUsed = true
	_, err = s.db.Exec("UPDATE otps SET is_used = true WHERE otp_code = ?", otpCode)
	if err != nil {
		log.Printf("Error marking OTP as used for code %s: %v", otpCode, err) // Improved logging
	}

	// Removed secret code generation and storage

	return &models.OTPValidationResponse{
		Success:   true,
		StudentID: otp.StudentID,
		Message:   "Authentication successful",
	}, nil
}

// VerifyDeviceAuth verifies if a device is authorized using student_id and secret_code
func (s *AuthService) VerifyDeviceAuth(studentID int, secretCode string) (bool, error) {
	var authDevice models.AuthorizedDevice

	// Check if the device exists
	err := s.db.QueryRow("SELECT * FROM authorized_devices WHERE student_id = ? AND secret_code = ?", studentID, secretCode).Scan(&authDevice.ID, &authDevice.StudentID, &authDevice.SecretCode)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	} else if err != nil {
		log.Printf("Database error while verifying device authorization: %v", err)
		return false, fmt.Errorf("database error: %w", err)
	}

	return true, nil
}

// RegisterRoutes registers the routes for AuthService
func (s *AuthService) RegisterRoutes(router *mux.Router) {
	// CORS Setup for AuthService routes
	corsMiddleware := handlers.CORS(
		handlers.AllowedHeaders([]string{
			"Content-Type",
			"Authorization",
			"Origin",
			"Accept",
			"X-Requested-With",
			"Test-Key",
			"testkey",
			"student-id",
			"otp-code",
		}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}),
		handlers.AllowedOrigins([]string{"*"}), // Adjust as needed
		handlers.AllowCredentials(),
		handlers.ExposedHeaders([]string{"Content-Length"}),
		handlers.MaxAge(86400),
	)

	// Apply CORS middleware to AuthService routes
	router.HandleFunc("/generate-otp", s.HandleGenerateOTP).Methods("POST")
	router.HandleFunc("/validate-otp", s.HandleValidateOTP).Methods("POST")
	router.HandleFunc("/verify-device-auth", s.HandleVerifyDeviceAuth).Methods("POST")
	router.Use(corsMiddleware)
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
