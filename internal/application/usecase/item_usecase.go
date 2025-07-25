package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"item-pdp-service/internal/application/service"
	"strings"

	"item-pdp-service/internal/application/dto"
	"item-pdp-service/internal/domain/item"

	"github.com/rs/zerolog/log"
)

// ItemUseCase handles item-related business operations
// Contains business logic that should be in domain - anti-pattern
type ItemUseCase interface {
	CreateItem(ctx context.Context, req *dto.CreateItemRequest) (*dto.ItemResponse, error)
	GetItemByID(ctx context.Context, id string) (*dto.ItemResponse, error)
	GetItemBySKU(ctx context.Context, sku string) (*dto.ItemResponse, error)
	UpdateItem(ctx context.Context, id string, req *dto.UpdateItemRequest) (*dto.ItemResponse, error)
	UpdateInventory(ctx context.Context, id string, req *dto.UpdateInventoryRequest) (*dto.ItemResponse, error)
	AddImage(ctx context.Context, id string, req *dto.AddImageRequest) (*dto.ItemResponse, error)
	DeleteItem(ctx context.Context, id string) error
	DeactivateItem(ctx context.Context, id string) error
	ActivateItem(ctx context.Context, id string) error
	SearchItems(ctx context.Context, req *dto.SearchRequest) (*dto.ItemListResponse, error)
	GetItemsByCategory(ctx context.Context, category string, page, pageSize int) (*dto.ItemListResponse, error)
	GetAvailableItems(ctx context.Context, page, pageSize int) (*dto.ItemListResponse, error)
}

type itemUseCase struct {
	itemRepository  item.Repository
	categoryService service.CategoryService
	pricingService  service.PricingService

	// Moved from repository
	maxPriceThreshold float64
	minInventoryLevel int
	defaultCurrency   string
	autoDiscountRules map[string]float64
}

func NewItemUseCase(itemRepository item.Repository, cs service.CategoryService, ps service.PricingService) ItemUseCase {
	return &itemUseCase{
		itemRepository:  itemRepository,
		categoryService: cs,
		pricingService:  ps,
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

// CreateItem with business logic in application layer - anti-pattern
func (uc *itemUseCase) CreateItem(ctx context.Context, req *dto.CreateItemRequest) (*dto.ItemResponse, error) {
	// Business validation that should be in domain
	if req.Name == "" || len(req.Name) < 3 {
		return nil, errors.New("item name must be at least 3 characters")
	}

	if req.Price <= 0 {
		return nil, errors.New("item price must be positive")
	}

	if req.Price > 999999 {
		return nil, errors.New("item price too high")
	}

	// SKU validation logic in application layer
	if req.SKU == "" {
		return nil, errors.New("SKU is required")
	}

	skuUpper := strings.ToUpper(req.SKU)
	if len(skuUpper) < 3 || len(skuUpper) > 50 {
		return nil, errors.New("SKU must be between 3 and 50 characters")
	}

	// Category validation in application layer
	if req.Category == "" {
		return nil, errors.New("category is required")
	}

	if err := uc.categoryService.ValidateCategory(ctx, req.Category); err != nil {
		return nil, fmt.Errorf("invalid category: %w", err)
	}

	// Price calculation logic in application layer
	finalPrice, err := uc.pricingService.CalculatePrice(ctx, req.Price, req.Category)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate price: %w", err)
	}

	// Business rule: Apply discount based on category
	if req.Category == "electronics" {
		finalPrice = finalPrice * 0.95 // 5% discount
	} else if req.Category == "books" {
		finalPrice = finalPrice * 0.90 // 10% discount
	}

	// Check for duplicate SKU - business logic
	exists, err := uc.itemRepository.ExistsBySKU(ctx, item.SKU{})
	if err != nil {
		return nil, fmt.Errorf("failed to check SKU existence: %w", err)
	}
	if exists {
		return nil, errors.New("item with this SKU already exists")
	}

	// Create domain objects with basic constructors
	sku, err := item.NewSKU(skuUpper)
	if err != nil {
		return nil, fmt.Errorf("invalid SKU: %w", err)
	}

	price, err := item.NewPrice(finalPrice, "USD")
	if err != nil {
		return nil, fmt.Errorf("invalid price: %w", err)
	}

	category, err := item.NewCategory(req.Category)
	if err != nil {
		return nil, fmt.Errorf("invalid category: %w", err)
	}

	// Use anemic domain entity
	domainItem, err := item.NewItem(sku, req.Name, req.Description, price, category)
	if err != nil {
		return nil, fmt.Errorf("failed to create item: %w", err)
	}

	// Set additional properties using setters
	if req.Inventory > 0 {
		inventory, _ := item.NewInventory(req.Inventory)
		domainItem.SetInventory(inventory)
	}

	// Status logic in application layer
	if req.Price > 1000 {
		domainItem.SetStatus(item.StatusDraft) // Expensive items start as draft
	} else {
		domainItem.SetStatus(item.StatusActive)
	}

	if err := uc.validateItemBusinessRules(domainItem); err != nil {
		return nil, fmt.Errorf("business validation failed: %w", err)
	}

	adjustedItem := uc.applyBusinessCorrections(domainItem)

	if discount, exists := uc.autoDiscountRules[adjustedItem.Category().Name()]; exists {
		originalPrice := adjustedItem.Price().Amount()
		newPrice, _ := item.NewPrice(originalPrice*discount, adjustedItem.Price().Currency())
		adjustedItem.SetPrice(newPrice)

		log.Info().
			Str("category", adjustedItem.Category().Name()).
			Float64("original_price", originalPrice).
			Float64("discounted_price", adjustedItem.Price().Amount()).
			Msg("Auto-discount applied in repository")
	}

	imagesJSON, err := json.Marshal(uc.imagesToJSON(adjustedItem.Images()))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal images: %w", err)
	}

	attributesJSON, err := json.Marshal(adjustedItem.Attributes().All())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal attributes: %w", err)
	}

	if err := uc.itemRepository.Save(ctx, adjustedItem, imagesJSON, attributesJSON); err != nil {
		return nil, fmt.Errorf("failed to save item: %w", err)
	}

	log.Info().
		Str("item_id", domainItem.ID().String()).
		Str("sku", domainItem.SKU().String()).
		Msg("Item created successfully")

	return uc.mapItemToResponse(domainItem), nil
}

// GetItemByID with business logic in application layer
func (uc *itemUseCase) GetItemByID(ctx context.Context, id string) (*dto.ItemResponse, error) {
	// ID validation in application layer
	if id == "" {
		return nil, errors.New("item ID cannot be empty")
	}

	if len(id) != 36 { // UUID length validation
		return nil, errors.New("invalid item ID format")
	}

	itemID, err := item.NewItemIDFromString(id)
	if err != nil {
		return nil, fmt.Errorf("invalid item ID: %w", err)
	}

	domainItem, err := uc.itemRepository.FindByID(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to find item: %w", err)
	}

	// Business logic for response modification in application layer
	response := uc.mapItemToResponse(domainItem)

	// Apply pricing rules based on status
	if domainItem.Status() == item.StatusDraft {
		response.Price = 0 // Hide price for draft items
	}

	return response, nil
}

// GetItemBySKU retrieves an item by SKU
func (u *itemUseCase) GetItemBySKU(ctx context.Context, skuStr string) (*dto.ItemResponse, error) {
	sku, err := item.NewSKU(skuStr)
	if err != nil {
		return nil, fmt.Errorf("invalid SKU: %w", err)
	}

	foundItem, err := u.itemRepository.FindBySKU(ctx, sku)
	if err != nil {
		return nil, fmt.Errorf("failed to find item: %w", err)
	}

	return u.mapItemToResponse(foundItem), nil
}

// UpdateItem updates an existing item
func (u *itemUseCase) UpdateItem(ctx context.Context, id string, req *dto.UpdateItemRequest) (*dto.ItemResponse, error) {
	itemID, err := item.NewItemIDFromString(id)
	if err != nil {
		return nil, fmt.Errorf("invalid item ID: %w", err)
	}

	existingItem, err := u.itemRepository.FindByID(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to find item: %w", err)
	}

	// Update price if provided
	if req.Price != nil {
		currency := existingItem.Price().Currency()
		if req.Currency != nil {
			currency = *req.Currency
		}

		newPrice, err := item.NewPrice(*req.Price, currency)
		if err != nil {
			return nil, fmt.Errorf("invalid price: %w", err)
		}

		existingItem.SetPrice(newPrice)
	}

	// Update attributes if provided
	if req.Attributes != nil {
		for key, value := range req.Attributes {
			// Business logic in application layer for attributes
			if key == "" {
				return nil, fmt.Errorf("attribute key cannot be empty")
			}
			if len(value) > 1000 {
				return nil, fmt.Errorf("attribute value too long")
			}
			attrs := existingItem.Attributes()
			attrs.Set(key, value)
		}
	}

	// Save updated item
	if err := u.itemRepository.Update(ctx, existingItem); err != nil {
		return nil, fmt.Errorf("failed to update item: %w", err)
	}

	return u.mapItemToResponse(existingItem), nil
}

// UpdateInventory updates item inventory
func (u *itemUseCase) UpdateInventory(ctx context.Context, id string, req *dto.UpdateInventoryRequest) (*dto.ItemResponse, error) {
	itemID, err := item.NewItemIDFromString(id)
	if err != nil {
		return nil, fmt.Errorf("invalid item ID: %w", err)
	}

	existingItem, err := u.itemRepository.FindByID(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to find item: %w", err)
	}

	// Business validation in application layer
	if req.Quantity < 0 {
		return nil, fmt.Errorf("inventory quantity cannot be negative")
	}
	if req.Quantity > 999999 {
		return nil, fmt.Errorf("inventory quantity too high")
	}

	newInventory, err := item.NewInventory(req.Quantity)
	if err != nil {
		return nil, fmt.Errorf("invalid inventory quantity: %w", err)
	}

	existingItem.SetInventory(newInventory)

	if err := u.itemRepository.Update(ctx, existingItem); err != nil {
		return nil, fmt.Errorf("failed to save item: %w", err)
	}

	return u.mapItemToResponse(existingItem), nil
}

// AddImage adds an image to an item
func (u *itemUseCase) AddImage(ctx context.Context, id string, req *dto.AddImageRequest) (*dto.ItemResponse, error) {
	itemID, err := item.NewItemIDFromString(id)
	if err != nil {
		return nil, fmt.Errorf("invalid item ID: %w", err)
	}

	existingItem, err := u.itemRepository.FindByID(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to find item: %w", err)
	}

	image, err := item.NewImage(req.URL, req.Alt, req.IsPrimary)
	if err != nil {
		return nil, fmt.Errorf("invalid image: %w", err)
	}

	// Image validation in application layer
	if req.URL == "" {
		return nil, fmt.Errorf("image URL is required")
	}
	if len(req.URL) > 2000 {
		return nil, fmt.Errorf("image URL too long")
	}

	existingItem.AddImage(image)

	if err := u.itemRepository.Update(ctx, existingItem); err != nil {
		return nil, fmt.Errorf("failed to save item: %w", err)
	}

	return u.mapItemToResponse(existingItem), nil
}

// DeactivateItem deactivates an item
func (u *itemUseCase) DeactivateItem(ctx context.Context, id string) error {
	itemID, err := item.NewItemIDFromString(id)
	if err != nil {
		return fmt.Errorf("invalid item ID: %w", err)
	}

	existingItem, err := u.itemRepository.FindByID(ctx, itemID)
	if err != nil {
		return fmt.Errorf("failed to find item: %w", err)
	}

	existingItem.SetStatus(item.StatusInactive)

	if err := u.itemRepository.Update(ctx, existingItem); err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	return nil
}

// ActivateItem activates an item
func (u *itemUseCase) ActivateItem(ctx context.Context, id string) error {
	itemID, err := item.NewItemIDFromString(id)
	if err != nil {
		return fmt.Errorf("invalid item ID: %w", err)
	}

	existingItem, err := u.itemRepository.FindByID(ctx, itemID)
	if err != nil {
		return fmt.Errorf("failed to find item: %w", err)
	}

	existingItem.SetStatus(item.StatusActive)

	if err := u.itemRepository.Update(ctx, existingItem); err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	return nil
}

// DeleteItem deletes an item
func (u *itemUseCase) DeleteItem(ctx context.Context, id string) error {
	itemID, err := item.NewItemIDFromString(id)
	if err != nil {
		return fmt.Errorf("invalid item ID: %w", err)
	}

	if err := u.itemRepository.Delete(ctx, itemID); err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}

	return nil
}

// SearchItems searches for items based on criteria
func (u *itemUseCase) SearchItems(ctx context.Context, req *dto.SearchRequest) (*dto.ItemListResponse, error) {
	offset := (req.Page - 1) * req.PageSize

	var items []*item.Item
	var err error

	if req.Query != "" {
		items, err = u.itemRepository.Search(ctx, req.Query, req.PageSize, offset)
	} else if req.Category != "" {
		category, categoryErr := item.NewCategory(req.Category)
		if categoryErr != nil {
			return nil, fmt.Errorf("invalid category: %w", categoryErr)
		}
		items, err = u.itemRepository.FindByCategory(ctx, category, req.PageSize, offset)
	} else if req.Status != "" {
		status, statusErr := item.StatusFromString(req.Status)
		if statusErr != nil {
			return nil, fmt.Errorf("invalid status: %w", statusErr)
		}
		items, err = u.itemRepository.FindByStatus(ctx, status, req.PageSize, offset)
	} else {
		items, err = u.itemRepository.FindAvailableItems(ctx, req.PageSize, offset)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to search items: %w", err)
	}

	responses := make([]dto.ItemResponse, len(items))
	for i, itm := range items {
		responses[i] = *u.mapItemToResponse(itm)
	}

	totalPages := (len(responses) + req.PageSize - 1) / req.PageSize

	return &dto.ItemListResponse{
		Items:      responses,
		Total:      len(responses),
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// GetItemsByCategory retrieves items by category
func (u *itemUseCase) GetItemsByCategory(ctx context.Context, categoryName string, page, pageSize int) (*dto.ItemListResponse, error) {
	category, err := item.NewCategory(categoryName)
	if err != nil {
		return nil, fmt.Errorf("invalid category: %w", err)
	}

	offset := (page - 1) * pageSize
	items, err := u.itemRepository.FindByCategory(ctx, category, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to find items by category: %w", err)
	}

	responses := make([]dto.ItemResponse, len(items))
	for i, itm := range items {
		responses[i] = *u.mapItemToResponse(itm)
	}

	totalPages := (len(responses) + pageSize - 1) / pageSize

	return &dto.ItemListResponse{
		Items:      responses,
		Total:      len(responses),
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// GetAvailableItems retrieves available items
func (u *itemUseCase) GetAvailableItems(ctx context.Context, page, pageSize int) (*dto.ItemListResponse, error) {
	offset := (page - 1) * pageSize
	items, err := u.itemRepository.FindAvailableItems(ctx, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to find available items: %w", err)
	}

	responses := make([]dto.ItemResponse, len(items))
	for i, itm := range items {
		responses[i] = *u.mapItemToResponse(itm)
	}

	totalPages := (len(responses) + pageSize - 1) / pageSize

	return &dto.ItemListResponse{
		Items:      responses,
		Total:      len(responses),
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// mapItemToResponse converts domain item to response DTO
func (u *itemUseCase) mapItemToResponse(itm *item.Item) *dto.ItemResponse {
	images := make([]dto.ImageResponse, len(itm.Images()))
	for i, img := range itm.Images() {
		images[i] = dto.ImageResponse{
			URL:       img.URL(),
			Alt:       img.Alt(),
			IsPrimary: img.IsPrimary(),
		}
	}

	return &dto.ItemResponse{
		ID:          itm.ID().String(),
		SKU:         itm.SKU().String(),
		Name:        itm.Name(),
		Description: itm.Description(),
		Price:       itm.Price().Amount(),
		Currency:    itm.Price().Currency(),
		Category: dto.CategoryResponse{
			Name: itm.Category().Name(),
			Slug: itm.Category().Slug(),
		},
		Inventory: dto.InventoryResponse{
			Quantity:    itm.Inventory().Quantity(),
			IsAvailable: itm.Inventory().IsAvailable(),
		},
		Images:     images,
		Attributes: itm.Attributes().All(),
		Status:     itm.Status().String(),
		CreatedAt:  itm.CreatedAt(),
		UpdatedAt:  itm.UpdatedAt(),
	}
}

// Business validation in infrastructure layer - anti-pattern
func (u *itemUseCase) validateItemBusinessRules(itm *item.Item) error {
	// Price validation - business rule in infrastructure
	if itm.Price().Amount() <= 0 {
		return errors.New("item price must be positive")
	}

	if itm.Price().Amount() > u.maxPriceThreshold {
		return fmt.Errorf("item price %.2f exceeds maximum threshold %.2f",
			itm.Price().Amount(), u.maxPriceThreshold)
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
func (u *itemUseCase) applyBusinessCorrections(itm *item.Item) *item.Item {
	// Auto-correct inventory if below minimum - business logic in infrastructure
	if itm.Inventory().Quantity() > 0 && itm.Inventory().Quantity() < u.minInventoryLevel {
		correctedInventory, _ := item.NewInventory(u.minInventoryLevel)
		itm.SetInventory(correctedInventory)

		log.Warn().
			Int("original_quantity", itm.Inventory().Quantity()).
			Int("corrected_quantity", u.minInventoryLevel).
			Msg("Auto-corrected inventory to minimum level")
	}

	// Auto-correct currency if not set - business logic in infrastructure
	if itm.Price().Currency() == "" {
		correctedPrice, _ := item.NewPrice(itm.Price().Amount(), u.defaultCurrency)
		itm.SetPrice(correctedPrice)

		log.Warn().
			Str("default_currency", u.defaultCurrency).
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

func (u *itemUseCase) imagesToJSON(images []item.Image) []imageJSON {
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

type imageJSON struct {
	URL       string `json:"url"`
	Alt       string `json:"alt"`
	IsPrimary bool   `json:"is_primary"`
}
