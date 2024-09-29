package controller

import (
	"userapi/internal/entity"

	"github.com/google/uuid"
)

// UserUseCase interface defines the methods used by the controller
type UserUseCase interface {
	CreateUser(user *entity.User) error
	GetUser(id uuid.UUID) (*entity.User, error)
	UpdateUser(user *entity.User) error
	DeleteUser(id uuid.UUID) error
}
