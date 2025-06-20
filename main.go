package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"server/controllers"
	"server/database"
	"server/models"
	"server/routes"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq" // PostgreSQL driver
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

	// Initialize DB connection
	db, err := sql.Open("postgres", "your_connection_string_here")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Define router
	router := mux.NewRouter()

	// CORS Setup with proper configuration
	corsMiddleware := handlers.CORS(
		handlers.AllowedHeaders([]string{
			"Content-Type",
			"Authorization",
			"Origin",
			"Accept",
			"X-Requested-With",
			"Test-Key",
			"testkey",
			"student-id", // Ensure this header is explicitly allowed
		}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}),
		handlers.AllowedOrigins([]string{"*"}), // Adjust as needed
		handlers.AllowCredentials(),
		handlers.ExposedHeaders([]string{
			"Content-Length",
		}),
		handlers.MaxAge(86400), // 24 hours
	)

	authService := controllers.NewAuthService()
	authService.RegisterRoutes(router)
	router.Use(corsMiddleware)

	// Register API routes
	routes.RegisterStudentRoutes(router)

	// Add middleware

	// Start the server
	log.Println("Server started on port", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
