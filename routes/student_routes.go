package routes

import (
	"server/controllers"

	"github.com/gorilla/mux"
)

func RegisterStudentRoutes(router *mux.Router) {
	router.HandleFunc("/students", controllers.GetStudents).Methods("GET")
	router.HandleFunc("/students", controllers.CreateStudent).Methods("POST")
	router.HandleFunc("/students/{id}", controllers.GetStudent).Methods("GET")
	router.HandleFunc("/students/{id}", controllers.UpdateStudent).Methods("PUT")
	router.HandleFunc("/students/{id}", controllers.DeleteStudent).Methods("DELETE")

	// Add attendance routes
	router.HandleFunc("/attendance/{id}", controllers.PostAttendance).Methods("POST")

	// Add mood routes
	router.HandleFunc("/moods/{id}", controllers.GetMood).Methods("GET")
	router.HandleFunc("/moods", controllers.CreateMood).Methods("POST")

	// Add employer routes
	router.HandleFunc("/employers", controllers.CreateEmployer).Methods("POST")
	router.HandleFunc("/employers/{id}", controllers.GetEmployer).Methods("GET")
	router.HandleFunc("/employers/{id}", controllers.UpdateEmployer).Methods("PUT")
	router.HandleFunc("/employers/{id}", controllers.DeleteEmployer).Methods("DELETE")

	// Add supervisor routes
	router.HandleFunc("/supervisors", controllers.GetSupervisors).Methods("GET")
	router.HandleFunc("/supervisors/{id}", controllers.GetSupervisor).Methods("GET")
	router.HandleFunc("/supervisors", controllers.CreateSupervisor).Methods("POST")
	router.HandleFunc("/supervisors/{id}", controllers.UpdateSupervisor).Methods("PUT")
	router.HandleFunc("/supervisors/{id}", controllers.DeleteSupervisor).Methods("DELETE")

	// Add student details route
	router.HandleFunc("/students/details", controllers.GetStudentDetails).Methods("GET")
	// router.HandleFunc("/students/{id}/employer", controllers.GetEmployerByStudentID).Methods("GET")
}
