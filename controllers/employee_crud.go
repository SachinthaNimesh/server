package controllers

import (
	"encoding/json"
	"net/http"
	"server/database"
	"strconv"
	"time"
)

// Employee represents the employee model (formerly Student)
type Employee struct {
	ID                    uint      `json:"id"`
	FirstName             string    `json:"first_name"`
	LastName              string    `json:"last_name"`
	DOB                   time.Time `json:"dob"`
	Gender                string    `json:"gender"`
	AddressLine1          string    `json:"address_line1"`
	AddressLine2          string    `json:"address_line2"`
	City                  string    `json:"city"`
	ContactNumber         string    `json:"contact_number"`
	ContactNumberGuardian string    `json:"contact_number_guardian"`
	SupervisorID          *uint     `json:"supervisor_id"`
	Remarks               string    `json:"remarks"`
	HomeLong              float64   `json:"home_long"`
	HomeLat               float64   `json:"home_lat"`
	EmployerID            *uint     `json:"employer_id"`
	CheckInTime           time.Time `json:"check_in_time"`
	CheckOutTime          time.Time `json:"check_out_time"`
}

// CreateEmployee creates a new employee using the database connection
func CreateEmployee(w http.ResponseWriter, r *http.Request) {
	var employee Employee
	if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := database.DB.Create(&employee).Error; err != nil {
		http.Error(w, "Failed to create employee", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"data": employee})
}

// UpdateEmployee updates an employee record using the database connection
func UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	idStr := r.Header.Get("student-id")
	if idStr == "" {
		http.Error(w, "Missing student-id header", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid student-id header", http.StatusBadRequest)
		return
	}
	var employee Employee
	if err := database.DB.First(&employee, id).Error; err != nil {
		http.Error(w, "Employee not found", http.StatusNotFound)
		return
	}
	var input Employee
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	database.DB.Model(&employee).Updates(input)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"data": employee})
}

// DeleteEmployee removes an employee using the database connection
func DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	idStr := r.Header.Get("student-id")
	if idStr == "" {
		http.Error(w, "Missing student-id header", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid student-id header", http.StatusBadRequest)
		return
	}
	if err := database.DB.Delete(&Employee{}, id).Error; err != nil {
		http.Error(w, "Failed to delete employee", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"data": "Employee deleted successfully"})
}
