package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"server/database"
	"server/models"

	"github.com/gorilla/mux"
)

// GetSupervisors godoc
// @Summary Get all supervisors
// @Description Get all supervisors
// @Tags supervisors
// @Produce json
// @Success 200 {array} models.Supervisor
// @Failure 500 {string} string "Internal Server Error"
// @Router /supervisors [get]
func GetSupervisors(w http.ResponseWriter, r *http.Request) {
	var supervisors []models.Supervisor
	rows, err := database.DB.Query("SELECT supervisor_id, student_id, first_name, last_name, email_address, contact_number FROM supervisor")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var s models.Supervisor
		if err := rows.Scan(&s.SupervisorID, &s.StudentID, &s.FirstName, &s.LastName, &s.EmailAddress, &s.ContactNumber); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		supervisors = append(supervisors, s)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(supervisors)
}

// GetSupervisor godoc
// @Summary Get a supervisor by ID
// @Description Get a supervisor by ID
// @Tags supervisors
// @Produce json
// @Param id path int true "Supervisor ID"
// @Success 200 {object} models.Supervisor
// @Failure 400 {string} string "Invalid supervisor ID"
// @Failure 404 {string} string "Supervisor not found"
// @Router /supervisors/{id} [get]
func GetSupervisor(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid supervisor ID", http.StatusBadRequest)
		return
	}
	var s models.Supervisor
	err = database.DB.QueryRow("SELECT supervisor_id, student_id, first_name, last_name, email_address, contact_number FROM supervisor WHERE supervisor_id = $1", id).Scan(&s.SupervisorID, &s.StudentID, &s.FirstName, &s.LastName, &s.EmailAddress, &s.ContactNumber)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

// CreateSupervisor godoc
// @Summary Create a new supervisor
// @Description Create a new supervisor
// @Tags supervisors
// @Accept json
// @Produce json
// @Param supervisor body models.Supervisor true "Supervisor"
// @Success 201 {object} models.Supervisor
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /supervisors [post]
func CreateSupervisor(w http.ResponseWriter, r *http.Request) {
	var s models.Supervisor
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	query := `INSERT INTO supervisor (student_id, first_name, last_name, email_address, contact_number) VALUES ($1, $2, $3, $4, $5) RETURNING supervisor_id`
	err := database.DB.QueryRow(query, s.StudentID, s.FirstName, s.LastName, s.EmailAddress, s.ContactNumber).Scan(&s.SupervisorID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

// UpdateSupervisor godoc
// @Summary Update a supervisor by ID
// @Description Update a supervisor by ID
// @Tags supervisors
// @Accept json
// @Produce json
// @Param id path int true "Supervisor ID"
// @Param supervisor body models.Supervisor true "Supervisor"
// @Success 200 {object} models.Supervisor
// @Failure 400 {string} string "Invalid supervisor ID"
// @Failure 404 {string} string "Supervisor not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /supervisors/{id} [put]
func UpdateSupervisor(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid supervisor ID", http.StatusBadRequest)
		return
	}
	var s models.Supervisor
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	query := `UPDATE supervisor SET student_id=$1, first_name=$2, last_name=$3, email_address=$4, contact_number=$5 WHERE supervisor_id=$6`
	_, err = database.DB.Exec(query, s.StudentID, s.FirstName, s.LastName, s.EmailAddress, s.ContactNumber, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.SupervisorID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

// DeleteSupervisor godoc
// @Summary Delete a supervisor by ID
// @Description Delete a supervisor by ID
// @Tags supervisors
// @Param id path int true "Supervisor ID"
// @Success 204 {string} string "No Content"
// @Failure 400 {string} string "Invalid supervisor ID"
// @Failure 500 {string} string "Internal Server Error"
// @Router /supervisors/{id} [delete]
func DeleteSupervisor(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid supervisor ID", http.StatusBadRequest)
		return
	}
	_, err = database.DB.Exec("DELETE FROM supervisor WHERE supervisor_id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GetAllSupervisorIDsAndNames godoc
// @Summary Get all supervisor IDs and names
// @Description Returns a list of all supervisor IDs and their names
// @Tags supervisors
// @Produce json
// @Success 200 {array} object
// @Router /supervisors/ids-names [get]
func GetAllSupervisorIDsAndNames(w http.ResponseWriter, r *http.Request) {
	type SupervisorIDName struct {
		SupervisorID uint64 `json:"supervisor_id"`
		FirstName    string `json:"first_name"`
		LastName     string `json:"last_name"`
	}
	rows, err := database.DB.Query("SELECT supervisor_id, first_name, last_name FROM supervisor")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var supervisors []SupervisorIDName
	for rows.Next() {
		var s SupervisorIDName
		if err := rows.Scan(&s.SupervisorID, &s.FirstName, &s.LastName); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		supervisors = append(supervisors, s)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(supervisors)
}
