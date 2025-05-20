package routes

import (
	"server/controllers"

	"github.com/gorilla/mux"
)

func RegisterStudentRoutes(router *mux.Router) {
	router.HandleFunc("/students", controllers.GetStudents).Methods("GET")
	// router.HandleFunc("/students", controllers.CreateStudent).Methods("POST")
	router.HandleFunc("/student", controllers.GetStudent).Methods("GET")
	router.HandleFunc("/student", controllers.UpdateStudent).Methods("PUT")
	router.HandleFunc("/student", controllers.DeleteStudent).Methods("DELETE")

	// Add attendance routes
	router.HandleFunc("/attendance", controllers.PostAttendance).Methods("POST")

	// Add mood routes
	router.HandleFunc("/mood", controllers.CreateMood).Methods("POST")
	router.HandleFunc("/mood", controllers.GetMoods).Methods("GET")

	// Add card routes
	router.HandleFunc("/dashboard", controllers.GetStudentDetails).Methods("GET")

	router.HandleFunc("/employees", controllers.GetEmployeeData).Methods("GET")

}
