package repository

import (
	"study-golang-backend/internal/domain/entity"
	"study-golang-backend/internal/domain/repository"

	"gorm.io/gorm"
)

type CategoryPostgresRepository struct {
	db *gorm.DB
}

// func return db
func NewCategoryPostgresRepository(db *gorm.DB) repository.CategoryRepository {
	return &CategoryPostgresRepository{db: db}
}

// Create
func (r *CategoryPostgresRepository) Create(category entity.Category) (*entity.Category, error) {
	newItem := entity.Category{
		Name:        category.Name,
		Description: category.Description,
		ParentID:    category.ParentID,
	}
	err := r.db.Create(&newItem).Error
	return &newItem, err
}

// GetAll
func (r *CategoryPostgresRepository) GetAll() ([]*entity.Category, error) {
	var items []*entity.Category
	err := r.db.Find(&items).Error
	return items, err
}

// Get by ID
func (r *CategoryPostgresRepository) GetByID(id uint64) (*entity.Category, error) {
	var item entity.Category
	err := r.db.First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// Update
func (r *CategoryPostgresRepository) Update(id uint64, category entity.Category) (*entity.Category, error) {
	var item entity.Category
	err := r.db.First(&item, id).Error
	if err != nil {
		return nil, err
	}
	if category.Name != "" {
		item.Name = category.Name
	}
	if category.Description != "" {
		item.Description = category.Description
	}
	if category.ParentID != 0 {
		item.ParentID = category.ParentID
	}
	return &item, r.db.Save(&item).Error
}

// Delete
func (r *CategoryPostgresRepository) Delete(id uint64) error {
	var item entity.Category
	err := r.db.First(&item, id).Error
	if err != nil {
		return err
	}
	return r.db.Delete(&item).Error
}
