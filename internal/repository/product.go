package repository

import (
	"study-golang-backend/internal/domain/entity"
	"study-golang-backend/internal/domain/repository"

	"gorm.io/gorm"
)

type ProductPostgresRepository struct {
	db *gorm.DB
}

// func return db
func NewProductPostgresRepository(db *gorm.DB) repository.ProductRepository {
	return &ProductPostgresRepository{db: db}
}

// Create
func (r *ProductPostgresRepository) Create(product entity.Product) (*entity.Product, error) {
	newItem := entity.Product{
		SKU:         product.SKU,
		Name:        product.Name,
		Description: product.Description,
		Stock:       product.Stock,
		Price:       product.Price,
		Category:    product.Category,
		Tag:         product.Tag,
		CreatedBy:   product.CreatedBy,
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
func (r *ProductPostgresRepository) GetByID(id uint64) (*entity.Product, error) {
	var item entity.Product
	err := r.db.First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// Update
func (r *ProductPostgresRepository) Update(id uint64, product entity.Product) (*entity.Product, error) {
	var item entity.Product
	err := r.db.First(&item, id).Error
	if err != nil {
		return nil, err
	}
	if product.SKU != "" {
		item.SKU = product.SKU
	}
	if product.Name != "" {
		item.Name = product.Name
	}
	if product.Description != "" {
		item.Description = product.Description
	}
	if product.Stock != 0 {
		item.Stock = product.Stock
	}
	if product.Price != 0 {
		item.Price = product.Price
	}
	if len(product.Category) > 0 {
		item.Category = product.Category
	}
	if len(product.Tag) > 0 {
		item.Tag = product.Tag
	}
	return &item, r.db.Save(&item).Error
}

// Delete
func (r *ProductPostgresRepository) Delete(id uint64) error {
	var item entity.Product
	err := r.db.First(&item, id).Error
	if err != nil {
		return err
	}
	return r.db.Delete(&item).Error
}
