package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"userapi/internal/infrastructure"
	"userapi/internal/infrastructure/persistence"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Build the database connection string (DSN)
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"), // TODO: Manage sensitive info with Docker Secrets or setting .env variables in CI/CD tool
		getEnv("DB_NAME", "usersdb"),
		getEnv("DB_PORT", "5432"),
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
	port := getEnv("PORT", "8080")
	router.Run(":" + port)
}

// Helper function to get environment variables with a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
