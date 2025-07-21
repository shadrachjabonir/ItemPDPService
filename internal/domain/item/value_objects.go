package item

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

// ItemID is a value object representing an item identifier
type ItemID struct {
	value string
}

func NewItemID() ItemID {
	return ItemID{value: uuid.New().String()}
}

func NewItemIDFromString(id string) (ItemID, error) {
	if _, err := uuid.Parse(id); err != nil {
		return ItemID{}, NewDomainError("invalid item ID format")
	}
	return ItemID{value: id}, nil
}

func (id ItemID) String() string {
	return id.value
}

func (id ItemID) Equals(other ItemID) bool {
	return id.value == other.value
}

// SKU is a value object representing a Stock Keeping Unit
type SKU struct {
	value string
}

func NewSKU(sku string) (SKU, error) {
	sku = strings.TrimSpace(strings.ToUpper(sku))
	if err := validateSKU(sku); err != nil {
		return SKU{}, err
	}
	return SKU{value: sku}, nil
}

func (s SKU) String() string {
	return s.value
}

func (s SKU) Validate() error {
	return validateSKU(s.value)
}

func validateSKU(sku string) error {
	if sku == "" {
		return NewDomainError("SKU cannot be empty")
	}
	if len(sku) < 3 || len(sku) > 20 {
		return NewDomainError("SKU must be between 3 and 20 characters")
	}
	matched, _ := regexp.MatchString("^[A-Z0-9-_]+$", sku)
	if !matched {
		return NewDomainError("SKU can only contain uppercase letters, numbers, hyphens, and underscores")
	}
	return nil
}

// Price is a value object representing monetary value
type Price struct {
	amount   int64 // stored in cents to avoid floating point issues
	currency string
}

func NewPrice(amount float64, currency string) (Price, error) {
	if amount < 0 {
		return Price{}, NewDomainError("price cannot be negative")
	}
	if currency == "" {
		currency = "USD"
	}
	currency = strings.ToUpper(currency)
	
	// Convert to cents
	amountInCents := int64(amount * 100)
	
	return Price{
		amount:   amountInCents,
		currency: currency,
	}, nil
}

func (p Price) Amount() float64 {
	return float64(p.amount) / 100
}

func (p Price) Currency() string {
	return p.currency
}

func (p Price) String() string {
	return fmt.Sprintf("%.2f %s", p.Amount(), p.currency)
}

func (p Price) Validate() error {
	if p.amount < 0 {
		return NewDomainError("price cannot be negative")
	}
	if p.currency == "" {
		return NewDomainError("currency cannot be empty")
	}
	return nil
}

// Category is a value object representing item category
type Category struct {
	name string
	slug string
}

func NewCategory(name string) (Category, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Category{}, NewDomainError("category name cannot be empty")
	}
	
	slug := strings.ToLower(strings.ReplaceAll(name, " ", "-"))
	
	return Category{
		name: name,
		slug: slug,
	}, nil
}

func (c Category) Name() string {
	return c.name
}

func (c Category) Slug() string {
	return c.slug
}

func (c Category) Validate() error {
	if c.name == "" {
		return NewDomainError("category name cannot be empty")
	}
	return nil
}

// Inventory is a value object representing stock quantity
type Inventory struct {
	quantity int
}

func NewInventory(quantity int) (Inventory, error) {
	if quantity < 0 {
		return Inventory{}, NewDomainError("inventory quantity cannot be negative")
	}
	return Inventory{quantity: quantity}, nil
}

func (i Inventory) Quantity() int {
	return i.quantity
}

func (i Inventory) IsAvailable() bool {
	return i.quantity > 0
}

func (i Inventory) CanReserve(quantity int) bool {
	return i.quantity >= quantity
}

// Image is a value object representing an item image
type Image struct {
	url         string
	alt         string
	isPrimary   bool
}

func NewImage(url, alt string, isPrimary bool) (Image, error) {
	if url == "" {
		return Image{}, NewDomainError("image URL cannot be empty")
	}
	return Image{
		url:       url,
		alt:       alt,
		isPrimary: isPrimary,
	}, nil
}

func (i Image) URL() string {
	return i.url
}

func (i Image) Alt() string {
	return i.alt
}

func (i Image) IsPrimary() bool {
	return i.isPrimary
}

func (i Image) Validate() error {
	if i.url == "" {
		return NewDomainError("image URL cannot be empty")
	}
	return nil
}

// Attributes is a value object representing item attributes
type Attributes struct {
	data map[string]string
}

func NewAttributes() Attributes {
	return Attributes{
		data: make(map[string]string),
	}
}

func (a *Attributes) Set(key, value string) error {
	key = strings.TrimSpace(key)
	value = strings.TrimSpace(value)
	
	if key == "" {
		return NewDomainError("attribute key cannot be empty")
	}
	
	a.data[key] = value
	return nil
}

func (a Attributes) Get(key string) (string, bool) {
	value, exists := a.data[key]
	return value, exists
}

func (a Attributes) All() map[string]string {
	result := make(map[string]string)
	for k, v := range a.data {
		result[k] = v
	}
	return result
}

// Status is a value object representing item status
type Status int

const (
	StatusActive Status = iota
	StatusInactive
	StatusDraft
	StatusArchived
)

func (s Status) String() string {
	switch s {
	case StatusActive:
		return "active"
	case StatusInactive:
		return "inactive"
	case StatusDraft:
		return "draft"
	case StatusArchived:
		return "archived"
	default:
		return "unknown"
	}
}

func StatusFromString(status string) (Status, error) {
	switch strings.ToLower(status) {
	case "active":
		return StatusActive, nil
	case "inactive":
		return StatusInactive, nil
	case "draft":
		return StatusDraft, nil
	case "archived":
		return StatusArchived, nil
	default:
		return StatusActive, NewDomainError("invalid status: " + status)
	}
} 