package repository

import (
	"errors"
	"study-golang-backend/internal/domain/entity"
	"gorm.io/gorm"
)

// Repository interface
type UserRepository interface {
	Create(u *entity.User) error
	GetByUsername(username string) (*entity.User, error)
}

// PostgresRepository implement Repository
type UserPostgreRepository struct {
	db *gorm.DB
}

//Dependency Injection
func NewUserPostgreRepository(db *gorm.DB) UserRepository {
	return &UserPostgreRepository{db: db}
}

// Implement Create
func (r *UserPostgreRepository) Create(u *entity.User) error {
	return r.db.Create(u).Error
}

// Implement GetByUsername
func (r *UserPostgreRepository) GetByUsername(username string) (*entity.User, error) {
	var u entity.User
	err := r.db.Where("username = ?", username).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("User not found")
		}
		return nil, err
	}
	return &u, nil
}