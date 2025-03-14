package controllers

import (
	"encoding/json"
	"net/http"
	"server/database"
	"server/models"
	"time"

	"strconv"

	"github.com/gorilla/mux"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/wkt"
)

// Attendance handles the creation and update of attendance records for a student
func Attendance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	studentIDStr := vars["id"]
	studentID, err := strconv.Atoi(studentIDStr)
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	var input struct {
		Action   string `json:"action"`
		Location string `json:"location"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var attendance models.Attendance
	location, err := wkt.Unmarshal(input.Location)
	if err != nil {
		http.Error(w, "Invalid location format", http.StatusBadRequest)
		return
	}

	point := location.(*geom.Point)
	point.SetSRID(4326)

	if input.Action == "checkIn" {
		attendance = models.Attendance{
			StudentID:       studentID,
			CheckInDateTime: time.Now(),
			CheckInLocation: *point,
		}
		if err := database.DB.Create(&attendance).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if input.Action == "checkOut" {
		if err := database.DB.Where("student_id = ? AND DATE(check_in_date_time) = ?", studentID, time.Now().Format("2006-01-02")).First(&attendance).Error; err != nil {
			http.Error(w, "Attendance record not found", http.StatusNotFound)
			return
		}
		attendance.CheckOutDateTime = time.Now()
		attendance.CheckOutLocation = *point
		if err := database.DB.Save(&attendance).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Invalid action", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(attendance)
}
