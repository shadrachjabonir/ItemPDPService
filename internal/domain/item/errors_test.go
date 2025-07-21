package item

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDomainError(t *testing.T) {
	message := "test error message"
	err := NewDomainError(message)
	
	assert.NotNil(t, err)
	assert.Equal(t, message, err.Error())
}

func TestDomainError_Is(t *testing.T) {
	err1 := NewDomainError("error 1")
	err2 := NewDomainError("error 2")
	standardErr := errors.New("standard error")
	
	// Domain errors should match other domain errors
	assert.True(t, err1.Is(err2))
	assert.True(t, err2.Is(err1))
	
	// Domain errors should not match standard errors
	assert.False(t, err1.Is(standardErr))
}

func TestPredefinedErrors(t *testing.T) {
	tests := []struct {
		name string
		err  *DomainError
	}{
		{"item not found", ErrItemNotFound},
		{"item already exists", ErrItemAlreadyExists},
		{"invalid SKU", ErrInvalidSKU},
		{"invalid price", ErrInvalidPrice},
		{"insufficient stock", ErrInsufficientStock},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNil(t, tt.err)
			assert.NotEmpty(t, tt.err.Error())
			assert.IsType(t, &DomainError{}, tt.err)
		})
	}
}

func TestItemNotFoundError(t *testing.T) {
	id := NewItemID()
	err := ItemNotFoundError(id)
	
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), id.String())
	assert.Contains(t, err.Error(), "not found")
}

func TestItemNotFoundBySKUError(t *testing.T) {
	sku, _ := NewSKU("TEST-001")
	err := ItemNotFoundBySKUError(sku)
	
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), sku.String())
	assert.Contains(t, err.Error(), "not found")
}

func TestDuplicateSKUError(t *testing.T) {
	sku, _ := NewSKU("TEST-001")
	err := DuplicateSKUError(sku)
	
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), sku.String())
	assert.Contains(t, err.Error(), "already exists")
} 