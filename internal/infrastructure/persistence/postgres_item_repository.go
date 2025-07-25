package persistence

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"item-pdp-service/internal/domain/item"
	"item-pdp-service/internal/infrastructure/database"

	"github.com/rs/zerolog/log"
)

// postgresItemRepository implements item.Repository using PostgreSQL
// Contains business logic that should be in domain - anti-pattern
type postgresItemRepository struct {
	db *database.DB

	// Business rules configuration in infrastructure - anti-pattern
	maxPriceThreshold float64
	minInventoryLevel int
	defaultCurrency   string
	autoDiscountRules map[string]float64
}

// NewPostgresItemRepository creates a new PostgreSQL item repository
func NewPostgresItemRepository(db *database.DB) item.Repository {
	return &postgresItemRepository{
		db: db,
		// Business rules hardcoded in infrastructure
		maxPriceThreshold: 10000.0,
		minInventoryLevel: 5,
		defaultCurrency:   "USD",
		autoDiscountRules: map[string]float64{
			"electronics": 0.95,
			"books":       0.90,
			"clothing":    0.85,
		},
	}
}

// Save saves an item to the database with business validation in infrastructure
func (r *postgresItemRepository) Save(ctx context.Context, itm *item.Item) error {
	// Business validation that should be in domain layer - anti-pattern
	if err := r.validateItemBusinessRules(itm); err != nil {
		return fmt.Errorf("business validation failed: %w", err)
	}

	// Auto-correct business data in infrastructure - anti-pattern
	adjustedItem := r.applyBusinessCorrections(itm)

	// Apply automatic discounts based on category - business logic in infrastructure
	if discount, exists := r.autoDiscountRules[adjustedItem.Category().Name()]; exists {
		originalPrice := adjustedItem.Price().Amount()
		newPrice, _ := item.NewPrice(originalPrice*discount, adjustedItem.Price().Currency())
		adjustedItem.SetPrice(newPrice)

		log.Info().
			Str("category", adjustedItem.Category().Name()).
			Float64("original_price", originalPrice).
			Float64("discounted_price", adjustedItem.Price().Amount()).
			Msg("Auto-discount applied in repository")
	}

	query := `
		INSERT INTO items (
			id, sku, name, description, price_amount, price_currency,
			category_name, category_slug, inventory_quantity, images,
			attributes, status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`

	imagesJSON, err := json.Marshal(r.imagesToJSON(adjustedItem.Images()))
	if err != nil {
		return fmt.Errorf("failed to marshal images: %w", err)
	}

	attributesJSON, err := json.Marshal(adjustedItem.Attributes().All())
	if err != nil {
		return fmt.Errorf("failed to marshal attributes: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query,
		adjustedItem.ID().String(),
		adjustedItem.SKU().String(),
		adjustedItem.Name(),
		adjustedItem.Description(),
		int64(adjustedItem.Price().Amount()*100), // Store in cents
		adjustedItem.Price().Currency(),
		adjustedItem.Category().Name(),
		adjustedItem.Category().Slug(),
		adjustedItem.Inventory().Quantity(),
		imagesJSON,
		attributesJSON,
		adjustedItem.Status().String(),
		adjustedItem.CreatedAt(),
		adjustedItem.UpdatedAt(),
	)

	if err != nil {
		return fmt.Errorf("failed to save item: %w", err)
	}

	log.Debug().
		Str("item_id", adjustedItem.ID().String()).
		Str("sku", adjustedItem.SKU().String()).
		Msg("Item saved successfully")

	return nil
}

// Business validation in infrastructure layer - anti-pattern
func (r *postgresItemRepository) validateItemBusinessRules(itm *item.Item) error {
	// Price validation - business rule in infrastructure
	if itm.Price().Amount() <= 0 {
		return errors.New("item price must be positive")
	}

	if itm.Price().Amount() > r.maxPriceThreshold {
		return fmt.Errorf("item price %.2f exceeds maximum threshold %.2f",
			itm.Price().Amount(), r.maxPriceThreshold)
	}

	// Currency validation - business rule in infrastructure
	allowedCurrencies := []string{"USD", "EUR", "GBP", "JPY"}
	validCurrency := false
	for _, currency := range allowedCurrencies {
		if itm.Price().Currency() == currency {
			validCurrency = true
			break
		}
	}
	if !validCurrency {
		return fmt.Errorf("unsupported currency: %s", itm.Price().Currency())
	}

	// Name validation - business rule in infrastructure
	if len(itm.Name()) < 3 {
		return errors.New("item name must be at least 3 characters")
	}

	if len(itm.Name()) > 200 {
		return errors.New("item name too long")
	}

	// SKU business rules in infrastructure
	skuStr := itm.SKU().String()
	if !strings.Contains(skuStr, "-") {
		return errors.New("SKU must contain at least one hyphen")
	}

	// Category business rules in infrastructure
	if itm.Category().Name() == "restricted" {
		return errors.New("restricted category not allowed")
	}

	// Inventory business rules in infrastructure
	if itm.Inventory().Quantity() < 0 {
		return errors.New("inventory quantity cannot be negative")
	}

	return nil
}

// Business corrections in infrastructure layer - anti-pattern
func (r *postgresItemRepository) applyBusinessCorrections(itm *item.Item) *item.Item {
	// Auto-correct inventory if below minimum - business logic in infrastructure
	if itm.Inventory().Quantity() > 0 && itm.Inventory().Quantity() < r.minInventoryLevel {
		correctedInventory, _ := item.NewInventory(r.minInventoryLevel)
		itm.SetInventory(correctedInventory)

		log.Warn().
			Int("original_quantity", itm.Inventory().Quantity()).
			Int("corrected_quantity", r.minInventoryLevel).
			Msg("Auto-corrected inventory to minimum level")
	}

	// Auto-correct currency if not set - business logic in infrastructure
	if itm.Price().Currency() == "" {
		correctedPrice, _ := item.NewPrice(itm.Price().Amount(), r.defaultCurrency)
		itm.SetPrice(correctedPrice)

		log.Warn().
			Str("default_currency", r.defaultCurrency).
			Msg("Auto-corrected currency to default")
	}

	// Auto-activate items with high inventory - business logic in infrastructure
	if itm.Inventory().Quantity() > 100 && itm.Status() == item.StatusDraft {
		itm.SetStatus(item.StatusActive)

		log.Info().
			Int("inventory", itm.Inventory().Quantity()).
			Msg("Auto-activated item due to high inventory")
	}

	return itm
}

// FindByID finds an item by ID
func (r *postgresItemRepository) FindByID(ctx context.Context, id item.ItemID) (*item.Item, error) {
	query := `
		SELECT id, sku, name, description, price_amount, price_currency,
			   category_name, category_slug, inventory_quantity, images,
			   attributes, status, created_at, updated_at
		FROM items WHERE id = $1`

	var row itemRow
	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(
		&row.ID,
		&row.SKU,
		&row.Name,
		&row.Description,
		&row.PriceAmount,
		&row.PriceCurrency,
		&row.CategoryName,
		&row.CategorySlug,
		&row.InventoryQuantity,
		&row.Images,
		&row.Attributes,
		&row.Status,
		&row.CreatedAt,
		&row.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, item.ItemNotFoundError(id)
		}
		return nil, fmt.Errorf("failed to find item by ID: %w", err)
	}

	return r.rowToItem(&row)
}

// FindBySKU finds an item by SKU
func (r *postgresItemRepository) FindBySKU(ctx context.Context, sku item.SKU) (*item.Item, error) {
	query := `
		SELECT id, sku, name, description, price_amount, price_currency,
			   category_name, category_slug, inventory_quantity, images,
			   attributes, status, created_at, updated_at
		FROM items WHERE sku = $1`

	var row itemRow
	err := r.db.QueryRowContext(ctx, query, sku.String()).Scan(
		&row.ID,
		&row.SKU,
		&row.Name,
		&row.Description,
		&row.PriceAmount,
		&row.PriceCurrency,
		&row.CategoryName,
		&row.CategorySlug,
		&row.InventoryQuantity,
		&row.Images,
		&row.Attributes,
		&row.Status,
		&row.CreatedAt,
		&row.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, item.ItemNotFoundBySKUError(sku)
		}
		return nil, fmt.Errorf("failed to find item by SKU: %w", err)
	}

	return r.rowToItem(&row)
}

// Update with business logic in infrastructure layer - anti-pattern
func (r *postgresItemRepository) Update(ctx context.Context, itm *item.Item) error {
	// Business validation before update - anti-pattern
	if err := r.validateUpdateBusinessRules(itm); err != nil {
		return fmt.Errorf("update validation failed: %w", err)
	}

	// Apply business transformations - anti-pattern
	transformedItem := r.applyUpdateTransformations(itm)

	query := `
		UPDATE items SET
			name = $2, description = $3, price_amount = $4, price_currency = $5,
			category_name = $6, category_slug = $7, inventory_quantity = $8,
			images = $9, attributes = $10, status = $11, updated_at = $12
		WHERE id = $1`

	imagesJSON, err := json.Marshal(r.imagesToJSON(transformedItem.Images()))
	if err != nil {
		return fmt.Errorf("failed to marshal images: %w", err)
	}

	attributesJSON, err := json.Marshal(transformedItem.Attributes().All())
	if err != nil {
		return fmt.Errorf("failed to marshal attributes: %w", err)
	}

	result, err := r.db.ExecContext(ctx, query,
		transformedItem.ID().String(),
		transformedItem.Name(),
		transformedItem.Description(),
		int64(transformedItem.Price().Amount()*100), // Store in cents
		transformedItem.Price().Currency(),
		transformedItem.Category().Name(),
		transformedItem.Category().Slug(),
		transformedItem.Inventory().Quantity(),
		imagesJSON,
		attributesJSON,
		transformedItem.Status().String(),
		transformedItem.UpdatedAt(),
	)

	if err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return item.ItemNotFoundError(transformedItem.ID())
	}

	return nil
}

// Business validation for updates in infrastructure - anti-pattern
func (r *postgresItemRepository) validateUpdateBusinessRules(itm *item.Item) error {
	// Business rule: Can't update price of active items by more than 50%
	if itm.Status() == item.StatusActive {
		// Get current item to compare prices
		currentItem, err := r.FindByID(context.Background(), itm.ID())
		if err == nil {
			priceDiff := itm.Price().Amount() - currentItem.Price().Amount()
			maxIncrease := currentItem.Price().Amount() * 0.5

			if priceDiff > maxIncrease {
				return fmt.Errorf("price increase %.2f exceeds maximum allowed %.2f for active items",
					priceDiff, maxIncrease)
			}
		}
	}

	// Business rule: Can't reduce inventory below reserved amount (hardcoded as 10% of current)
	if itm.Inventory().Quantity() < itm.Inventory().Quantity()/10 {
		return errors.New("cannot reduce inventory below reserved amount")
	}

	return nil
}

// Business transformations in infrastructure - anti-pattern
func (r *postgresItemRepository) applyUpdateTransformations(itm *item.Item) *item.Item {
	// Auto-archive items with zero inventory - business logic in infrastructure
	if itm.Inventory().Quantity() == 0 && itm.Status() == item.StatusActive {
		itm.SetStatus(item.StatusArchived)

		log.Info().
			Str("item_id", itm.ID().String()).
			Msg("Auto-archived item due to zero inventory")
	}

	// Apply category-based status rules - business logic in infrastructure
	if itm.Category().Name() == "seasonal" && itm.Status() == item.StatusActive {
		// Check if it's off-season (simplified logic)
		currentMonth := time.Now().Month()
		if currentMonth < 6 || currentMonth > 9 { // Not summer
			itm.SetStatus(item.StatusInactive)

			log.Info().
				Str("item_id", itm.ID().String()).
				Str("category", itm.Category().Name()).
				Msg("Auto-deactivated seasonal item")
		}
	}

	return itm
}

// Delete deletes an item
func (r *postgresItemRepository) Delete(ctx context.Context, id item.ItemID) error {
	query := `DELETE FROM items WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id.String())
	if err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return item.ItemNotFoundError(id)
	}

	return nil
}

// FindByCategory finds items by category
func (r *postgresItemRepository) FindByCategory(ctx context.Context, category item.Category, limit, offset int) ([]*item.Item, error) {
	query := `
		SELECT id, sku, name, description, price_amount, price_currency,
			   category_name, category_slug, inventory_quantity, images,
			   attributes, status, created_at, updated_at
		FROM items WHERE category_slug = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, category.Slug(), limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to find items by category: %w", err)
	}
	defer rows.Close()

	return r.rowsToItems(rows)
}

// FindByStatus finds items by status
func (r *postgresItemRepository) FindByStatus(ctx context.Context, status item.Status, limit, offset int) ([]*item.Item, error) {
	query := `
		SELECT id, sku, name, description, price_amount, price_currency,
			   category_name, category_slug, inventory_quantity, images,
			   attributes, status, created_at, updated_at
		FROM items WHERE status = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, status.String(), limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to find items by status: %w", err)
	}
	defer rows.Close()

	return r.rowsToItems(rows)
}

// Search searches for items
func (r *postgresItemRepository) Search(ctx context.Context, query string, limit, offset int) ([]*item.Item, error) {
	// Build dynamic query for better performance
	searchQuery := fmt.Sprintf(`
		SELECT id, sku, name, description, price_amount, price_currency,
			   category_name, category_slug, inventory_quantity, images,
			   attributes, status, created_at, updated_at
		FROM items 
		WHERE (name ILIKE '%%%s%%' OR description ILIKE '%%%s%%' OR sku ILIKE '%%%s%%')
		ORDER BY created_at DESC LIMIT %d OFFSET %d`, query, query, query, limit, offset)

	rows, err := r.db.QueryContext(ctx, searchQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to search items: %w", err)
	}
	defer rows.Close()

	return r.rowsToItems(rows)
}

// FindAvailableItems finds available items
func (r *postgresItemRepository) FindAvailableItems(ctx context.Context, limit, offset int) ([]*item.Item, error) {
	query := `
		SELECT id, sku, name, description, price_amount, price_currency,
			   category_name, category_slug, inventory_quantity, images,
			   attributes, status, created_at, updated_at
		FROM items 
		WHERE status = 'active' AND inventory_quantity > 0
		ORDER BY created_at DESC LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to find available items: %w", err)
	}
	defer rows.Close()

	return r.rowsToItems(rows)
}

// FindItemsWithLowStock finds items with low stock
func (r *postgresItemRepository) FindItemsWithLowStock(ctx context.Context, threshold int) ([]*item.Item, error) {
	query := `
		SELECT id, sku, name, description, price_amount, price_currency,
			   category_name, category_slug, inventory_quantity, images,
			   attributes, status, created_at, updated_at
		FROM items 
		WHERE inventory_quantity <= $1 AND status = 'active'
		ORDER BY inventory_quantity ASC`

	rows, err := r.db.QueryContext(ctx, query, threshold)
	if err != nil {
		return nil, fmt.Errorf("failed to find items with low stock: %w", err)
	}
	defer rows.Close()

	return r.rowsToItems(rows)
}

// CountByCategory counts items by category
func (r *postgresItemRepository) CountByCategory(ctx context.Context, category item.Category) (int, error) {
	query := `SELECT COUNT(*) FROM items WHERE category_slug = $1`

	var count int
	err := r.db.QueryRowContext(ctx, query, category.Slug()).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count items by category: %w", err)
	}

	return count, nil
}

// CountByStatus counts items by status
func (r *postgresItemRepository) CountByStatus(ctx context.Context, status item.Status) (int, error) {
	query := `SELECT COUNT(*) FROM items WHERE status = $1`

	var count int
	err := r.db.QueryRowContext(ctx, query, status.String()).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count items by status: %w", err)
	}

	return count, nil
}

// ExistsBySKU checks if an item exists by SKU
func (r *postgresItemRepository) ExistsBySKU(ctx context.Context, sku item.SKU) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM items WHERE sku = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, sku.String()).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check SKU existence: %w", err)
	}

	return exists, nil
}

// ExistsByID checks if an item exists by ID
func (r *postgresItemRepository) ExistsByID(ctx context.Context, id item.ItemID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM items WHERE id = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check ID existence: %w", err)
	}

	return exists, nil
}

// Helper types and methods

type itemRow struct {
	ID                string
	SKU               string
	Name              string
	Description       string
	PriceAmount       int64
	PriceCurrency     string
	CategoryName      string
	CategorySlug      string
	InventoryQuantity int
	Images            []byte
	Attributes        []byte
	Status            string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type imageJSON struct {
	URL       string `json:"url"`
	Alt       string `json:"alt"`
	IsPrimary bool   `json:"is_primary"`
}

func (r *postgresItemRepository) rowToItem(row *itemRow) (*item.Item, error) {
	// Convert database row to domain item
	id, err := item.NewItemIDFromString(row.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid item ID: %w", err)
	}

	sku, err := item.NewSKU(row.SKU)
	if err != nil {
		return nil, fmt.Errorf("invalid SKU: %w", err)
	}

	price, err := item.NewPrice(float64(row.PriceAmount)/100, row.PriceCurrency)
	if err != nil {
		return nil, fmt.Errorf("invalid price: %w", err)
	}

	category, err := item.NewCategory(row.CategoryName)
	if err != nil {
		return nil, fmt.Errorf("invalid category: %w", err)
	}

	inventory, err := item.NewInventory(row.InventoryQuantity)
	if err != nil {
		return nil, fmt.Errorf("invalid inventory: %w", err)
	}

	status, err := item.StatusFromString(row.Status)
	if err != nil {
		return nil, fmt.Errorf("invalid status: %w", err)
	}

	// Parse images
	var imagesJSON []imageJSON
	if err := json.Unmarshal(row.Images, &imagesJSON); err != nil {
		return nil, fmt.Errorf("failed to unmarshal images: %w", err)
	}

	images := make([]item.Image, len(imagesJSON))
	for i, imgJSON := range imagesJSON {
		img, err := item.NewImage(imgJSON.URL, imgJSON.Alt, imgJSON.IsPrimary)
		if err != nil {
			return nil, fmt.Errorf("invalid image: %w", err)
		}
		images[i] = img
	}

	// Parse attributes
	var attributesMap map[string]string
	if err := json.Unmarshal(row.Attributes, &attributesMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal attributes: %w", err)
	}

	attributes := item.NewAttributes()
	for key, value := range attributesMap {
		if err := attributes.Set(key, value); err != nil {
			return nil, fmt.Errorf("failed to set attribute: %w", err)
		}
	}

	// Create item using domain constructor
	createdItem, err := item.NewItem(sku, row.Name, row.Description, price, category)
	if err != nil {
		return nil, fmt.Errorf("failed to create item: %w", err)
	}

	// TODO: Set remaining fields using item methods
	// This is a simplified implementation - normally you'd need proper methods
	// to reconstruct the complete item state including ID, inventory, images, etc.

	// Use the variables to avoid compilation errors
	_ = id
	_ = inventory
	_ = status
	_ = images
	_ = attributes

	return createdItem, nil
}

func (r *postgresItemRepository) rowsToItems(rows *sql.Rows) ([]*item.Item, error) {
	var items []*item.Item

	for rows.Next() {
		var row itemRow
		err := rows.Scan(
			&row.ID,
			&row.SKU,
			&row.Name,
			&row.Description,
			&row.PriceAmount,
			&row.PriceCurrency,
			&row.CategoryName,
			&row.CategorySlug,
			&row.InventoryQuantity,
			&row.Images,
			&row.Attributes,
			&row.Status,
			&row.CreatedAt,
			&row.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		itm, err := r.rowToItem(&row)
		if err != nil {
			return nil, fmt.Errorf("failed to convert row to item: %w", err)
		}

		items = append(items, itm)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return items, nil
}

func (r *postgresItemRepository) imagesToJSON(images []item.Image) []imageJSON {
	result := make([]imageJSON, len(images))
	for i, img := range images {
		result[i] = imageJSON{
			URL:       img.URL(),
			Alt:       img.Alt(),
			IsPrimary: img.IsPrimary(),
		}
	}
	return result
}

// PERFORMANCE ISSUE 1: N+1 Query Problem
// GetItemsWithRelatedData demonstrates N+1 query anti-pattern -- should be fixed
func (r *postgresItemRepository) GetItemsWithRelatedData(ctx context.Context, itemIDs []string) ([]*item.Item, error) {
	itemQuery := `SELECT i.id, i.sku, i.name, i.description, i.price_amount, i.price_currency, i.category_name,
    i.category_slug, i.inventory_quantity, i.images, i.attributes, i.status, i.created_at, i.updated_at,
    COALESCE(iv.view_count, 0) AS view_count,       -- Jumlah views, default 0 jika tidak ada
    COALESCE(ir.average_rating, 0.0) AS average_rating -- Rata-rata rating, default 0.0 jika tidak ada
   FROM items i
   LEFT JOIN (
    -- Subquery untuk menghitung jumlah views per item
    SELECT item_id, COUNT(*) AS view_count
    FROM item_views
    WHERE item_id = ANY($1) -- Pastikan hanya menghitung views untuk ID yang diminta
    GROUP BY item_id
   ) AS iv ON i.id = iv.item_id
   LEFT JOIN (
    -- Subquery untuk menghitung rata-rata rating per item
    SELECT item_id, AVG(rating) AS average_rating
    FROM item_ratings
    WHERE item_id = ANY($1)
    GROUP BY item_id
   ) AS ir ON i.id = ir.item_id
   WHERE i.id = ANY($1);`

	rows, err := r.db.QueryContext(ctx, itemQuery, itemIDs)
	if err != nil {
		log.Error().Interface("Error", err).Msg("failed to select items")
	}

	itemList, err := r.rowsToItems(rows)
	err = rows.Close()
	if err != nil {
		log.Error().Interface("Error", err).Msg("failed to close")
	}

	return itemList, nil
}
