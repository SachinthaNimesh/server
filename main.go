package main

import (
	"log"
	"net/http"
	"os"
	"server/database"
	"server/models"
	"server/routes"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	// Fetch port from environment (default to 8000 for Choreo)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Connect to DB with environment variables
	database.ConnectDB()
	database.DB.AutoMigrate(&models.Student{}, &models.Supervisor{}, &models.Mood{}, &models.Employer{}, &models.Attendance{})
	log.Println("Database migrated")

	// Define router
	r := mux.NewRouter()

	// Register API routes
	routes.RegisterStudentRoutes(r)

	// CORS Setup - Set allowed origins to Choreo API Gateway
	corsMiddleware := handlers.CORS(
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}),
		handlers.AllowedOrigins([]string{"*"}), // Replace "*" with Choreo domain in production
	)

	// Start the server
	log.Println("Server started on port", port)
	if err := http.ListenAndServe(":"+port, corsMiddleware(r)); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
