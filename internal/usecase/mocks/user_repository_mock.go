package mocks

import (
	"userapi/internal/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// UserRepositoryMock is a mock implementation of UserRepository
type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) Create(user *entity.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *UserRepositoryMock) GetByID(id uuid.UUID) (*entity.User, error) {
	args := m.Called(id)
	if result := args.Get(0); result != nil {
		return result.(*entity.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *UserRepositoryMock) GetByEmail(email string) (*entity.User, error) {
	args := m.Called(email)
	if result := args.Get(0); result != nil {
		return result.(*entity.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *UserRepositoryMock) Update(user *entity.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *UserRepositoryMock) Delete(user *entity.User) error {
	args := m.Called(user)
	return args.Error(0)
}
