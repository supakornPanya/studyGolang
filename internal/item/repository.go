package item

import (
	"errors"
	"sync"
)

// Repository defines the interface for item storage operations.
type Repository interface {
	Create(sku string, name string, qty int, price float64) Item
	GetAll() []Item
	GetByID(id int) (Item, error)
	Update(id int, sku string, name string, qty int, price float64) (Item, error)
	Delete(id int) error
}

// Create struct
type MemoryRepository struct {
	mu     sync.RWMutex
	items  []Item
	nextID int
}

// Create NewRepository
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		nextID: 3,
		items: []Item{
			{ID: 1, SKU: "SKU001", Name: "Item 1", Quantity: 10, Price: 100},
			{ID: 2, SKU: "SKU002", Name: "Item 2", Quantity: 20, Price: 200},
		},
	}
}

func (r *MemoryRepository) Create(sku string, name string, qty int, price float64) Item {
	r.mu.Lock()
	defer r.mu.Unlock()

	newItem := Item{
		ID:       r.nextID,
		SKU:      sku,
		Name:     name,
		Quantity: qty,
		Price:    price,
	}

	r.items = append(r.items, newItem)
	r.nextID++
	return newItem
}

func (r *MemoryRepository) GetAll() []Item {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.items
}
func (r *MemoryRepository) GetByID(id int) (Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, item := range r.items {
		if item.ID == id {
			return item, nil
		}
	}
	return Item{}, errors.New("item not found")
}
func (r *MemoryRepository) Update(id int, sku string, name string, qty int, price float64) (Item, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, item := range r.items {
		if item.ID == id {
			if name != "" {
				r.items[i].Name = name
			}
			if sku != "" {
				r.items[i].SKU = sku
			}
			if qty > 0 {
				r.items[i].Quantity = qty
			}
			if price > 0 {
				r.items[i].Price = price
			}
			return r.items[i], nil
		}
	}
	return Item{}, errors.New("item not found")
}

func (r *MemoryRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, item := range r.items {
		if item.ID == id {
			r.items = append(r.items[:i], r.items[i+1:]...)
			return nil
		}
	}
	return errors.New("item not found")
}
