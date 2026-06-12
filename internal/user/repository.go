package user

import (
	"gorm.io/gorm"
)

type PostgresRepository struct {
	db *gorm.DB
}

// Return PostgresRepository
func NewPostgresRepository(db *gorm.DB) Repository {
	return &PostgresRepository{db: db}
}

// Create new user
func (r *PostgresRepository) Create(u *User) error {
	return r.db.Create(u).Error
}

// Get user by username
func (r *PostgresRepository) GetByUsername(username string) (*User, error) {
	var u User
	err := r.db.Where("username = ?", username).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}
