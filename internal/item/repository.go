package item

import (
	"gorm.io/gorm"
)

// Repository defines the interface for item storage operations.
type Repository interface {
	Create(sku string, name string, qty int, price float64) (*Item, error)
	GetAll() ([]*Item, error)
	GetByID(id int) (*Item, error)
	Update(id int, sku string, name string, qty int, price float64) (*Item, error)
	Delete(id int) error
}

type PostgresRepository struct {
	db *gorm.DB
}

// func return db
func NewPostgresRepository(db *gorm.DB) Repository {
	return &PostgresRepository{db: db}
}

// Create
func (r *PostgresRepository) Create(sku string, name string, qty int, price float64) (*Item, error) {
	newItem := Item{
		SKU:      sku,
		Name:     name,
		Quantity: qty,
		Price:    price,
	}
	err := r.db.Create(&newItem).Error
	return &newItem, err
}

// GetAll
func (r *PostgresRepository) GetAll() ([]*Item, error) {
	var items []*Item
	err := r.db.Find(&items).Error
	return items, err
}

// Get by ID
func (r *PostgresRepository) GetByID(id int) (*Item, error) {
	var item Item
	err := r.db.First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// Update
func (r *PostgresRepository) Update(id int, sku string, name string, qty int, price float64) (*Item, error) {
	var item Item
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
func (r *PostgresRepository) Delete(id int) error {
	var item Item
	err := r.db.First(&item, id).Error
	if err != nil {
		return err
	}
	return r.db.Delete(&item).Error
}
