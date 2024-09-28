package usecase

import (
	"userapi/internal/entity"

	"github.com/google/uuid"
)

// UserUseCase interface defines the methods for user use cases
type UserUseCase interface {
	CreateUser(user *entity.User) error
	GetUser(id uuid.UUID) (*entity.User, error)
	UpdateUser(user *entity.User) error
	DeleteUser(id uuid.UUID) error
}

// UserRepository interface defines the methods that any
// data storage provider must implement to get and store users
type UserRepository interface {
	Create(user *entity.User) error
	GetByID(id uuid.UUID) (*entity.User, error)
	GetByEmail(email string) (*entity.User, error)
	Update(user *entity.User) error
	Delete(user *entity.User) error
}
