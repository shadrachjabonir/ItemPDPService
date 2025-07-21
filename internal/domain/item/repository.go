package item

import (
	"context"
)

// Repository defines the interface for item persistence
type Repository interface {
	// Basic CRUD operations
	Save(ctx context.Context, item *Item) error
	FindByID(ctx context.Context, id ItemID) (*Item, error)
	FindBySKU(ctx context.Context, sku SKU) (*Item, error)
	Update(ctx context.Context, item *Item) error
	Delete(ctx context.Context, id ItemID) error
	
	// Query operations
	FindByCategory(ctx context.Context, category Category, limit, offset int) ([]*Item, error)
	FindByStatus(ctx context.Context, status Status, limit, offset int) ([]*Item, error)
	Search(ctx context.Context, query string, limit, offset int) ([]*Item, error)
	
	// Business-specific queries
	FindAvailableItems(ctx context.Context, limit, offset int) ([]*Item, error)
	FindItemsWithLowStock(ctx context.Context, threshold int) ([]*Item, error)
	
	// Aggregations
	CountByCategory(ctx context.Context, category Category) (int, error)
	CountByStatus(ctx context.Context, status Status) (int, error)
	
	// Existence checks
	ExistsBySKU(ctx context.Context, sku SKU) (bool, error)
	ExistsByID(ctx context.Context, id ItemID) (bool, error)
}

// ReadOnlyRepository defines a read-only interface for queries
type ReadOnlyRepository interface {
	FindByID(ctx context.Context, id ItemID) (*Item, error)
	FindBySKU(ctx context.Context, sku SKU) (*Item, error)
	FindByCategory(ctx context.Context, category Category, limit, offset int) ([]*Item, error)
	FindByStatus(ctx context.Context, status Status, limit, offset int) ([]*Item, error)
	Search(ctx context.Context, query string, limit, offset int) ([]*Item, error)
	FindAvailableItems(ctx context.Context, limit, offset int) ([]*Item, error)
	CountByCategory(ctx context.Context, category Category) (int, error)
	CountByStatus(ctx context.Context, status Status) (int, error)
	ExistsBySKU(ctx context.Context, sku SKU) (bool, error)
	ExistsByID(ctx context.Context, id ItemID) (bool, error)
} 