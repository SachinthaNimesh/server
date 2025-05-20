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
	router := mux.NewRouter()

	// CORS Setup with proper configuration
	corsMiddleware := handlers.CORS(
		handlers.AllowedHeaders([]string{
			"Content-Type",
			"Authorization",

			"Access-Control-Allow-Headers",
			"Access-Control-Allow-Origin",
			"Origin",
			"Accept",
			"X-Requested-With",
			"Test-Key",
			"student-id", // Ensure this header is explicitly allowed
		}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}),
		handlers.AllowedOrigins([]string{"https://87abc270-1269-4d98-8dad-e53781a1ae52.e1-us-east-azure.choreoapps.dev"}),
		handlers.AllowCredentials(),
		handlers.ExposedHeaders([]string{
			"Content-Length",
			"Access-Control-Allow-Headers", // Add this header to exposed headers
			"student-id",                   // Ensure this header is explicitly exposed
		}),
		handlers.MaxAge(86400), // 24 hours
	)
	router.Use(corsMiddleware)

	// Add OPTIONS handler for preflight requests
	router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Register API routes
	routes.RegisterStudentRoutes(router)

	// Start the server
	log.Println("Server started on port", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
