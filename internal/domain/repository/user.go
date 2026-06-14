package repository

import (
	"study-golang-backend/internal/domain/entity"
)

// Repository interface
type UserRepository interface {
	Create(u *entity.User) error
	GetByUsername(username string) (*entity.User, error)
}
