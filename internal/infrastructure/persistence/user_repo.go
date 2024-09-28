package persistence

import (
	"time"
	"userapi/internal/entity"
	"userapi/internal/usecase"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserGorm represents the GORM model for the User entity
type UserGorm struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Firstname string
	Lastname  string
	Email     string `gorm:"uniqueIndex"`
	Age       uint
	Created   int64
}

// ToEntity converts UserGorm to entity.User
func (ug *UserGorm) ToEntity() *entity.User {
	return &entity.User{
		ID:        ug.ID,
		Firstname: ug.Firstname,
		Lastname:  ug.Lastname,
		Email:     ug.Email,
		Age:       ug.Age,
		Created:   ug.CreatedTime(),
	}
}

// FromEntity updates UserGorm fields from entity.User
func (ug *UserGorm) FromEntity(user *entity.User) {
	ug.ID = user.ID
	ug.Firstname = user.Firstname
	ug.Lastname = user.Lastname
	ug.Email = user.Email
	ug.Age = user.Age
	ug.Created = user.Created.Unix()
}

// CreatedTime returns the Created field as time.Time
func (ug *UserGorm) CreatedTime() time.Time {
	return time.Unix(ug.Created, 0)
}

// userRepository implements the UserRepository interface
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db *gorm.DB) usecase.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(user *entity.User) error {
	ug := &UserGorm{}
	ug.FromEntity(user)
	return r.db.Create(ug).Error
}

func (r *userRepository) GetByID(id uuid.UUID) (*entity.User, error) {
	var ug UserGorm
	if err := r.db.First(&ug, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return ug.ToEntity(), nil
}

func (r *userRepository) GetByEmail(email string) (*entity.User, error) {
	var ug UserGorm
	if err := r.db.First(&ug, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return ug.ToEntity(), nil
}

func (r *userRepository) Update(user *entity.User) error {
	ug := &UserGorm{}
	ug.FromEntity(user)
	return r.db.Save(ug).Error
}

func (r *userRepository) Delete(user *entity.User) error {
	ug := &UserGorm{}
	ug.FromEntity(user)
	return r.db.Delete(ug).Error
}
