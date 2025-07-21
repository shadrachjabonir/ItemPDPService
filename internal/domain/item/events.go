package item

import (
	"time"

	"github.com/google/uuid"
)

// DomainEvent represents a domain event interface
type DomainEvent interface {
	EventID() string
	EventType() string
	AggregateID() string
	OccurredAt() time.Time
	EventData() interface{}
}

// BaseDomainEvent provides common event functionality
type BaseDomainEvent struct {
	eventID     string
	eventType   string
	aggregateID string
	occurredAt  time.Time
}

func NewBaseDomainEvent(eventType, aggregateID string) BaseDomainEvent {
	return BaseDomainEvent{
		eventID:     uuid.New().String(),
		eventType:   eventType,
		aggregateID: aggregateID,
		occurredAt:  time.Now(),
	}
}

func (e BaseDomainEvent) EventID() string     { return e.eventID }
func (e BaseDomainEvent) EventType() string   { return e.eventType }
func (e BaseDomainEvent) AggregateID() string { return e.aggregateID }
func (e BaseDomainEvent) OccurredAt() time.Time { return e.occurredAt }

// ItemCreatedEvent is raised when a new item is created
type ItemCreatedEvent struct {
	BaseDomainEvent
	Item *Item
}

func NewItemCreatedEvent(item *Item) *ItemCreatedEvent {
	return &ItemCreatedEvent{
		BaseDomainEvent: NewBaseDomainEvent("ItemCreated", item.ID().String()),
		Item:            item,
	}
}

func (e *ItemCreatedEvent) EventData() interface{} {
	return map[string]interface{}{
		"id":          e.Item.ID().String(),
		"sku":         e.Item.SKU().String(),
		"name":        e.Item.Name(),
		"price":       e.Item.Price().Amount(),
		"currency":    e.Item.Price().Currency(),
		"category":    e.Item.Category().Name(),
		"status":      e.Item.Status().String(),
	}
}

// ItemPriceChangedEvent is raised when item price is updated
type ItemPriceChangedEvent struct {
	BaseDomainEvent
	ItemID   ItemID
	OldPrice Price
	NewPrice Price
}

func NewItemPriceChangedEvent(itemID ItemID, oldPrice, newPrice Price) *ItemPriceChangedEvent {
	return &ItemPriceChangedEvent{
		BaseDomainEvent: NewBaseDomainEvent("ItemPriceChanged", itemID.String()),
		ItemID:          itemID,
		OldPrice:        oldPrice,
		NewPrice:        newPrice,
	}
}

func (e *ItemPriceChangedEvent) EventData() interface{} {
	return map[string]interface{}{
		"itemId":      e.ItemID.String(),
		"oldPrice":    e.OldPrice.Amount(),
		"newPrice":    e.NewPrice.Amount(),
		"currency":    e.NewPrice.Currency(),
		"changeAmount": e.NewPrice.Amount() - e.OldPrice.Amount(),
	}
}

// ItemInventoryUpdatedEvent is raised when item inventory is updated
type ItemInventoryUpdatedEvent struct {
	BaseDomainEvent
	ItemID      ItemID
	OldQuantity int
	NewQuantity int
}

func NewItemInventoryUpdatedEvent(itemID ItemID, oldQuantity, newQuantity int) *ItemInventoryUpdatedEvent {
	return &ItemInventoryUpdatedEvent{
		BaseDomainEvent: NewBaseDomainEvent("ItemInventoryUpdated", itemID.String()),
		ItemID:          itemID,
		OldQuantity:     oldQuantity,
		NewQuantity:     newQuantity,
	}
}

func (e *ItemInventoryUpdatedEvent) EventData() interface{} {
	return map[string]interface{}{
		"itemId":       e.ItemID.String(),
		"oldQuantity":  e.OldQuantity,
		"newQuantity":  e.NewQuantity,
		"changeAmount": e.NewQuantity - e.OldQuantity,
	}
}

// ItemStatusChangedEvent is raised when item status changes
type ItemStatusChangedEvent struct {
	BaseDomainEvent
	ItemID    ItemID
	OldStatus Status
	NewStatus Status
}

func NewItemStatusChangedEvent(itemID ItemID, oldStatus, newStatus Status) *ItemStatusChangedEvent {
	return &ItemStatusChangedEvent{
		BaseDomainEvent: NewBaseDomainEvent("ItemStatusChanged", itemID.String()),
		ItemID:          itemID,
		OldStatus:       oldStatus,
		NewStatus:       newStatus,
	}
}

func (e *ItemStatusChangedEvent) EventData() interface{} {
	return map[string]interface{}{
		"itemId":    e.ItemID.String(),
		"oldStatus": e.OldStatus.String(),
		"newStatus": e.NewStatus.String(),
	}
}

// ItemDeletedEvent is raised when an item is deleted
type ItemDeletedEvent struct {
	BaseDomainEvent
	ItemID ItemID
	SKU    SKU
}

func NewItemDeletedEvent(itemID ItemID, sku SKU) *ItemDeletedEvent {
	return &ItemDeletedEvent{
		BaseDomainEvent: NewBaseDomainEvent("ItemDeleted", itemID.String()),
		ItemID:          itemID,
		SKU:             sku,
	}
}

func (e *ItemDeletedEvent) EventData() interface{} {
	return map[string]interface{}{
		"itemId": e.ItemID.String(),
		"sku":    e.SKU.String(),
	}
} 