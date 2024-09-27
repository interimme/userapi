package usecase

import (
	"errors"
	"time"
	"userapi/internal/entity"

	"github.com/google/uuid"
)

// UserRepository defines the expected behavior from a user repository
type UserRepository interface {
	Create(user *entity.User) error
	GetByID(id uuid.UUID) (*entity.User, error)
	GetByEmail(email string) (*entity.User, error)
	Update(user *entity.User) error
	Delete(user *entity.User) error
}

// UserUseCase defines the business logic for user operations
type UserUseCase interface {
	CreateUser(user *entity.User) error
	GetUser(id uuid.UUID) (*entity.User, error)
	UpdateUser(user *entity.User) error
	DeleteUser(id uuid.UUID) error
}

// userUseCase implements the UserUseCase interface
type userUseCase struct {
	repo UserRepository
}

// NewUserUseCase creates a new instance of UserUseCase
func NewUserUseCase(repo UserRepository) UserUseCase {
	return &userUseCase{
		repo: repo,
	}
}

// CreateUser handles the business logic for creating a user
func (uc *userUseCase) CreateUser(user *entity.User) error {
	// Set the user ID and creation timestamp
	user.ID = uuid.New()
	user.Created = time.Now().UTC()

	// Validate the user entity
	if err := user.Validate(); err != nil {
		return err
	}

	// Check if the email already exists
	existingUser, _ := uc.repo.GetByEmail(user.Email)
	if existingUser != nil {
		return errors.New("email already exists")
	}

	// Create the user in the repository
	return uc.repo.Create(user)
}

// GetUser retrieves a user by ID
func (uc *userUseCase) GetUser(id uuid.UUID) (*entity.User, error) {
	return uc.repo.GetByID(id)
}

// UpdateUser updates an existing user's information
func (uc *userUseCase) UpdateUser(user *entity.User) error {
	// Validate the user entity
	if err := user.Validate(); err != nil {
		return err
	}

	// Retrieve the existing user
	existingUser, err := uc.repo.GetByID(user.ID)
	if err != nil {
		return errors.New("user not found")
	}

	// Update the user's information
	existingUser.Firstname = user.Firstname
	existingUser.Lastname = user.Lastname
	existingUser.Email = user.Email
	existingUser.Age = user.Age

	// Save the updated user
	return uc.repo.Update(existingUser)
}

// DeleteUser removes a user from the repository
func (uc *userUseCase) DeleteUser(id uuid.UUID) error {
	user, err := uc.repo.GetByID(id)
	if err != nil {
		return errors.New("user not found")
	}
	return uc.repo.Delete(user)
}
