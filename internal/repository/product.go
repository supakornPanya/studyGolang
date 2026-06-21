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
func (r *ProductPostgresRepository) Update(id uint64, product entity.Product, userID uint64, role string) (*entity.Product, error) {
	updateData := make(map[string]interface{})

	if product.SKU != "" { updateData["sku"] = product.SKU }
	if product.Name != "" { updateData["name"] = product.Name }
	if product.Description != "" { updateData["description"] = product.Description }
	if product.Price != 0 { updateData["price"] = product.Price }
	if product.Stock != 0 { updateData["stock"] = product.Stock }
	if len(product.Category) > 0 { updateData["category"] = product.Category }
	if len(product.Tag) > 0 { updateData["tag"] = product.Tag }

	// Create query base on id
	query := r.db.Model(&entity.Product{}).Where("id = ?", id)
	// if role == seller => add where created_by = userID
	if role == "seller" {
		query = query.Where("created_by = ?", userID)
	}

	result := query.Updates(updateData)

	// check error
	if result.Error != nil {
		return nil, result.Error
	}
	
	// Get the updated product
	var updatedProduct entity.Product
	if err := r.db.First(&updatedProduct, id).Error; err != nil {
		return nil, err
	}

	return &updatedProduct, nil
}

// Delete
func (r *ProductPostgresRepository) Delete(id uint64, userID uint64, role string) error {
	var item entity.Product
	// Create query base on id
	query := r.db.Model(&entity.Product{}).Where("id = ?", id)
	// if role == seller => add where created_by = userID
	if role == "seller" {
		query = query.Where("created_by = ?", userID)
	}
	err := query.Delete(&item).Error
	if err != nil {
		return err
	}
	return r.db.Delete(&item).Error
}
