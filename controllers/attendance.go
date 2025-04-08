// There is an issue in model and this controller
package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"server/database" // Make sure to import your database package
	"server/models"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// PostAttendance handles the attendance check-in and check-out
// @Summary Record attendance
// @Description Record check-in or check-out for a student
// @Tags attendance
// @Accept json
// @Produce json
// @Param id path int true "Student ID"
// @Param attendance body struct{CheckIn bool `json:"check_in"`; Latitude float64 `json:"check_in_lat"`; Longitude float64 `json:"check_in_long"`} true "Attendance data"
// @Success 200 {object} models.Attendance
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Check-in record not found for today"
// @Failure 500 {string} string "Internal server error"
// @Router /attendance/{id} [post]
func PostAttendance(w http.ResponseWriter, r *http.Request) {
	log.Println("Received attendance request")

	vars := mux.Vars(r)
	studentID, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Invalid student ID: %v", err)
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	log.Printf("Processing attendance for student ID: %d", studentID)

	var requestData struct {
		CheckIn   bool    `json:"check_in"`
		Latitude  float64 `json:"check_in_lat"`
		Longitude float64 `json:"check_in_long"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Request data: check_in=%v, lat=%f, long=%f",
		requestData.CheckIn, requestData.Latitude, requestData.Longitude)

	var attendance models.Attendance
	if requestData.CheckIn {
		attendance.StudentID = studentID
		attendance.CheckInLatitude = requestData.Latitude
		attendance.CheckInLongitude = requestData.Longitude
		attendance.CheckInDateTime = time.Now()

		log.Println("Creating new check-in record")
		if err := database.DB.Create(&attendance).Error; err != nil {
			log.Printf("Database error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		log.Println("Looking for existing check-in record")
		if err := database.DB.Where("student_id = ? AND DATE(check_in_date_time) = ?",
			studentID, time.Now().Format("2006-01-02")).First(&attendance).Error; err != nil {
			log.Printf("Check-in record not found: %v", err)
			http.Error(w, "Check-in record not found for today", http.StatusNotFound)
			return
		}

		attendance.CheckOutLatitude = requestData.Latitude
		attendance.CheckOutLongitude = requestData.Longitude
		attendance.CheckOutDateTime = time.Now()

		log.Println("Updating existing record with check-out data")
		if err := database.DB.Save(&attendance).Error; err != nil {
			log.Printf("Failed to save record: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	log.Println("Successfully processed request, sending response")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(attendance)
}
