package item

import "fmt"

// DomainError represents an error in the item domain
type DomainError struct {
	message string
}

func NewDomainError(message string) *DomainError {
	return &DomainError{message: message}
}

func (e *DomainError) Error() string {
	return e.message
}

func (e *DomainError) Is(target error) bool {
	_, ok := target.(*DomainError)
	return ok
}

// Specific domain errors
var (
	ErrItemNotFound     = &DomainError{message: "item not found"}
	ErrItemAlreadyExists = &DomainError{message: "item already exists"}
	ErrInvalidSKU       = &DomainError{message: "invalid SKU format"}
	ErrInvalidPrice     = &DomainError{message: "invalid price"}
	ErrInsufficientStock = &DomainError{message: "insufficient stock"}
)

// ItemNotFoundError creates a specific error for item not found by ID
func ItemNotFoundError(id ItemID) error {
	return &DomainError{message: fmt.Sprintf("item with ID %s not found", id.String())}
}

// ItemNotFoundBySKUError creates a specific error for item not found by SKU
func ItemNotFoundBySKUError(sku SKU) error {
	return &DomainError{message: fmt.Sprintf("item with SKU %s not found", sku.String())}
}

// DuplicateSKUError creates a specific error for duplicate SKU
func DuplicateSKUError(sku SKU) error {
	return &DomainError{message: fmt.Sprintf("item with SKU %s already exists", sku.String())}
} 