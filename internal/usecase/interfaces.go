package usecase

import (
	"github.com/interimme/userapi/internal/entity"

	"github.com/google/uuid"
)

// UserRepository interface defines the methods that any
// data storage provider must implement to get and store users
type UserRepository interface {
	Create(user *entity.User) error
	GetByID(id uuid.UUID) (*entity.User, error)
	GetByEmail(email string) (*entity.User, error)
	Update(user *entity.User) error
	Delete(user *entity.User) error
}
