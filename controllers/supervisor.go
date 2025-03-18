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
// @Description Get details of all supervisors
// @Tags supervisors
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Supervisor
// @Router /supervisors [get]
func GetSupervisors(w http.ResponseWriter, r *http.Request) {
	var supervisors []models.Supervisor
	if err := database.DB.Find(&supervisors).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(supervisors)
}

// GetSupervisor godoc
// @Summary Get a supervisor by ID
// @Description Get details of a supervisor by ID
// @Tags supervisors
// @Accept  json
// @Produce  json
// @Param id path int true "Supervisor ID"
// @Success 200 {object} models.Supervisor
// @Router /supervisors/{id} [get]
func GetSupervisor(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid supervisor ID", http.StatusBadRequest)
		return
	}

	var supervisor models.Supervisor
	if err := database.DB.First(&supervisor, id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(supervisor)
}

// CreateSupervisor godoc
// @Summary Create a new supervisor
// @Description Create a new supervisor with the input payload
// @Tags supervisors
// @Accept  json
// @Produce  json
// @Param supervisor body models.Supervisor true "Supervisor"
// @Success 200 {object} models.Supervisor
// @Router /supervisors [post]
func CreateSupervisor(w http.ResponseWriter, r *http.Request) {
	var supervisor models.Supervisor
	if err := json.NewDecoder(r.Body).Decode(&supervisor); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := database.DB.Create(&supervisor).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(supervisor)
}

// UpdateSupervisor godoc
// @Summary Update a supervisor by ID
// @Description Update details of a supervisor by ID
// @Tags supervisors
// @Accept  json
// @Produce  json
// @Param id path int true "Supervisor ID"
// @Param supervisor body models.Supervisor true "Supervisor"
// @Success 200 {object} models.Supervisor
// @Router /supervisors/{id} [put]
func UpdateSupervisor(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid supervisor ID", http.StatusBadRequest)
		return
	}

	var supervisor models.Supervisor
	if err := database.DB.First(&supervisor, id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&supervisor); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := database.DB.Save(&supervisor).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(supervisor)
}

// DeleteSupervisor godoc
// @Summary Delete a supervisor by ID
// @Description Delete a supervisor by ID
// @Tags supervisors
// @Accept  json
// @Produce  json
// @Param id path int true "Supervisor ID"
// @Success 204
// @Router /supervisors/{id} [delete]
func DeleteSupervisor(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid supervisor ID", http.StatusBadRequest)
		return
	}

	if err := database.DB.Delete(&models.Supervisor{}, id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
