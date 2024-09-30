package mocks

import (
	"github.com/interimme/userapi/internal/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// UserRepository is a mock type for the UserRepository interface
type UserRepository struct {
	mock.Mock
}

func (m *UserRepository) Create(user *entity.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *UserRepository) GetByID(id uuid.UUID) (*entity.User, error) {
	args := m.Called(id)
	if user, ok := args.Get(0).(*entity.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *UserRepository) GetByEmail(email string) (*entity.User, error) {
	args := m.Called(email)
	if user, ok := args.Get(0).(*entity.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *UserRepository) Update(user *entity.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *UserRepository) Delete(user *entity.User) error {
	args := m.Called(user)
	return args.Error(0)
}
