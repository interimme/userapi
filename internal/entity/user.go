package entity

import (
	"errors"
	"regexp"
	"time"

	"github.com/google/uuid"
)

// User represents the user entity in the application
type User struct {
	ID        uuid.UUID // Unique identifier
	Firstname string    // User's first name
	Lastname  string    // User's last name
	Email     string    // User's email address
	Age       uint      // User's age
	Created   time.Time // Timestamp of user creation
}

// Validate checks the fields of the User entity for correctness
func (u *User) Validate() error {
	if u.Firstname == "" {
		return errors.New("firstname is required")
	}
	if u.Lastname == "" {
		return errors.New("lastname is required")
	}
	if u.Email == "" {
		return errors.New("email is required")
	}
	if !isValidEmail(u.Email) {
		return errors.New("invalid email format")
	}
	if u.Age == 0 || u.Age > 150 {
		return errors.New("age must be between 1 and 150")
	}
	return nil
}

// Helper function to validate email format using a regex
func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	return re.MatchString(email)
}
