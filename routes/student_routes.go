package routes

import (
	"server/controllers"

	"github.com/gorilla/mux"
)

func RegisterStudentRoutes(router *mux.Router) {
	router.HandleFunc("/students", controllers.GetStudents).Methods("GET")
	// router.HandleFunc("/students", controllers.CreateStudent).Methods("POST")
	router.HandleFunc("/students/{id}", controllers.GetStudent).Methods("GET")
	router.HandleFunc("/students/{id}", controllers.UpdateStudent).Methods("PUT")
	router.HandleFunc("/students/{id}", controllers.DeleteStudent).Methods("DELETE")

	// Add attendance routes
	router.HandleFunc("/attendance", controllers.PostAttendance).Methods("POST")

	// Add mood routes
	router.HandleFunc("/moods/{id}", controllers.GetMood).Methods("GET")
	router.HandleFunc("/moods", controllers.CreateMood).Methods("POST")

	// Add card routes
	router.HandleFunc("/dashboard", controllers.GetStudentDetails).Methods("GET")

	// RegisterRoutes sets up the routes for the AuthService

	// Initialize AuthService
	authService := controllers.NewAuthService()

	// Register AuthService routes
	authService.RegisterRoutes(router)

	router.HandleFunc("/employees", controllers.GetEmployeeData).Methods("GET")

}
