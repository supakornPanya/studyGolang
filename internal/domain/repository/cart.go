package repository

import (
	"study-golang-backend/internal/domain/entity"
)

// Repository defines the interface for item storage operations.
type CartRepository interface {
	Create(sku string, name string, qty int, price float64) (*entity.Item, error)
	GetAll() ([]*entity.Item, error)
	GetByID(id int) (*entity.Item, error)
	Update(id int, sku string, name string, qty int, price float64) (*entity.Item, error)
	Delete(id int) error
}
