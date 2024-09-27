package persistence

import (
	"time"
	"userapi/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserGorm represents the GORM model for the User entity
type UserGorm struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	Firstname string
	Lastname  string
	Email     string `gorm:"uniqueIndex"`
	Age       uint
	Created   time.Time
}

// UserRepository handles database operations for users
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository instance
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db}
}

// Create inserts a new user into the database
func (r *UserRepository) Create(user *entity.User) error {
	userGorm := fromEntity(user)
	return r.db.Create(&userGorm).Error
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id uuid.UUID) (*entity.User, error) {
	var userGorm UserGorm
	if err := r.db.First(&userGorm, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return toEntity(&userGorm), nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(email string) (*entity.User, error) {
	var userGorm UserGorm
	if err := r.db.First(&userGorm, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return toEntity(&userGorm), nil
}

// Update modifies an existing user in the database
func (r *UserRepository) Update(user *entity.User) error {
	userGorm := fromEntity(user)
	return r.db.Save(&userGorm).Error
}

// Delete removes a user from the database
func (r *UserRepository) Delete(user *entity.User) error {
	userGorm := fromEntity(user)
	return r.db.Delete(&userGorm).Error
}

// Helper function to convert from entity.User to UserGorm
func fromEntity(user *entity.User) UserGorm {
	return UserGorm{
		ID:        user.ID,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		Age:       user.Age,
		Created:   user.Created,
	}
}

// Helper function to convert from UserGorm to entity.User
func toEntity(userGorm *UserGorm) *entity.User {
	return &entity.User{
		ID:        userGorm.ID,
		Firstname: userGorm.Firstname,
		Lastname:  userGorm.Lastname,
		Email:     userGorm.Email,
		Age:       userGorm.Age,
		Created:   userGorm.Created,
	}
}
