package controllers

import (
	"encoding/json"
	"net/http"
	"server/database"
	"server/models"
)

func GetStudentDetails(w http.ResponseWriter, r *http.Request) {
	query := `
    SELECT 
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
            FROM attendance 
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

	var students []models.StudentDetailedResponse
	rows, err := database.DB.Raw(query).Rows() // Use Raw to execute the query
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Failed to execute query"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var student models.StudentDetailedResponse
		err := rows.Scan(
			&student.StudentID,
			&student.FirstName,
			&student.LastName,
			&student.EmployerName,
			&student.CheckInDateTime,
			&student.CheckOutDateTime,
			&student.Emotion,
		)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Failed to scan row"})
			return
		}
		students = append(students, student)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}
