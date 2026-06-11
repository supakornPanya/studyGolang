package user

import (
	"errors"
	"gorm.io/gorm"
)

// Struct of User 
type User struct {
	ID           string `gorm:"primaryKey" json:"id"`
	Username     string `gorm:"uniqueIndex;not null" json:"username"`
	PasswordHash string `gorm:"not null" json:"-"`
	CanRead      bool   `gorm:"default:true" json:"can_read"`
	CanWrite     bool   `gorm:"default:false" json:"can_write"`
	CanUpdate    bool   `gorm:"default:false" json:"can_update"`
	CanDelete    bool   `gorm:"default:false" json:"can_delete"`
}

// Repository interface
type Repository interface {
	Create(u *User) error
	GetByUsername(username string) (*User, error)
}

// PostgresRepository implement Repository
type PostgresRepository struct {
	db *gorm.DB
}

//Dependency Injection
func NewPostgresRepository(db *gorm.DB) Repository {
	return &PostgresRepository{db: db}
}

// Implement Create
func (r *PostgresRepository) Create(u *User) error {
	return r.db.Create(u).Error
}

// Implement GetByUsername
func (r *PostgresRepository) GetByUsername(username string) (*User, error) {
	var u User
	err := r.db.Where("username = ?", username).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("User not found")
		}
		return nil, err
	}
	return &u, nil
}