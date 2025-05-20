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

	// // Route to handle location updates
	// r.HandleFunc("/location", controllers.HandleLocationUpdate).Methods("POST")

	// // Route for WebSocket connections
	// r.HandleFunc("/ws", controllers.HandleWebSocket)

	// // Serve static files for the React web app
	// r.PathPrefix("/").Handler(http.FileServer(http.Dir("./web")))

	// CORS Setup - Set allowed origins to Choreo API Gateway
	corsMiddleware := handlers.CORS(
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "Student-ID"}), // Ensure "Student-ID" is included
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}),
		handlers.AllowedOrigins([]string{"https://87abc270-1269-4d98-8dad-e53781a1ae52.e1-us-east-azure.choreoapps.dev"}), // Replace "*" with the actual origin in production
	)

	// Start the server
	log.Println("Server started on port", port)
	if err := http.ListenAndServe(":"+port, corsMiddleware(r)); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
