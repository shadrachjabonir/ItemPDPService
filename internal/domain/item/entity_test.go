package item

import (
	"testing"
	"time"
)

func TestNewItem(t *testing.T) {
	t.Run("valid item creation", func(t *testing.T) {
		sku, _ := NewSKU("TEST-001")
		price, _ := NewPrice(99.99, "USD")
		category, _ := NewCategory("Electronics")

		item, err := NewItem(sku, "Test Item", "Test Description", price, category)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if item == nil {
			t.Fatal("Expected item to be created")
		}

		if item.Name() != "Test Item" {
			t.Errorf("Expected name 'Test Item', got %s", item.Name())
		}

		if item.SKU().String() != "TEST-001" {
			t.Errorf("Expected SKU 'TEST-001', got %s", item.SKU().String())
		}
	})

	t.Run("empty name should fail", func(t *testing.T) {
		sku, _ := NewSKU("TEST-001")
		price, _ := NewPrice(99.99, "USD")
		category, _ := NewCategory("Electronics")

		item, err := NewItem(sku, "", "Test Description", price, category)

		if err == nil {
			t.Error("Expected error for empty name")
		}

		if item != nil {
			t.Error("Expected no item to be created")
		}
	})
}

func TestItem_SetPrice(t *testing.T) {
	sku, _ := NewSKU("TEST-001")
	price, _ := NewPrice(99.99, "USD")
	category, _ := NewCategory("Electronics")
	item, _ := NewItem(sku, "Test Item", "Test Description", price, category)

	newPrice, _ := NewPrice(149.99, "USD")
	item.SetPrice(newPrice)

	if item.Price().Amount() != 149.99 {
		t.Errorf("Expected price 149.99, got %f", item.Price().Amount())
	}
}

func TestItem_SetInventory(t *testing.T) {
	sku, _ := NewSKU("TEST-001")
	price, _ := NewPrice(99.99, "USD")
	category, _ := NewCategory("Electronics")
	item, _ := NewItem(sku, "Test Item", "Test Description", price, category)

	t.Run("valid positive quantity", func(t *testing.T) {
		newInventory, _ := NewInventory(50)
		item.SetInventory(newInventory)

		if item.Inventory().Quantity() != 50 {
			t.Errorf("Expected inventory 50, got %d", item.Inventory().Quantity())
		}
	})

	t.Run("valid zero quantity", func(t *testing.T) {
		newInventory, _ := NewInventory(0)
		item.SetInventory(newInventory)

		if item.Inventory().Quantity() != 0 {
			t.Errorf("Expected inventory 0, got %d", item.Inventory().Quantity())
		}
	})
}

func TestItem_AddImage(t *testing.T) {
	sku, _ := NewSKU("TEST-001")
	price, _ := NewPrice(99.99, "USD")
	category, _ := NewCategory("Electronics")
	item, _ := NewItem(sku, "Test Item", "Test Description", price, category)

	image, _ := NewImage("http://example.com/image.jpg", "Test Image", false)
	item.AddImage(image)

	if len(item.Images()) != 1 {
		t.Errorf("Expected 1 image, got %d", len(item.Images()))
	}
}

func TestItem_SetAttribute(t *testing.T) {
	sku, _ := NewSKU("TEST-001")
	price, _ := NewPrice(99.99, "USD")
	category, _ := NewCategory("Electronics")
	item, _ := NewItem(sku, "Test Item", "Test Description", price, category)

	t.Run("valid attribute", func(t *testing.T) {
		attrs := item.Attributes()
		attrs.Set("color", "red")

		value, exists := item.Attributes().Get("color")
		if !exists {
			t.Error("Expected attribute 'color' to exist")
		}
		if value != "red" {
			t.Errorf("Expected attribute value 'red', got %s", value)
		}
	})
}

func TestItem_StatusMethods(t *testing.T) {
	sku, _ := NewSKU("TEST-001")
	price, _ := NewPrice(99.99, "USD")
	category, _ := NewCategory("Electronics")
	item, _ := NewItem(sku, "Test Item", "Test Description", price, category)

	// Test setting status using setter
	item.SetStatus(StatusActive)
	if !item.IsActive() {
		t.Error("Expected item to be active")
	}

	item.SetStatus(StatusInactive)
	if !item.IsInactive() {
		t.Error("Expected item to be inactive")
	}

	item.SetStatus(StatusDraft)
	if !item.IsDraft() {
		t.Error("Expected item to be draft")
	}

	item.SetStatus(StatusArchived)
	if !item.IsArchived() {
		t.Error("Expected item to be archived")
	}
}

func TestItem_BasicProperties(t *testing.T) {
	sku, _ := NewSKU("TEST-001")
	price, _ := NewPrice(99.99, "USD")
	category, _ := NewCategory("Electronics")
	item, _ := NewItem(sku, "Test Item", "Test Description", price, category)

	// Test basic getters
	if item.Name() != "Test Item" {
		t.Errorf("Expected name 'Test Item', got %s", item.Name())
	}

	if item.Description() != "Test Description" {
		t.Errorf("Expected description 'Test Description', got %s", item.Description())
	}

	if item.Price().Amount() != 99.99 {
		t.Errorf("Expected price 99.99, got %f", item.Price().Amount())
	}

	if item.Category().Name() != "Electronics" {
		t.Errorf("Expected category 'Electronics', got %s", item.Category().Name())
	}

	// Test basic setters
	item.SetName("Updated Item")
	if item.Name() != "Updated Item" {
		t.Errorf("Expected name 'Updated Item', got %s", item.Name())
	}

	item.SetDescription("Updated Description")
	if item.Description() != "Updated Description" {
		t.Errorf("Expected description 'Updated Description', got %s", item.Description())
	}
}

func TestItem_Timestamps(t *testing.T) {
	sku, _ := NewSKU("TEST-001")
	price, _ := NewPrice(99.99, "USD")
	category, _ := NewCategory("Electronics")

	beforeCreate := time.Now()
	item, _ := NewItem(sku, "Test Item", "Test Description", price, category)
	afterCreate := time.Now()

	if item.CreatedAt().Before(beforeCreate) || item.CreatedAt().After(afterCreate) {
		t.Error("CreatedAt timestamp should be set during creation")
	}

	if item.UpdatedAt().Before(beforeCreate) || item.UpdatedAt().After(afterCreate) {
		t.Error("UpdatedAt timestamp should be set during creation")
	}

	// Test that UpdatedAt changes when modifying item
	originalUpdatedAt := item.UpdatedAt()
	time.Sleep(1 * time.Millisecond) // Ensure time difference

	item.SetName("New Name")

	if !item.UpdatedAt().After(originalUpdatedAt) {
		t.Error("UpdatedAt should be updated when item is modified")
	}
}
