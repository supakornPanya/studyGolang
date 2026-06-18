package repository
import (
	"study-golang-backend/internal/domain/entity"
)

type ProductRepository interface{
	Create(sku string, name string, description string, qty int, price float64, category []uint64, tag []uint64, createdBy uint64) (*entity.Product, error)
	GetAll() ([]*entity.Product, error)
	GetByID(id int) (*entity.Product, error)
	Update(id int, sku string, name string, description string, qty int, price float64, category []uint64, tag []uint64) (*entity.Product, error)
	Delete(id int) error
}