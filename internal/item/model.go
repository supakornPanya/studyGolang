package item

// Create a struct for items
type Item struct {
	ID       int     `json:"id"`
	SKU      string  `json:"sku"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

// Create template for payload that is json can't using in code like a struct
type CreateItemPayload struct {
	SKU      string  `json:"sku" binding:"required"`
	Name     string  `json:"name" binding:"required"`
	Quantity int     `json:"quantity" binding:"required,gt=0"`
	Price    float64 `json:"price" binding:"required,gt=0"`
}

//Create template for payload that is json can't using in code like a struct
type UpdateItemPayload struct {
	SKU      string  `json:"sku"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}
