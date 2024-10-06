package usecase

import (
	"errors"
	appErrors "github.com/interimme/userapi/internal/apperrors"
	"github.com/interimme/userapi/internal/entity"
	"net/http"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserUseCase struct implements the methods required by the controller's UserUseCase interface
type UserUseCase struct {
	repo UserRepository
}

// NewUserUseCase creates a new instance of UserUseCase
func NewUserUseCase(repo UserRepository) *UserUseCase {
	return &UserUseCase{
		repo: repo,
	}
}

// CreateUser creates a new user
func (uc *UserUseCase) CreateUser(user *entity.User) error {
	user.ID = uuid.New()
	user.Created = time.Now().UTC()

	// Validate the user entity
	if err := user.Validate(); err != nil {
		return &appErrors.AppError{Code: http.StatusBadRequest, Message: err.Error()}
	}

	// Check if the email already exists
	existingUser, _ := uc.repo.GetByEmail(user.Email)
	if existingUser != nil {
		return appErrors.ErrConflict
	}

	// Create the user in the repository
	if err := uc.repo.Create(user); err != nil {
		return appErrors.ErrInternalServerError
	}

	return nil
}

// GetUser retrieves a user by ID
func (uc *UserUseCase) GetUser(id uuid.UUID) (*entity.User, error) {
	user, err := uc.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.ErrNotFound
		}
		return nil, appErrors.ErrInternalServerError
	}
	return user, nil
}

// UpdateUser updates an existing user
func (uc *UserUseCase) UpdateUser(user *entity.User) error {
	// Validate the user entity
	if err := user.Validate(); err != nil {
		return &appErrors.AppError{Code: http.StatusBadRequest, Message: err.Error()}
	}

	// Retrieve the existing user
	existingUser, err := uc.repo.GetByID(user.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return appErrors.ErrNotFound
		}
		return appErrors.ErrInternalServerError
	}

	// Update the user's information
	existingUser.Firstname = user.Firstname
	existingUser.Lastname = user.Lastname
	existingUser.Email = user.Email
	existingUser.Age = user.Age

	// Save the updated user
	if err := uc.repo.Update(existingUser); err != nil {
		return appErrors.ErrInternalServerError
	}

	return nil
}

// DeleteUser deletes a user by ID
func (uc *UserUseCase) DeleteUser(id uuid.UUID) error {
	user, err := uc.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return appErrors.ErrNotFound
		}
		return appErrors.ErrInternalServerError
	}

	if err := uc.repo.Delete(user); err != nil {
		return appErrors.ErrInternalServerError
	}

	return nil
}
