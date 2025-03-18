package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"server/database"
	"server/models"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

func CreateEmployer(w http.ResponseWriter, r *http.Request) {
	var employerInput struct {
		Name          string  `json:"name"`
		StudentID     int     `json:"student_id"`
		ContactNumber string  `json:"contact_number"`
		AddressLine1  string  `json:"address_line_1"`
		AddressLine2  string  `json:"address_line_2"`
		AddressLine3  string  `json:"address_line_3"`
		Longitude     float64 `json:"addr_long"`
		Latitude      float64 `json:"addr_lat"`
	}
	if err := json.NewDecoder(r.Body).Decode(&employerInput); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	employer := models.Employer{
		Name:          employerInput.Name,
		StudentID:     employerInput.StudentID,
		ContactNumber: employerInput.ContactNumber,
		AddressLine1:  employerInput.AddressLine1,
		AddressLine2:  employerInput.AddressLine2,
		AddressLine3:  employerInput.AddressLine3,
		Longitude:     employerInput.Longitude,
		Latitude:      employerInput.Latitude,
	}
	if err := database.DB.Create(&employer).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(employer)
}

func GetEmployer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var employer models.Employer
	if err := database.DB.First(&employer, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			http.Error(w, "Employer not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(employer)
}

func UpdateEmployer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var employer models.Employer
	if err := database.DB.First(&employer, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			http.Error(w, "Employer not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	var employerInput struct {
		Name          string  `json:"name"`
		StudentID     int     `json:"student_id"`
		ContactNumber string  `json:"contact_number"`
		AddressLine1  string  `json:"address_line1"`
		AddressLine2  string  `json:"address_line2"`
		AddressLine3  string  `json:"address_line3"`
		Longitude     float64 `json:"addr_long"`
		Latitude      float64 `json:"addr_lat"`
	}
	if err := json.NewDecoder(r.Body).Decode(&employerInput); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	employer.Name = employerInput.Name
	employer.StudentID = employerInput.StudentID
	employer.ContactNumber = employerInput.ContactNumber
	employer.AddressLine1 = employerInput.AddressLine1
	employer.AddressLine2 = employerInput.AddressLine2
	employer.AddressLine3 = employerInput.AddressLine3
	employer.Longitude = employerInput.Longitude
	employer.Latitude = employerInput.Latitude

	if err := database.DB.Save(&employer).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(employer)
}

func DeleteEmployer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := database.DB.Delete(&models.Employer{}, id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
