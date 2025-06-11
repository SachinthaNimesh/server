package routes

import (
	"server/controllers"

	"github.com/gorilla/mux"
)

func RegisterStudentRoutes(router *mux.Router) {
	router.HandleFunc("/get-students", controllers.GetStudents).Methods("GET")
	router.HandleFunc("/post-student", controllers.CreateStudent).Methods("POST")

	router.HandleFunc("/get-student", controllers.GetStudent).Methods("GET")

	router.HandleFunc("/delete-student", controllers.DeleteStudent).Methods("DELETE")

	// Add attendance routes
	router.HandleFunc("/attendance", controllers.PostAttendance).Methods("POST")

	// Add mood routes
	router.HandleFunc("/post-mood", controllers.CreateMood).Methods("POST")
	router.HandleFunc("/get-mood", controllers.GetMoods).Methods("GET")

	// Add card routes
	router.HandleFunc("/dashboard", controllers.GetStudentDetails).Methods("GET")

	router.HandleFunc("/employees", controllers.GetEmployeeData).Methods("GET")
	router.HandleFunc("/management", controllers.GetManagementTable).Methods("GET")
	router.HandleFunc("/trainee-profile", controllers.GetTraineeProfile).Methods("GET")
}
