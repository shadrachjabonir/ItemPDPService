package item

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewItemID(t *testing.T) {
	id := NewItemID()
	assert.NotEmpty(t, id.String())
	assert.Equal(t, 36, len(id.String())) // UUID format
}

func TestNewItemIDFromString(t *testing.T) {
	tests := []struct {
		name    string
		idStr   string
		wantErr bool
	}{
		{"valid UUID", "550e8400-e29b-41d4-a716-446655440000", false},
		{"invalid UUID", "invalid-uuid", true},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := NewItemIDFromString(tt.idStr)
			
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.idStr, id.String())
			}
		})
	}
}

func TestItemID_Equals(t *testing.T) {
	id1 := NewItemID()
	id2 := NewItemID()
	id3, _ := NewItemIDFromString(id1.String())

	assert.False(t, id1.Equals(id2))
	assert.True(t, id1.Equals(id3))
}

func TestNewSKU(t *testing.T) {
	tests := []struct {
		name    string
		sku     string
		wantErr bool
		errMsg  string
	}{
		{"valid SKU", "ABC-123", false, ""},
		{"lowercase converted", "abc-123", false, ""},
		{"with underscores", "ABC_123", false, ""},
		{"empty SKU", "", true, "SKU cannot be empty"},
		{"too short", "AB", true, "SKU must be between 3 and 20 characters"},
		{"too long", "ABCDEFGHIJKLMNOPQRSTUVWXYZ", true, "SKU must be between 3 and 20 characters"},
		{"invalid characters", "ABC@123", true, "SKU can only contain uppercase letters, numbers, hyphens, and underscores"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sku, err := NewSKU(tt.sku)
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
				// Should be uppercase
				assert.Equal(t, strings.ToUpper(strings.TrimSpace(tt.sku)), sku.String())
			}
		})
	}
}

func TestNewPrice(t *testing.T) {
	tests := []struct {
		name     string
		amount   float64
		currency string
		wantErr  bool
		errMsg   string
	}{
		{"valid price USD", 99.99, "USD", false, ""},
		{"valid price EUR", 85.50, "EUR", false, ""},
		{"zero price", 0.0, "USD", false, ""},
		{"empty currency defaults to USD", 50.00, "", false, ""},
		{"negative price", -10.0, "USD", true, "price cannot be negative"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			price, err := NewPrice(tt.amount, tt.currency)
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.amount, price.Amount())
				expectedCurrency := tt.currency
				if expectedCurrency == "" {
					expectedCurrency = "USD"
				}
				assert.Equal(t, strings.ToUpper(expectedCurrency), price.Currency())
			}
		})
	}
}

func TestPrice_String(t *testing.T) {
	price, _ := NewPrice(99.99, "USD")
	assert.Equal(t, "99.99 USD", price.String())
}

func TestNewCategory(t *testing.T) {
	tests := []struct {
		name         string
		categoryName string
		wantErr      bool
		expectedSlug string
	}{
		{"valid category", "Electronics", false, "electronics"},
		{"category with spaces", "Home & Garden", false, "home-&-garden"},
		{"category with trimming", "  Books  ", false, "books"},
		{"empty category", "", true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			category, err := NewCategory(tt.categoryName)
			
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, strings.TrimSpace(tt.categoryName), category.Name())
				assert.Equal(t, tt.expectedSlug, category.Slug())
			}
		})
	}
}

func TestNewInventory(t *testing.T) {
	tests := []struct {
		name     string
		quantity int
		wantErr  bool
	}{
		{"valid positive quantity", 100, false},
		{"valid zero quantity", 0, false},
		{"invalid negative quantity", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inventory, err := NewInventory(tt.quantity)
			
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.quantity, inventory.Quantity())
				assert.Equal(t, tt.quantity > 0, inventory.IsAvailable())
				assert.Equal(t, tt.quantity >= tt.quantity, inventory.CanReserve(tt.quantity))
			}
		})
	}
}

func TestInventory_CanReserve(t *testing.T) {
	inventory, _ := NewInventory(10)
	
	assert.True(t, inventory.CanReserve(5))
	assert.True(t, inventory.CanReserve(10))
	assert.False(t, inventory.CanReserve(15))
}

func TestNewImage(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		alt       string
		isPrimary bool
		wantErr   bool
	}{
		{"valid image", "https://example.com/image.jpg", "Test Image", true, false},
		{"empty URL", "", "Test Image", false, true},
		{"valid without alt", "https://example.com/image.jpg", "", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			image, err := NewImage(tt.url, tt.alt, tt.isPrimary)
			
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.url, image.URL())
				assert.Equal(t, tt.alt, image.Alt())
				assert.Equal(t, tt.isPrimary, image.IsPrimary())
			}
		})
	}
}

func TestAttributes(t *testing.T) {
	attributes := NewAttributes()
	
	// Test setting attributes
	err := attributes.Set("color", "red")
	assert.NoError(t, err)
	
	err = attributes.Set("size", "large")
	assert.NoError(t, err)
	
	// Test getting attributes
	color, exists := attributes.Get("color")
	assert.True(t, exists)
	assert.Equal(t, "red", color)
	
	_, exists = attributes.Get("nonexistent")
	assert.False(t, exists)
	
	// Test empty key
	err = attributes.Set("", "value")
	assert.Error(t, err)
	
	// Test All method
	all := attributes.All()
	assert.Len(t, all, 2)
	assert.Equal(t, "red", all["color"])
	assert.Equal(t, "large", all["size"])
}

func TestStatus(t *testing.T) {
	tests := []struct {
		status   Status
		expected string
	}{
		{StatusActive, "active"},
		{StatusInactive, "inactive"},
		{StatusDraft, "draft"},
		{StatusArchived, "archived"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.status.String())
		})
	}
}

func TestStatusFromString(t *testing.T) {
	tests := []struct {
		name     string
		statusStr string
		expected Status
		wantErr  bool
	}{
		{"active", "active", StatusActive, false},
		{"inactive", "inactive", StatusInactive, false},
		{"draft", "draft", StatusDraft, false},
		{"archived", "archived", StatusArchived, false},
		{"uppercase", "ACTIVE", StatusActive, false},
		{"invalid", "invalid", StatusActive, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, err := StatusFromString(tt.statusStr)
			
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, status)
			}
		})
	}
} 