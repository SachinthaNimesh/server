package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	if host == "" || user == "" || password == "" || dbname == "" || port == "" {
		log.Fatal("Database connection environment variables are not set properly")
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host,
		user,
		password,
		dbname,
		port,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Fail to connect to database:", err)
	}

	DB = db

	// Register the geom.Point type
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Fail to get sql.DB from gorm.DB:", err)
	}

	connConfig, err := pgxpool.ParseConfig(sqlDB.DSN())
	if err != nil {
		log.Fatal("Fail to parse pgxpool config:", err)
	}

	connConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		conn.ConnInfo().RegisterDataType(&pgtype.DataType{
			Value: &pgtype.Geometry{},
			Name:  "geometry",
			OID:   17515,
		})
		return nil
	}

	fmt.Println("Database connected successfully!")
	DB.Exec(`CREATE EXTENSION IF NOT EXISTS postgis`)
	DB.Exec(`ALTER TABLE attendances ADD COLUMN IF NOT EXISTS check_in_location geometry(Point, 4326)`)
	DB.Exec(`ALTER TABLE attendances ADD COLUMN IF NOT EXISTS check_out_location geometry(Point, 4326)`)

}
