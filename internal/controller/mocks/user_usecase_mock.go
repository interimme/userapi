package mocks

import (
	"github.com/interimme/userapi/internal/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// UserUseCase is a mock type for the UserUseCase interface
type UserUseCase struct {
	mock.Mock
}

func (m *UserUseCase) CreateUser(user *entity.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *UserUseCase) GetUser(id uuid.UUID) (*entity.User, error) {
	args := m.Called(id)
	if user, ok := args.Get(0).(*entity.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *UserUseCase) UpdateUser(user *entity.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *UserUseCase) DeleteUser(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}
