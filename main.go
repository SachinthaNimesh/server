// filepath: c:\worky\server\main.go
package main

import (
	"log"
	"net/http"
	"server/database"
	"server/models"
	"server/routes"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	database.ConnectDB()
	database.DB.AutoMigrate(&models.Student{}, &models.Supervisor{}, &models.Mood{}, &models.Employer{}, &models.Attendance{})
	log.Println("Database migrated")

	// Define the server routes
	r := mux.NewRouter()

	// Register the routes
	routes.RegisterStudentRoutes(r)

	// CORS setup
	corsMiddleware := handlers.CORS(
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}),
		handlers.AllowedOrigins([]string{"*"}),
	)

	// Start the server
	log.Println("Server started at http://localhost:8080")
	if err := http.ListenAndServe(":8080", corsMiddleware(r)); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}

/*
curl -X POST http://localhost:8080/students -H "Content-Type: application/json" -d '{"name": "John Doe","age": 21,"email": "john.doe@example.com"}'
curl -X POST http://localhost:8080/moods \-H "Content-Type: application/json" \-d '{    "student_id": 1,    "emotion": "happy",    "is_daily": true}'
curl -X GET http://localhost:8080/students
*/
