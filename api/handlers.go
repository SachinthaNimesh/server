package api

import (
	"encoding/json"
	"net/http"

	"server/controllers"
	"server/models"
)

// APIHandler handles HTTP requests
type APIHandler struct {
	authService *controllers.AuthService
}

// NewAPIHandler creates a new API handler
func NewAPIHandler() *APIHandler {
	return &APIHandler{
		authService: controllers.NewAuthService(),
	}
}

// GenerateOTP handles OTP generation requests
func (h *APIHandler) GenerateOTP(w http.ResponseWriter, r *http.Request) {
	var req models.OTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.StudentID <= 0 {
		respondError(w, http.StatusBadRequest, "Invalid employer ID")
		return
	}

	resp, err := h.authService.GenerateOTP(req.StudentID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

// ValidateOTP handles OTP validation requests
func (h *APIHandler) ValidateOTP(w http.ResponseWriter, r *http.Request) {
	var req models.OTPValidationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.OTPCode == "" {
		respondError(w, http.StatusBadRequest, "OTP code is required")
		return
	}

	resp, err := h.authService.ValidateOTP(req.OTPCode)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !resp.Success {
		respondJSON(w, http.StatusUnauthorized, resp)
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

// VerifyAuth handles authentication verification
func (h *APIHandler) VerifyAuth(w http.ResponseWriter, r *http.Request) {
	var req models.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	isValid, err := h.authService.VerifyDeviceAuth(req.StudentID, req.SecretCode)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !isValid {
		respondError(w, http.StatusUnauthorized, "Invalid authentication credentials")
		return
	}

	respondJSON(w, http.StatusOK, map[string]bool{"authenticated": true})
}

// Helper function to respond with JSON
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

// Helper function to respond with error
func respondError(w http.ResponseWriter, code int, message string) {
	respondJSON(w, code, map[string]string{"error": message})
}
