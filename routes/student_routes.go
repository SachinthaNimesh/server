package routes

import (
	"server/controllers"

	"github.com/gorilla/mux"
)

func RegisterStudentRoutes(router *mux.Router) {
	router.HandleFunc("/students", controllers.GetStudents).Methods("GET")
	// router.HandleFunc("/students", controllers.CreateStudent).Methods("POST")
	router.HandleFunc("/students", controllers.GetStudent).Methods("GET")
	router.HandleFunc("/students", controllers.UpdateStudent).Methods("PUT")
	router.HandleFunc("/students", controllers.DeleteStudent).Methods("DELETE")

	// Add attendance routes
	router.HandleFunc("/attendance", controllers.PostAttendance).Methods("POST")

	// Add mood routes
	router.HandleFunc("/moods", controllers.CreateMood).Methods("POST")

	// Add card routes
	router.HandleFunc("/dashboard", controllers.GetStudentDetails).Methods("GET")

	router.HandleFunc("/employees", controllers.GetEmployeeData).Methods("GET")

}
