package usecase

import (
	"errors"
	"net/http"
	"time"
	"userapi/internal/entity"
	appErrors "userapi/internal/errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userUseCase struct {
	repo UserRepository
}

func NewUserUseCase(repo UserRepository) UserUseCase {
	return &userUseCase{
		repo: repo,
	}
}

func (uc *userUseCase) CreateUser(user *entity.User) error {
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
