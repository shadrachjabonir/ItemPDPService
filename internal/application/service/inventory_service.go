package service

import "context"

type InventoryService interface {
	ReserveInventory(ctx context.Context, itemID string, quantity int) error
	ReleaseInventory(ctx context.Context, itemID string, quantity int) error
}