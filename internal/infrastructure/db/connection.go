package db

import (
	"errors"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connect establishes a connection to the database using GORM.
// It returns a pointer to the GORM DB instance and any error encountered.
func Connect(dsn string) (*gorm.DB, error) {
	if dsn == "" {
		return nil, errors.New("no Postgres DSN provided")
	}

	var db *gorm.DB
	var err error
	maxAttempts := 10

	for attempts := 1; attempts <= maxAttempts; attempts++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			sqlDB, err := db.DB()
			if err != nil {
				return nil, fmt.Errorf("failed to get sqlDB from GORM: %w", err)
			}
			// Optionally set up database connection pool
			sqlDB.SetMaxOpenConns(10)
			sqlDB.SetMaxIdleConns(10)
			sqlDB.SetConnMaxLifetime(time.Minute)
			// Ping to test the connection
			if err := sqlDB.Ping(); err != nil {
				return nil, fmt.Errorf("db.Ping failed: %w", err)
			}
			return db, nil
		}
		log.Printf("Attempt %d: Unable to connect to database. Retrying in 2 seconds...", attempts)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxAttempts, err)
}
