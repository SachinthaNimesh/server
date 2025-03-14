// filepath: c:\worky\server\main.go
package main

import (
	"log"
	"net/http"
	"server/controllers"
	"server/database"
	"server/models"

	"github.com/gorilla/mux"
)

func main() {
	database.ConnectDB()
	database.DB.AutoMigrate(&models.Student{}, &models.Supervisor{}, &models.Mood{}, &models.Employer{}, &models.Attendance{})
	log.Println("Database migrated")

	// Define the server routes
	r := mux.NewRouter()

	// Register the routes
	r.HandleFunc("/students", controllers.GetStudents).Methods("GET")
	r.HandleFunc("/students", controllers.CreateStudent).Methods("POST")
	r.HandleFunc("/students/{id}", controllers.GetStudent).Methods("GET")
	r.HandleFunc("/students/{id}", controllers.UpdateStudent).Methods("PUT")
	r.HandleFunc("/students/{id}", controllers.DeleteStudent).Methods("DELETE")
	r.HandleFunc("/attendance/{id}", controllers.Attendance).Methods("POST")
	//mark mood <if checkout button is true then mark mood.isDaily = true>
	//r.HandleFunc("/students/{id}/mood", controllers.mood).Methods("POST")

	// Start the server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}

/*
curl -X POST http://localhost:8080/students -H "Content-Type: application/json" -d '{"name": "John Doe","age": 21,"email": "john.doe@example.com"}'
*/
