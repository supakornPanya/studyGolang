package repository

import (
	"study-golang-backend/internal/domain/entity"
)

type TagRepository interface {
	Create(entity.Tag) (*entity.Tag, error)
	GetAll() ([]*entity.Tag, error)
	GetByID(id uint64) (*entity.Tag, error)
	Update(uint64, entity.Tag) (*entity.Tag, error)
	Delete(id uint64) error
}