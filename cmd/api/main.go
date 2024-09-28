package main

import (
	"fmt"
	"log"
	"time"

	"userapi/internal/infrastructure"
	"userapi/internal/infrastructure/persistence"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Build the database connection string (DSN) with hardcoded values
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		"db",       // DB_HOST
		"postgres", // DB_USER
		"postgres", // DB_PASSWORD // TODO:
		"usersdb",  // DB_NAME
		"5432",     // DB_PORT
	)

	// Implement retry logic to wait for the database to be ready
	var db *gorm.DB
	var err error
	maxAttempts := 10
	for attempts := 1; attempts <= maxAttempts; attempts++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("Attempt %d: Unable to connect to database. Retrying in 2 seconds...", attempts)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		panic("failed to connect to database after multiple attempts")
	}

	// Migrate the database schema
	err = db.AutoMigrate(&persistence.UserGorm{})
	if err != nil {
		panic("failed to migrate database")
	}

	// Initialize the router with dependencies
	router := infrastructure.NewRouter(db)

	// Start the server
	router.Run(":8080")
}
