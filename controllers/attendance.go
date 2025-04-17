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

		attendance.CheckInLat = requestData.Latitude
		attendance.CheckInLong = requestData.Longitude
		attendance.CheckInDateTime = time.Now()

		log.Println("Creating new check-in record")
		if err := database.DB.Create(&attendance).Error; err != nil {
			log.Printf("Database error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		if err := database.DB.Where("student_id = ? AND DATE(check_in_date_time) = ?",
			studentID, time.Now().Format("2006-01-02")).First(&attendance).Error; err != nil {
			log.Printf("Check-in record not found: %v", err)
			http.Error(w, "Check-in record not found for today", http.StatusNotFound)
			return
		}

		attendance.CheckOutLat = requestData.Latitude
		attendance.CheckOutLong = requestData.Longitude
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
