package service

import "context"

type CategoryService interface {
	ValidateCategory(ctx context.Context, category string) error
	GetCategoryDiscounts(ctx context.Context, category string) (float64, error)
}
