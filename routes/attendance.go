package routes

import (
	"server/controllers"

	"github.com/gorilla/mux"
)

// SetupAttendanceRoutes sets up the routes for attendance management
func SetupAttendanceRoutes(r *mux.Router) {
	r.HandleFunc("/attendance/{id}", controllers.Attendance).Methods("POST")
}
