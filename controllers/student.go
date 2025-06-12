package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"server/database"
	"server/models"
	"strconv"
)

func getStudentIDFromHeader(r *http.Request) (int, error) {
	StudentIDHeader := r.Header.Get("student-id")
	if StudentIDHeader == "" {
		return 0, http.ErrMissingFile
	}
	studentID, err := strconv.Atoi(StudentIDHeader)
	if err != nil {
		return 0, err
	}
	return studentID, nil
}

// GetStudents godoc
// @Summary Get all students
// @Description Get all students
// @Tags students
// @Produce json
// @Success 200 {array} models.Student
// @Failure 500 {string} string "Internal Server Error"
// @Router /students [get]
func GetStudents(w http.ResponseWriter, r *http.Request) {
	var students []models.Student
	if err := database.DB.Find(&students).Error; err != nil {
		log.Printf("Error fetching students: %v", err) // Log error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Fetched %d students", len(students)) // Log success
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}

// GetStudent godoc
// @Summary Get a student by ID
// @Description Get a student by ID
// @Tags students
// @Produce json
// @Param student-id header string true "Student ID"
// @Success 200 {object} models.Student
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Router /get-student [get]
func GetStudent(w http.ResponseWriter, r *http.Request) {
	studentID, err := getStudentIDFromHeader(r)
	if err != nil {
		log.Printf("Error extracting student-id: %v", err)
		http.Error(w, "Invalid or missing student-id header", http.StatusBadRequest)
		return
	}

	var student models.Student
	if err := database.DB.First(&student, studentID).Error; err != nil {
		log.Printf("Error fetching student with ID %d: %v", studentID, err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	log.Printf("Fetched student with ID %d: %+v", studentID, student) // Log success
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)
}

// CreateStudent godoc
// @Summary Create a new student
// @Description Create a new student
// @Tags students
// @Accept json
// @Produce json
// @Param student body models.Student true "Student"
// @Success 201 {object} models.Student
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /students [post]
func CreateStudent(w http.ResponseWriter, r *http.Request) {
	var student models.Student
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := database.DB.Create(&student).Error; err != nil {
		http.Error(w, "Failed to create student", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"data": student})
}

// UpdateStudent godoc
// @Summary Update a student by ID
// @Description Update a student by ID
// @Tags students
// @Accept json
// @Produce json
// @Param id path string true "Student ID"
// @Param student body models.Student true "Student"
// @Success 200 {object} models.Student
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /students/{id} [put]
func UpdateStudent(w http.ResponseWriter, r *http.Request) {
	idStr := r.Header.Get("student-id")
	if idStr == "" {
		http.Error(w, "Missing student-id header", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid student-id header", http.StatusBadRequest)
		return
	}
	var student models.Student
	if err := database.DB.First(&student, id).Error; err != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}
	var input models.Student
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	database.DB.Model(&student).Updates(input)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"data": student})
}

// DeleteStudent godoc
// @Summary Delete a student by ID
// @Description Delete a student by ID
// @Tags students
// @Param id path string true "Student ID"
// @Success 204 {string} string "No Content"
// @Failure 500 {string} string "Internal Server Error"
// @Router /students/{id} [delete]
func DeleteStudent(w http.ResponseWriter, r *http.Request) {
	idStr := r.Header.Get("student-id")
	if idStr == "" {
		http.Error(w, "Missing student-id header", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid student-id header", http.StatusBadRequest)
		return
	}
	if err := database.DB.Delete(&models.Student{}, id).Error; err != nil {
		http.Error(w, "Failed to delete student", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"data": "Student deleted successfully"})
}
