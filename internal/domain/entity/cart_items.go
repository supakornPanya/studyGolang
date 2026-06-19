package entity

import "time"

// Create a struct for cart
type CartItems struct {
	UserID        int   `json:"user_id"`
	ListProductID []uint64 `gorm:"type:json;serializer:json" json:"product_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
