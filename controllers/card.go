package controllers

import (
	"encoding/json"
	"net/http"
	"server/database"
	"server/models"
<<<<<<< HEAD
	"time"
=======
>>>>>>> bd8e16452a5c16d58284408d4bf74ac477abd168
)

func GetStudentDetails(w http.ResponseWriter, r *http.Request) {
	query := `
<<<<<<< HEAD
    SELECT
=======
    SELECT 
>>>>>>> bd8e16452a5c16d58284408d4bf74ac477abd168
        s.id AS student_id,
        s.first_name,
        s.last_name,
        e.name AS employer_name,
        a.check_in_date_time,
        a.check_out_date_time,
        m.emotion
    FROM student s
    LEFT JOIN employer e ON s.id = e.student_id
    LEFT JOIN (
        SELECT a1.*
        FROM attendance a1
        INNER JOIN (
            SELECT student_id, MAX(check_in_date_time) AS latest_check_in
<<<<<<< HEAD
            FROM attendance
=======
            FROM attendance 
>>>>>>> bd8e16452a5c16d58284408d4bf74ac477abd168
            GROUP BY student_id
        ) a2 ON a1.student_id = a2.student_id AND a1.check_in_date_time = a2.latest_check_in
    ) a ON s.id = a.student_id
    LEFT JOIN (
        SELECT m1.*
        FROM mood m1
        INNER JOIN (
            SELECT student_id, MAX(recorded_at) AS latest_update
            FROM mood
            GROUP BY student_id
        ) m2 ON m1.student_id = m2.student_id AND m1.recorded_at = m2.latest_update
    ) m ON s.id = m.student_id;
    `

<<<<<<< HEAD
	var students []models.StudentCard
=======
	var students []models.StudentDetailedResponse
>>>>>>> bd8e16452a5c16d58284408d4bf74ac477abd168
	rows, err := database.DB.Raw(query).Rows() // Use Raw to execute the query
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Failed to execute query"})
		return
	}
	defer rows.Close()

	for rows.Next() {
<<<<<<< HEAD
		var student models.StudentCard
		var checkInDateTime, checkOutDateTime *time.Time
		var emotion *string

=======
		var student models.StudentDetailedResponse
>>>>>>> bd8e16452a5c16d58284408d4bf74ac477abd168
		err := rows.Scan(
			&student.StudentID,
			&student.FirstName,
			&student.LastName,
			&student.EmployerName,
<<<<<<< HEAD
			&checkInDateTime,
			&checkOutDateTime,
			&emotion,
=======
			&student.CheckInDateTime,
			&student.CheckOutDateTime,
			&student.Emotion,
>>>>>>> bd8e16452a5c16d58284408d4bf74ac477abd168
		)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Failed to scan row"})
			return
		}
<<<<<<< HEAD

		// Handle NULL values
		if checkInDateTime != nil {
			student.CheckInDateTime = *checkInDateTime
		}
		if checkOutDateTime != nil {
			student.CheckOutDateTime = *checkOutDateTime
		}
		if emotion != nil {
			student.Emotion = *emotion
		}

=======
>>>>>>> bd8e16452a5c16d58284408d4bf74ac477abd168
		students = append(students, student)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}
