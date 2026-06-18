package repository

import (
	"gorm.io/gorm"
	"study-golang-backend/internal/domain/entity"
	"study-golang-backend/internal/domain/repository"
)

type ProductPostgresRepository struct {
	db *gorm.DB
}

// func return db
func NewProductPostgresRepository(db *gorm.DB) repository.ProductRepository {
	return &ProductPostgresRepository{db: db}
}

// Create
func (r *ProductPostgresRepository) Create(sku string, name string, description string, qty int, price float64, category []uint64, tag []uint64, createdBy uint64) (*entity.Product, error) {
	newItem := entity.Product{
		SKU:         sku,
		Name:        name,
		Description: description,
		Stock:       qty,
		Price:       price,
		Category:    category,
		Tag:         tag,
		CreatedBy:   createdBy,
	}
	err := r.db.Create(&newItem).Error
	return &newItem, err
}

// GetAll
func (r *ProductPostgresRepository) GetAll() ([]*entity.Product, error) {
	var items []*entity.Product
	err := r.db.Find(&items).Error
	return items, err
}

// Get by ID
func (r *ProductPostgresRepository) GetByID(id int) (*entity.Product, error) {
	var item entity.Product
	err := r.db.First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// Update
func (r *ProductPostgresRepository) Update(id int, sku string, name string, description string, qty int, price float64, category []uint64, tag []uint64) (*entity.Product, error) {
	var item entity.Product
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
	if description != "" {
		item.Description = description
	}
	if qty != 0 {
		item.Stock = qty
	}
	if price != 0 {
		item.Price = price
	}
	if category != nil {
		item.Category = category
	}
	if tag != nil {
		item.Tag = tag
	}
	return &item, r.db.Save(&item).Error
}

// Delete
func (r *ProductPostgresRepository) Delete(id int) error {
	var item entity.Product
	err := r.db.First(&item, id).Error
	if err != nil {
		return err
	}
	return r.db.Delete(&item).Error
}
