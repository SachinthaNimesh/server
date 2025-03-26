package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	// Fetching database credentials from environment variables
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	sslmode := os.Getenv("DB_SSLMODE")    // Optional: Default to "require" or "verify-ca"
	sslrootcert := "/mnt/pg_creds/ca.pem" // Mounted certificate path

	if host == "" || user == "" || password == "" || dbname == "" || port == "" {
		log.Fatal("❌ Database connection environment variables are not set properly")
	}

	// Define DSN connection string with SSL root certificate
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s sslrootcert=%s",
		host,
		user,
		password,
		dbname,
		port,
		sslmode,
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
