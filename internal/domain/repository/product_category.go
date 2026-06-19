package repository

import (
	"study-golang-backend/internal/domain/entity"
)

type CategoryRepository interface {
	Create(entity.Category) (*entity.Category, error)
	GetAll() ([]*entity.Category, error)
	GetByID(id uint64) (*entity.Category, error)
	Update(uint64, entity.Category) (*entity.Category, error)
	Delete(id uint64) error
}