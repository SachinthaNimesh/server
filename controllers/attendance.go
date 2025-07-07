package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"server/database"
	"server/models"
	"strconv"
	"time"
)

func PostAttendance(w http.ResponseWriter, r *http.Request) {
	log.Println("Received attendance request")

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
		// Delete any existing records for today first
		now := time.Now().UTC()
		startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		endOfDay := startOfDay.Add(24 * time.Hour)

		deleteQuery := `DELETE FROM attendance WHERE student_id = $1 AND check_in_date_time >= $2 AND check_in_date_time < $3`
		log.Printf("Delete Query: %s", deleteQuery)
		log.Printf("Delete Params: student_id=%d, start=%v, end=%v", studentID, startOfDay, endOfDay)
		_, err := database.DB.Exec(deleteQuery, studentID, startOfDay, endOfDay)
		if err != nil {
			log.Printf("Failed to delete existing records: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		attendance.StudentID = studentID
		attendance.CheckInLat = requestData.Latitude
		attendance.CheckInLong = requestData.Longitude
		attendance.CheckInDateTime = time.Now()

		log.Println("Creating new check-in record")
		query := `INSERT INTO attendance (student_id, check_in_lat, check_in_long, check_in_date_time) VALUES ($1, $2, $3, $4) RETURNING id, student_id, check_in_lat, check_in_long, check_in_date_time`
		log.Printf("Insert Query: %s", query)
		log.Printf("Insert Params: student_id=%d, lat=%f, long=%f, datetime=%v", attendance.StudentID, attendance.CheckInLat, attendance.CheckInLong, attendance.CheckInDateTime)
		row := database.DB.QueryRow(query, attendance.StudentID, attendance.CheckInLat, attendance.CheckInLong, attendance.CheckInDateTime)
		err = row.Scan(&attendance.ID, &attendance.StudentID, &attendance.CheckInLat, &attendance.CheckInLong, &attendance.CheckInDateTime)
		if err != nil {
			log.Printf("Database error on insert: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Printf("Check-in record created: %+v", attendance)
	} else {
		now := time.Now().UTC()
		startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		endOfDay := startOfDay.Add(24 * time.Hour)

		// Delete any existing records for today first (to ensure only latest record)
		deleteQuery := `DELETE FROM attendance WHERE student_id = $1 AND (check_in_date_time >= $2 AND check_in_date_time < $3) OR (check_out_date_time >= $2 AND check_out_date_time < $3)`
		log.Printf("Delete Query: %s", deleteQuery)
		log.Printf("Delete Params: student_id=%d, start=%v, end=%v", studentID, startOfDay, endOfDay)
		_, err := database.DB.Exec(deleteQuery, studentID, startOfDay, endOfDay)
		if err != nil {
			log.Printf("Failed to delete existing records: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Create a new record with zero check-in values and actual checkout data
		log.Println("Creating checkout record with zero check-in values")
		attendance.StudentID = studentID
		attendance.CheckInLat = 0
		attendance.CheckInLong = 0
		attendance.CheckInDateTime = time.Time{} // zero timestamp
		attendance.CheckOutLat = sql.NullFloat64{Float64: requestData.Latitude, Valid: true}
		attendance.CheckOutLong = sql.NullFloat64{Float64: requestData.Longitude, Valid: true}
		attendance.CheckOutDateTime = sql.NullTime{Time: time.Now(), Valid: true}

		insertQuery := `INSERT INTO attendance (student_id, check_in_lat, check_in_long, check_in_date_time, check_out_lat, check_out_long, check_out_date_time) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
		log.Printf("Insert Query: %s", insertQuery)
		log.Printf("Insert Params: student_id=%d, check_in_lat=0, check_in_long=0, check_in_time=zero, check_out_lat=%v, check_out_long=%v, check_out_time=%v",
			attendance.StudentID, attendance.CheckOutLat, attendance.CheckOutLong, attendance.CheckOutDateTime)
		row := database.DB.QueryRow(insertQuery, attendance.StudentID, attendance.CheckInLat, attendance.CheckInLong, attendance.CheckInDateTime, attendance.CheckOutLat, attendance.CheckOutLong, attendance.CheckOutDateTime)
		err = row.Scan(&attendance.ID)
		if err != nil {
			log.Printf("Database error on checkout insert: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Printf("Checkout record with zero check-in values created: %+v", attendance)
	}

	log.Println("Successfully processed request, sending response")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(attendance)
}
