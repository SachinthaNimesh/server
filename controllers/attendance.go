package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"server/database"
	"server/models"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

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

	// Restructured request data to match expected input
	var requestData struct {
		CheckIn   bool    `json:"check_in"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Request data: check_in=%v, lat=%f, long=%f",
		requestData.CheckIn, requestData.Latitude, requestData.Longitude)

	db := database.DB

	// Using explicit field assignment to avoid any name resolution issues
	if requestData.CheckIn {
		// Create a new check-in record
		attendance := models.Attendance{
			StudentID:        studentID,
			CheckInDateTime:  time.Now(),
			CheckInLatitude:  requestData.Latitude,
			CheckInLongitude: requestData.Longitude,
		}

		log.Println("Creating new check-in record")
		if err := db.Create(&attendance).Error; err != nil {
			log.Printf("Database error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(attendance)
	} else {
		// Find the existing check-in record for today
		var attendance models.Attendance
		log.Println("Looking for existing check-in record")
		if err := db.Where("student_id = ? AND DATE(check_in_date_time) = ?",
			studentID, time.Now().Format("2006-01-02")).First(&attendance).Error; err != nil {
			log.Printf("Check-in record not found: %v", err)
			http.Error(w, "Check-in record not found for today", http.StatusNotFound)
			return
		}

		// Update check-out information
		attendance.CheckOutDateTime = time.Now()
		attendance.CheckOutLatitude = requestData.Latitude
		attendance.CheckOutLongitude = requestData.Longitude

		log.Println("Updating existing record with check-out data")
		if err := db.Save(&attendance).Error; err != nil {
			log.Printf("Failed to save record: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(attendance)
	}

	log.Println("Successfully processed request")
}
