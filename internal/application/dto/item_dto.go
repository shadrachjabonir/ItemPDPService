package dto

import "time"

// CreateItemRequest represents the request to create a new item
type CreateItemRequest struct {
	SKU         string            `json:"sku" validate:"required,min=3,max=20"`
	Name        string            `json:"name" validate:"required,min=1,max=255"`
	Description string            `json:"description" validate:"max=1000"`
	Price       float64           `json:"price" validate:"required,min=0"`
	Currency    string            `json:"currency" validate:"len=3"`
	Category    string            `json:"category" validate:"required,min=1,max=100"`
	Inventory   int               `json:"inventory" validate:"min=0"`
	Attributes  map[string]string `json:"attributes,omitempty"`
}

// UpdateItemRequest represents the request to update an item
type UpdateItemRequest struct {
	Name        *string           `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string           `json:"description,omitempty" validate:"omitempty,max=1000"`
	Price       *float64          `json:"price,omitempty" validate:"omitempty,min=0"`
	Currency    *string           `json:"currency,omitempty" validate:"omitempty,len=3"`
	Category    *string           `json:"category,omitempty" validate:"omitempty,min=1,max=100"`
	Attributes  map[string]string `json:"attributes,omitempty"`
}

// UpdateInventoryRequest represents the request to update item inventory
type UpdateInventoryRequest struct {
	Quantity int `json:"quantity" validate:"min=0"`
}

// AddImageRequest represents the request to add an image to an item
type AddImageRequest struct {
	URL       string `json:"url" validate:"required,url"`
	Alt       string `json:"alt" validate:"max=255"`
	IsPrimary bool   `json:"is_primary"`
}

// ItemResponse represents the response for item queries
type ItemResponse struct {
	ID          string            `json:"id"`
	SKU         string            `json:"sku"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Price       float64           `json:"price"`
	Currency    string            `json:"currency"`
	Category    CategoryResponse  `json:"category"`
	Inventory   InventoryResponse `json:"inventory"`
	Images      []ImageResponse   `json:"images"`
	Attributes  map[string]string `json:"attributes"`
	Status      string            `json:"status"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// CategoryResponse represents category information in responses
type CategoryResponse struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// InventoryResponse represents inventory information in responses
type InventoryResponse struct {
	Quantity    int  `json:"quantity"`
	IsAvailable bool `json:"is_available"`
}

// ImageResponse represents image information in responses
type ImageResponse struct {
	URL       string `json:"url"`
	Alt       string `json:"alt"`
	IsPrimary bool   `json:"is_primary"`
}

// ItemListResponse represents paginated list of items
type ItemListResponse struct {
	Items      []ItemResponse `json:"items"`
	Total      int            `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

// SearchRequest represents search parameters
type SearchRequest struct {
	Query    string `json:"query,omitempty"`
	Category string `json:"category,omitempty"`
	Status   string `json:"status,omitempty"`
	Page     int    `json:"page" validate:"min=1"`
	PageSize int    `json:"page_size" validate:"min=1,max=100"`
}

// ItemSummaryResponse represents a lightweight item response for lists
type ItemSummaryResponse struct {
	ID       string  `json:"id"`
	SKU      string  `json:"sku"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Currency string  `json:"currency"`
	Category string  `json:"category"`
	Status   string  `json:"status"`
	InStock  bool    `json:"in_stock"`
} 