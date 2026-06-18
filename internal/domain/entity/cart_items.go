package entity

import "time"

// Create a struct for cart
type CartItems struct {
	ID            int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        uint64    `gorm:"uniqueIndex;not null" json:"user_id"`
	ListProductID []uint64  `gorm:"type:json;not null" json:"product_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Create template for payload that is json can't using in code like a struct
type CartItemsPayload struct {
	UserID        uint64   `gorm:"user_id" binding:"required"`
	ListProductID []uint64 `gorm:"list_product_id" binding:"required"`
}
