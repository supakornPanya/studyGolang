package repository

import (
	"study-golang-backend/internal/domain/entity"
	"study-golang-backend/internal/domain/repository"

	"gorm.io/gorm"
)

type TagPostgresRepository struct {
	db *gorm.DB
}

// func return db
func NewTagPostgresRepository(db *gorm.DB) repository.TagRepository {
	return &TagPostgresRepository{db: db}
}

// Create
func (r *TagPostgresRepository) Create(tag entity.Tag) (*entity.Tag, error) {
	newItem := entity.Tag{
		Name: tag.Name,
	}
	err := r.db.Create(&newItem).Error
	return &newItem, err
}

// GetAll
func (r *TagPostgresRepository) GetAll() ([]*entity.Tag, error) {
	var items []*entity.Tag
	err := r.db.Find(&items).Error
	return items, err
}

// Get by ID
func (r *TagPostgresRepository) GetByID(id uint64) (*entity.Tag, error) {
	var item entity.Tag
	err := r.db.First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// Update
func (r *TagPostgresRepository) Update(id uint64, tag entity.Tag) (*entity.Tag, error) {
	var item entity.Tag
	err := r.db.First(&item, id).Error
	if err != nil {
		return nil, err
	}
	if tag.Name != "" {
		item.Name = tag.Name
	}
	return &item, r.db.Save(&item).Error
}

// Delete
func (r *TagPostgresRepository) Delete(id uint64) error {
	var item entity.Tag
	err := r.db.First(&item, id).Error
	if err != nil {
		return err
	}
	return r.db.Delete(&item).Error
}
