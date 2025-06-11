package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv" // Import godotenv
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ Could not load .env file, relying on system environment variables")
	}

	// Fetching database credentials from environment variables
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")

	// Set the timezone for the database connection
	timezone := "Asia/Kolkata"
	os.Setenv("PGTZ", timezone)
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	sslrootcert := "/server/config/ca.pem"

	// sslrootcert := "C:/worky/server/config/ca.pem"

	if host == "" || user == "" || password == "" || dbname == "" || port == "" {
		log.Fatal("❌ Database connection environment variables are not set properly")
	}

	// Define DSN connection string with SSL root certificate
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=verify-ca sslrootcert=%s",
		host,
		user,
		password,
		dbname,
		port,
		sslrootcert,
	)

	log.Println("ℹ️ Attempting to connect to the database...")

	// Open a new database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	log.Println("✅ Database connection established successfully!")

	DB = db
	log.Println("✅ Database connected successfully!")
}

// Initialize sets up the database connection
func Initialize(connectionString string) {
	var err error

	DB, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection established")
}
