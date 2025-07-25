package service

import "context"

type PricingService interface {
	CalculatePrice(ctx context.Context, basePrice float64, category string) (float64, error)
	ApplyDiscounts(ctx context.Context, price float64, itemID string) (float64, error)
}
