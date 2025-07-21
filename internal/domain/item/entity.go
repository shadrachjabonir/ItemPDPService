package item

import (
	"time"
)

// Item represents a product item in the system
// Simple data container with basic getters and setters
type Item struct {
	id          ItemID
	sku         SKU
	name        string
	description string
	price       Price
	category    Category
	inventory   Inventory
	images      []Image
	attributes  Attributes
	status      Status
	createdAt   time.Time
	updatedAt   time.Time
}

// NewItem creates a new item with basic validation
func NewItem(sku SKU, name, description string, price Price, category Category) (*Item, error) {
	if name == "" {
		return nil, NewDomainError("item name cannot be empty")
	}

	item := &Item{
		id:          NewItemID(),
		sku:         sku,
		name:        name,
		description: description,
		price:       price,
		category:    category,
		inventory:   Inventory{quantity: 0},
		images:      make([]Image, 0),
		attributes:  NewAttributes(),
		status:      StatusDraft,
		createdAt:   time.Now(),
		updatedAt:   time.Now(),
	}

	return item, nil
}

// Basic getters - anemic model pattern
func (i *Item) ID() ItemID             { return i.id }
func (i *Item) SKU() SKU               { return i.sku }
func (i *Item) Name() string           { return i.name }
func (i *Item) Description() string    { return i.description }
func (i *Item) Price() Price           { return i.price }
func (i *Item) Category() Category     { return i.category }
func (i *Item) Inventory() Inventory   { return i.inventory }
func (i *Item) Images() []Image        { return i.images }
func (i *Item) Attributes() Attributes { return i.attributes }
func (i *Item) Status() Status         { return i.status }
func (i *Item) CreatedAt() time.Time   { return i.createdAt }
func (i *Item) UpdatedAt() time.Time   { return i.updatedAt }

// Basic setters - anemic model pattern
func (i *Item) SetName(name string)              { i.name = name; i.updatedAt = time.Now() }
func (i *Item) SetDescription(desc string)       { i.description = desc; i.updatedAt = time.Now() }
func (i *Item) SetPrice(price Price)             { i.price = price; i.updatedAt = time.Now() }
func (i *Item) SetCategory(category Category)    { i.category = category; i.updatedAt = time.Now() }
func (i *Item) SetInventory(inventory Inventory) { i.inventory = inventory; i.updatedAt = time.Now() }
func (i *Item) SetStatus(status Status)          { i.status = status; i.updatedAt = time.Now() }
func (i *Item) SetImages(images []Image)         { i.images = images; i.updatedAt = time.Now() }

// Simple utility methods without business logic
func (i *Item) AddImage(image Image) {
	i.images = append(i.images, image)
	i.updatedAt = time.Now()
}

func (i *Item) ClearImages() {
	i.images = make([]Image, 0)
	i.updatedAt = time.Now()
}

// Basic status checks without business rules
func (i *Item) IsActive() bool   { return i.status == StatusActive }
func (i *Item) IsDraft() bool    { return i.status == StatusDraft }
func (i *Item) IsInactive() bool { return i.status == StatusInactive }
func (i *Item) IsArchived() bool { return i.status == StatusArchived }
