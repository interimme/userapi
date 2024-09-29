package main

import (
	"fmt"
	"log"
	"userapi/internal/controller"
	"userapi/internal/infrastructure"
	"userapi/internal/infrastructure/db"
	"userapi/internal/infrastructure/persistence"
	"userapi/internal/usecase"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// Build the database connection string (DSN) with hardcoded values
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		"db",       // DB_HOST
		"postgres", // DB_USER
		"postgres", // DB_PASSWORD TODO: Handle sensitive information using Docker Secretes or Set environment variables in your CI/CD tool
		"usersdb",  // DB_NAME
		"5432",     // DB_PORT
	)

	// Connect to the database using the new db package
	dbConn, err := db.Connect(dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Migrate the database schema
	err = dbConn.AutoMigrate(&persistence.UserGorm{})
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	// Initialize the repository
	userRepo := persistence.NewUserRepository(dbConn)

	// Initialize the use case
	userUseCase := usecase.NewUserUseCase(userRepo)

	// Initialize the controller
	userController := controller.NewUserController(userUseCase)

	// Initialize the router with dependencies
	router := infrastructure.NewRouter(userController)

	// Start the server
	if err := router.Run(":8080"); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}
