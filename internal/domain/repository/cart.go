package repository

import (
	"study-golang-backend/internal/domain/entity"
	"gorm.io/gorm"
)

// Repository defines the interface for item storage operations.
type CartRepository interface {
	Create(sku string, name string, qty int, price float64) (*entity.Item, error)
	GetAll() ([]*entity.Item, error)
	GetByID(id int) (*entity.Item, error)
	Update(id int, sku string, name string, qty int, price float64) (*entity.Item, error)
	Delete(id int) error
}

type CartPostgresRepository struct {
	db *gorm.DB
}

// func return db
func NewCartPostgresRepository(db *gorm.DB) CartRepository {
	return &CartPostgresRepository{db: db}
}

// Create
func (r *CartPostgresRepository) Create(sku string, name string, qty int, price float64) (*entity.Item, error) {
	newItem := entity.Item{
		SKU:      sku,
		Name:     name,
		Quantity: qty,
		Price:    price,
	}
	err := r.db.Create(&newItem).Error
	return &newItem, err
}

// GetAll
func (r *CartPostgresRepository) GetAll() ([]*entity.Item, error) {
	var items []*entity.Item
	err := r.db.Find(&items).Error
	return items, err
}

// Get by ID
func (r *CartPostgresRepository) GetByID(id int) (*entity.Item, error) {
	var item entity.Item
	err := r.db.First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// Update
func (r *CartPostgresRepository) Update(id int, sku string, name string, qty int, price float64) (*entity.Item, error) {
	var item entity.Item
	err := r.db.First(&item, id).Error
	if err != nil {
		return nil, err
	}
	if sku != "" {
		item.SKU = sku
	}
	if name != "" {
		item.Name = name
	}
	if qty != 0 {
		item.Quantity = qty
	}
	if price != 0 {
		item.Price = price
	}
	return &item, r.db.Save(&item).Error
}

// Delete
func (r *CartPostgresRepository) Delete(id int) error {
	var item entity.Item
	err := r.db.First(&item, id).Error
	if err != nil {
		return err
	}
	return r.db.Delete(&item).Error
}
