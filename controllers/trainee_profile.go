package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"server/database"
	"server/models"
	"strconv"
)

// GetTraineeProfile handles the request to get a trainee's profile information
func GetTraineeProfile(w http.ResponseWriter, r *http.Request) {
	log.Println("Received trainee profile request")

	// Get student ID from header
	studentIDHeader := r.Header.Get("student-id")
	if studentIDHeader == "" {
		log.Println("Missing student-id header")
		http.Error(w, "Missing student-id header", http.StatusBadRequest)
		return
	}

	studentID, err := strconv.Atoi(studentIDHeader)
	if err != nil {
		log.Printf("Invalid student-id header: %v", err)
		http.Error(w, "Invalid student-id header", http.StatusBadRequest)
		return
	}

	log.Printf("Processing trainee profile for student ID: %d", studentID)

	// Fetch student info
	var student models.Student
	if err := database.DB.First(&student, studentID).Error; err != nil {
		log.Printf("Failed to find student: %v", err)
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	// Fetch recent moods
	var recentMoods []models.Mood
	if err := database.DB.Where("student_id = ? AND is_daily = ?", studentID, true).
		Order("recorded_at DESC").Limit(5).Find(&recentMoods).Error; err != nil {
		log.Printf("Error fetching mood data: %v", err)
		// Continue execution even if mood data can't be fetched
	}

	// Fetch recent attendance
	var recentAttendanceRecords []struct {
		ScheduledCheckIn  string `json:"scheduled_check_in"`
		ScheduledCheckOut string `json:"scheduled_check_out"`
		ActualCheckIn     string `json:"actual_check_in"`
		ActualCheckOut    string `json:"actual_check_out"`
	}

	query := `
		SELECT 
			s.check_in_time AS scheduled_check_in,
			s.check_out_time AS scheduled_check_out,
			a.check_in_date_time AS actual_check_in,
			a.check_out_date_time AS actual_check_out
		FROM student s
		JOIN attendance a ON s.id = a.student_id
		WHERE s.id = ?
		ORDER BY a.check_in_date_time DESC
		LIMIT 5`

	if err := database.DB.Raw(query, studentID).Scan(&recentAttendanceRecords).Error; err != nil {
		log.Printf("Error fetching attendance data: %v", err)
		// Continue execution even if attendance data can't be fetched
	}

	// Prepare response
	studentInfo := struct {
		FirstName             string `json:"first_name"`
		LastName              string `json:"last_name"`
		Gender                string `json:"gender"`
		ContactNumber         string `json:"contact_number"`
		ContactNumberGuardian string `json:"contact_number_guardian"`
		Remarks               string `json:"remarks"`
	}{
		FirstName:             student.FirstName,
		LastName:              student.LastName,
		Gender:                student.Gender,
		ContactNumber:         student.ContactNumber,
		ContactNumberGuardian: student.ContactNumberGuardian,
		Remarks:               student.Remarks,
	}

	response := struct {
		StudentInfo      interface{} `json:"student_info"`
		RecentMoods      interface{} `json:"recent_moods"`
		RecentAttendance interface{} `json:"recent_attendance"`
	}{
		StudentInfo:      studentInfo,
		RecentMoods:      recentMoods,
		RecentAttendance: recentAttendanceRecords,
	}

	log.Println("Successfully processed request, sending response")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
