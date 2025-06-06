package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"server/database"
	"server/models"

	"gorm.io/gorm"
)

type Attendance struct {
	CheckIn  time.Time  `json:"check_in_date_time"`
	CheckOut *time.Time `json:"check_out_date_time"`
}

type Mood struct {
	Emotion    string    `json:"emotion"`
	RecordedAt time.Time `json:"recorded_at"`
}

type EmployeeSummary struct {
	Attendances []Attendance `json:"attendances"`
	Remarks     string       `json:"remarks"`
	Moods       []Mood       `json:"moods"`
}

func GetEmployeeSummary(w http.ResponseWriter, r *http.Request) {
	idStr := r.Header.Get("student-id")
	if idStr == "" {
		http.Error(w, `{"error":"Missing student-id header"}`, http.StatusBadRequest)
		return
	}
	studentID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error":"Invalid student-id header"}`, http.StatusBadRequest)
		return
	}

	var summary EmployeeSummary

	// 1. Last 5 attendance records (before today)
	var attendances []models.Attendance
	if err := database.DB.
		Where("student_id = ? AND DATE(check_in_date_time) < CURRENT_DATE", studentID).
		Order("check_in_date_time DESC").
		Limit(5).
		Find(&attendances).Error; err != nil && err != gorm.ErrRecordNotFound {
		http.Error(w, `{"error":"Failed to fetch attendance"}`, http.StatusInternalServerError)
		return
	}
	for _, att := range attendances {
		summary.Attendances = append(summary.Attendances, Attendance{
			CheckIn:  att.CheckInDateTime,
			CheckOut: &att.CheckOutDateTime,
		})
	}

	// 2. Remarks
	var student models.Student
	if err := database.DB.Select("remarks").First(&student, studentID).Error; err != nil && err != gorm.ErrRecordNotFound {
		http.Error(w, `{"error":"Failed to fetch remarks"}`, http.StatusInternalServerError)
		return
	}
	summary.Remarks = student.Remarks

	// 3. Last 5 daily mood entries
	var moods []models.Mood
	if err := database.DB.
		Where("student_id = ? AND is_daily = ?", studentID, true).
		Order("recorded_at ASC").
		Limit(5).
		Find(&moods).Error; err != nil && err != gorm.ErrRecordNotFound {
		http.Error(w, `{"error":"Failed to fetch moods"}`, http.StatusInternalServerError)
		return
	}
	for _, m := range moods {
		summary.Moods = append(summary.Moods, Mood{
			Emotion:    m.Emotion,
			RecordedAt: m.RecordedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(summary)
}
