package repository

import (
	"study-golang-backend/internal/domain/entity"
)

type ProductRepository interface {
	Create(entity.Product) (*entity.Product, error)
	GetAll() ([]*entity.Product, error)
	GetByID(id uint64) (*entity.Product, error)
	Update(uint64, entity.Product) (*entity.Product, error)
	Delete(id uint64) error
}
