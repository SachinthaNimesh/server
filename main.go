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

	// Middleware to set Access-Control-Allow-Origin header
	originMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
			next.ServeHTTP(w, r)
		})
	}
	router.Use(originMiddleware)

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
			"student-id",
		}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.ExposedHeaders([]string{
			"Content-Length",
		}),
		handlers.MaxAge(86400), // 24 hours
	)

	// Register API routes
	routes.RegisterStudentRoutes(router)

	// Start the server
	log.Println("Server started on port", port)
	if err := http.ListenAndServe(":"+port, corsMiddleware(router)); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
