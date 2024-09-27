package usecase

import (
	"errors"
	"gorm.io/gorm"
	"net/http"
	"time"
	"userapi/internal/entity"
	appErrors "userapi/internal/errors"

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
		return &appErrors.AppError{Code: http.StatusBadRequest, Message: err.Error()}
	}

	// Check if the email already exists
	existingUser, _ := uc.repo.GetByEmail(user.Email)
	if existingUser != nil {
		return &appErrors.AppError{Code: http.StatusConflict, Message: "Email already exists"}
	}

	// Create the user in the repository
	if err := uc.repo.Create(user); err != nil {
		return &appErrors.AppError{Code: http.StatusInternalServerError, Message: "Failed to create user"}
	}

	return nil
}

func (uc *userUseCase) GetUser(id uuid.UUID) (*entity.User, error) {
	user, err := uc.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.ErrNotFound
		}
		return nil, appErrors.ErrInternalServerError
	}
	return user, nil
}

func (uc *userUseCase) UpdateUser(user *entity.User) error {
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

func (uc *userUseCase) DeleteUser(id uuid.UUID) error {
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
