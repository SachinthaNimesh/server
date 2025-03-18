package controllers

import (
	"encoding/json"
	"net/http"
	"server/models"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// Global DB instance
var db *gorm.DB

func PostAttendance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	studentID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	var requestData struct {
		CheckIn   bool    `json:"check_in"`
		Latitude  float64 `json:"check_in_lat"`
		Longitude float64 `json:"check_in_long"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var attendance models.Attendance
	if requestData.CheckIn {
		attendance.StudentID = studentID
		attendance.CheckInLatitude = requestData.Latitude
		attendance.CheckInLongitude = requestData.Longitude
		attendance.CheckInDateTime = time.Now()

		if err := db.Create(&attendance).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		if err := db.Where("student_id = ? AND DATE(check_in_date_time) = ?", studentID, time.Now().Format("2006-01-02")).First(&attendance).Error; err != nil {
			http.Error(w, "Check-in record not found for today", http.StatusNotFound)
			return
		}

		attendance.CheckOutLatitude = requestData.Latitude
		attendance.CheckOutLongitude = requestData.Longitude
		attendance.CheckOutDateTime = time.Now()

		if err := db.Save(&attendance).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(attendance)
}
