package entity

// 1. Product
type Product struct {
	ID          uint64  `gorm:"primaryKey;autoIncrement" json:"id"`
	SKU         string   `gorm:"not null" json:"sku"`
	Name        string   `gorm:"not null" json:"name"`
	Description string   `gorm:"not null" json:"description"`
	Price       float64  `gorm:"not null" json:"price"`
	Stock       int      `gorm:"not null" json:"stock"`
	Category    []uint64 `gorm:"type:json;serializer:json" json:"category"`
	Tag         []uint64 `gorm:"type:json;serializer:json" json:"tag"`
	CreatedBy   uint64   `gorm:"not null" json:"created_by"`
}

// 2. Tag
type Tag struct {
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"not null" json:"name"`
}

// 3. Category
type Category struct {
	ID          uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string `gorm:"not null" json:"name"`
	Description string `gorm:"not null" json:"description"`
	ParentID    uint64 `gorm:"not null" json:"parent_id"`
}
