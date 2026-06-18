package repository
import (
	"study-golang-backend/internal/domain/entity"
)

type CartRepository interface {
	Create(cart *entity.CartItems) (*entity.CartItems, error)
	GetByID(id uint64) (*entity.CartItems, error)
	GetAll() ([]*entity.CartItems, error)
	Update(id uint64, cart *entity.CartItems) (*entity.CartItems, error)
	Delete(id uint64) error
}