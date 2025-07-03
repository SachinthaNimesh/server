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
	now := time.Now().UTC()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)

	// Always look for today's record first
	query := `SELECT id, student_id, check_in_lat, check_in_long, check_in_date_time, check_out_lat, check_out_long, check_out_date_time 
			  FROM attendance 
			  WHERE student_id = $1 AND check_in_date_time >= $2 AND check_in_date_time < $3 LIMIT 1`
	row := database.DB.QueryRow(query, studentID, startOfDay, endOfDay)
	err = row.Scan(&attendance.ID, &attendance.StudentID, &attendance.CheckInLat, &attendance.CheckInLong, &attendance.CheckInDateTime, &attendance.CheckOutLat, &attendance.CheckOutLong, &attendance.CheckOutDateTime)
	recordExists := err != sql.ErrNoRows && err == nil

	if requestData.CheckIn {
		// If record exists, update check-in fields with latest values
		if recordExists {
			attendance.CheckInLat = requestData.Latitude
			attendance.CheckInLong = requestData.Longitude
			attendance.CheckInDateTime = now

			updateQuery := `UPDATE attendance SET check_in_lat = $1, check_in_long = $2, check_in_date_time = $3 WHERE id = $4`
			_, err = database.DB.Exec(updateQuery, attendance.CheckInLat, attendance.CheckInLong, attendance.CheckInDateTime, attendance.ID)
			if err != nil {
				log.Printf("Failed to update check-in: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// Fetch updated record for response
			row = database.DB.QueryRow(query, studentID, startOfDay, endOfDay)
			_ = row.Scan(&attendance.ID, &attendance.StudentID, &attendance.CheckInLat, &attendance.CheckInLong, &attendance.CheckInDateTime, &attendance.CheckOutLat, &attendance.CheckOutLong, &attendance.CheckOutDateTime)
		} else {
			attendance.StudentID = studentID
			attendance.CheckInLat = requestData.Latitude
			attendance.CheckInLong = requestData.Longitude
			attendance.CheckInDateTime = now

			insertQuery := `INSERT INTO attendance (student_id, check_in_lat, check_in_long, check_in_date_time) VALUES ($1, $2, $3, $4) RETURNING id, student_id, check_in_lat, check_in_long, check_in_date_time, check_out_lat, check_out_long, check_out_date_time`
			row = database.DB.QueryRow(insertQuery, attendance.StudentID, attendance.CheckInLat, attendance.CheckInLong, attendance.CheckInDateTime)
			err = row.Scan(&attendance.ID, &attendance.StudentID, &attendance.CheckInLat, &attendance.CheckInLong, &attendance.CheckInDateTime, &attendance.CheckOutLat, &attendance.CheckOutLong, &attendance.CheckOutDateTime)
			if err != nil {
				log.Printf("Database error on insert: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	} else {
		// Must check out only if a record exists for today
		if !recordExists {
			// Try to find the latest previous check-in record without a check-out
			previousQuery := `SELECT id, student_id, check_in_lat, check_in_long, check_in_date_time, check_out_lat, check_out_long, check_out_date_time
							  FROM attendance
							  WHERE student_id = $1 AND check_out_date_time IS NULL
							  ORDER BY check_in_date_time DESC LIMIT 1`
			row = database.DB.QueryRow(previousQuery, studentID)
			err = row.Scan(&attendance.ID, &attendance.StudentID, &attendance.CheckInLat, &attendance.CheckInLong, &attendance.CheckInDateTime, &attendance.CheckOutLat, &attendance.CheckOutLong, &attendance.CheckOutDateTime)
			if err != nil {
				log.Printf("No open check-in record found for checkout: %v", err)
				http.Error(w, "No open check-in record found for checkout", http.StatusNotFound)
				return
			}
		}
		attendance.CheckOutLat = sql.NullFloat64{Float64: requestData.Latitude, Valid: true}
		attendance.CheckOutLong = sql.NullFloat64{Float64: requestData.Longitude, Valid: true}
		attendance.CheckOutDateTime = sql.NullTime{Time: now, Valid: true}

		updateQuery := `UPDATE attendance SET check_out_lat = $1, check_out_long = $2, check_out_date_time = $3 WHERE id = $4`
		_, err = database.DB.Exec(updateQuery, attendance.CheckOutLat, attendance.CheckOutLong, attendance.CheckOutDateTime, attendance.ID)
		if err != nil {
			log.Printf("Failed to update check-out: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Fetch updated record for response
		row = database.DB.QueryRow(query, studentID, startOfDay, endOfDay)
		_ = row.Scan(&attendance.ID, &attendance.StudentID, &attendance.CheckInLat, &attendance.CheckInLong, &attendance.CheckInDateTime, &attendance.CheckOutLat, &attendance.CheckOutLong, &attendance.CheckOutDateTime)
	}

	log.Println("Successfully processed request, sending response")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(attendance)
}
