package repository
import (
	"study-golang-backend/internal/domain/entity"
)

type CartRepository interface {
	Save(id int, cart *entity.CartItems) (*entity.CartItems, error)
	GetByID(id int) (*entity.CartItems, error)
	Delete(id int) error
}