package controllers

//implement post method for mood saving
import (
	"encoding/json"
	"net/http"
	"server/database"
	"server/models"
	"time"

	"github.com/gorilla/mux"
)

// GetMoods godoc
// @Summary Get all moods
// @Description Get details of all moods
// @Tags moods
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Mood
// @Router /moods [get]
func GetMoods(w http.ResponseWriter, r *http.Request) {
	var moods []models.Mood
	if err := database.DB.Find(&moods).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(moods)
}

// GetMood godoc
// @Summary Get a mood by ID
// @Description Get details of a mood by ID
// @Tags moods
// @Accept  json
// @Produce  json
// @Param id path int true "Mood ID"
// @Success 200 {object} models.Mood
// @Router /moods/{id} [get]
func GetMood(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var mood models.Mood
	if err := database.DB.First(&mood, id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mood)
}

// CreateMood godoc
// @Summary Create a new mood
// @Description Create a new mood with the input payload
// @Tags moods
// @Accept  json
// @Produce  json
// @Param mood body models.Mood true "Mood"
// @Success 200 {object} models.Mood
// @Router /moods [post]
func CreateMood(w http.ResponseWriter, r *http.Request) {
	var mood models.Mood
	if err := json.NewDecoder(r.Body).Decode(&mood); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mood.RecordedAt = time.Now()
	if err := database.DB.Create(&mood).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mood)
}

// UpdateMood godoc
// @Summary Update a mood by ID
// @Description Update details of a mood by ID
// @Tags moods
// @Accept  json
// @Produce  json
// @Param id path int true "Mood ID"
// @Param mood body models.Mood true "Mood"
// @Success 200 {object} models.Mood
// @Router /moods/{id} [put]
func UpdateMood(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var mood models.Mood
	if err := database.DB.First(&mood, id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&mood); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := database.DB.Save(&mood).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(mood)
}

// DeleteMood godoc
// @Summary Delete a mood by ID
// @Description Delete a mood by ID
// @Tags moods
// @Accept  json
// @Produce  json
// @Param id path int true "Mood ID"
// @Success 204
// @Router /moods/{id} [delete]
func DeleteMood(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := database.DB.Delete(&models.Mood{}, id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
