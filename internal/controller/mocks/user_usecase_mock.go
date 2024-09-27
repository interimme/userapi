package mocks

import (
	"userapi/internal/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// UserUseCaseMock is a mock implementation of UserUseCase
type UserUseCaseMock struct {
	mock.Mock
}

func (m *UserUseCaseMock) CreateUser(user *entity.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *UserUseCaseMock) GetUser(id uuid.UUID) (*entity.User, error) {
	args := m.Called(id)
	if result := args.Get(0); result != nil {
		return result.(*entity.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *UserUseCaseMock) UpdateUser(user *entity.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *UserUseCaseMock) DeleteUser(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}
