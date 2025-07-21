package usecase

import (
	"context"
	"testing"

	"item-pdp-service/internal/application/dto"
	"item-pdp-service/internal/domain/item"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock implementations for the new service dependencies
type MockInventoryService struct {
	mock.Mock
}

func (m *MockInventoryService) ReserveInventory(ctx context.Context, itemID string, quantity int) error {
	args := m.Called(ctx, itemID, quantity)
	return args.Error(0)
}

func (m *MockInventoryService) ReleaseInventory(ctx context.Context, itemID string, quantity int) error {
	args := m.Called(ctx, itemID, quantity)
	return args.Error(0)
}

type MockCategoryService struct {
	mock.Mock
}

func (m *MockCategoryService) ValidateCategory(ctx context.Context, category string) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockCategoryService) GetCategoryDiscounts(ctx context.Context, category string) (float64, error) {
	args := m.Called(ctx, category)
	return args.Get(0).(float64), args.Error(1)
}

type MockPricingService struct {
	mock.Mock
}

func (m *MockPricingService) CalculatePrice(ctx context.Context, basePrice float64, category string) (float64, error) {
	args := m.Called(ctx, basePrice, category)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockPricingService) ApplyDiscounts(ctx context.Context, price float64, itemID string) (float64, error) {
	args := m.Called(ctx, price, itemID)
	return args.Get(0).(float64), args.Error(1)
}

// MockItemRepository implementation
type MockItemRepository struct {
	mock.Mock
}

func (m *MockItemRepository) Save(ctx context.Context, itm *item.Item) error {
	args := m.Called(ctx, itm)
	return args.Error(0)
}

func (m *MockItemRepository) FindByID(ctx context.Context, id item.ItemID) (*item.Item, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*item.Item), args.Error(1)
}

func (m *MockItemRepository) FindBySKU(ctx context.Context, sku item.SKU) (*item.Item, error) {
	args := m.Called(ctx, sku)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*item.Item), args.Error(1)
}

func (m *MockItemRepository) Update(ctx context.Context, itm *item.Item) error {
	args := m.Called(ctx, itm)
	return args.Error(0)
}

func (m *MockItemRepository) Delete(ctx context.Context, id item.ItemID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockItemRepository) FindByCategory(ctx context.Context, category item.Category, limit, offset int) ([]*item.Item, error) {
	args := m.Called(ctx, category, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*item.Item), args.Error(1)
}

func (m *MockItemRepository) FindByStatus(ctx context.Context, status item.Status, limit, offset int) ([]*item.Item, error) {
	args := m.Called(ctx, status, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*item.Item), args.Error(1)
}

func (m *MockItemRepository) Search(ctx context.Context, query string, limit, offset int) ([]*item.Item, error) {
	args := m.Called(ctx, query, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*item.Item), args.Error(1)
}

func (m *MockItemRepository) FindAvailableItems(ctx context.Context, limit, offset int) ([]*item.Item, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*item.Item), args.Error(1)
}

func (m *MockItemRepository) ExistsBySKU(ctx context.Context, sku item.SKU) (bool, error) {
	args := m.Called(ctx, sku)
	return args.Bool(0), args.Error(1)
}

// Additional repository methods to satisfy interface
func (m *MockItemRepository) ExistsByID(ctx context.Context, id item.ItemID) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockItemRepository) CountByCategory(ctx context.Context, category item.Category) (int, error) {
	args := m.Called(ctx, category)
	return args.Int(0), args.Error(1)
}

func (m *MockItemRepository) CountByStatus(ctx context.Context, status item.Status) (int, error) {
	args := m.Called(ctx, status)
	return args.Int(0), args.Error(1)
}

func (m *MockItemRepository) FindItemsWithLowStock(ctx context.Context, threshold int) ([]*item.Item, error) {
	args := m.Called(ctx, threshold)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*item.Item), args.Error(1)
}

func TestItemUseCase_CreateItem(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		mockRepo := &MockItemRepository{}
		mockInventory := &MockInventoryService{}
		mockCategory := &MockCategoryService{}
		mockPricing := &MockPricingService{}

		useCase := NewItemUseCase(mockRepo, mockInventory, mockCategory, mockPricing)

		req := &dto.CreateItemRequest{
			SKU:         "TEST-001",
			Name:        "Test Item",
			Description: "Test Description",
			Price:       99.99,
			Category:    "electronics",
			Inventory:   10,
		}

		// Setup mock expectations for fat application service
		mockCategory.On("ValidateCategory", mock.Anything, "electronics").Return(nil)
		mockPricing.On("CalculatePrice", mock.Anything, 99.99, "electronics").Return(99.99, nil)
		mockRepo.On("ExistsBySKU", mock.Anything, mock.AnythingOfType("item.SKU")).Return(false, nil)
		mockRepo.On("Save", mock.Anything, mock.AnythingOfType("*item.Item")).Return(nil)

		result, err := useCase.CreateItem(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Test Item", result.Name)
		mockRepo.AssertExpectations(t)
		mockCategory.AssertExpectations(t)
		mockPricing.AssertExpectations(t)
	})

	t.Run("duplicate SKU error", func(t *testing.T) {
		mockRepo := &MockItemRepository{}
		mockInventory := &MockInventoryService{}
		mockCategory := &MockCategoryService{}
		mockPricing := &MockPricingService{}

		useCase := NewItemUseCase(mockRepo, mockInventory, mockCategory, mockPricing)

		req := &dto.CreateItemRequest{
			SKU:         "TEST-001",
			Name:        "Test Item",
			Description: "Test Description",
			Price:       99.99,
			Category:    "electronics",
		}

		// Setup expectations for business validation in application layer
		mockCategory.On("ValidateCategory", mock.Anything, "electronics").Return(nil)
		mockPricing.On("CalculatePrice", mock.Anything, 99.99, "electronics").Return(99.99, nil)
		mockRepo.On("ExistsBySKU", mock.Anything, mock.AnythingOfType("item.SKU")).Return(true, nil)

		result, err := useCase.CreateItem(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "already exists")
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid SKU", func(t *testing.T) {
		mockRepo := &MockItemRepository{}
		mockInventory := &MockInventoryService{}
		mockCategory := &MockCategoryService{}
		mockPricing := &MockPricingService{}

		useCase := NewItemUseCase(mockRepo, mockInventory, mockCategory, mockPricing)

		req := &dto.CreateItemRequest{
			SKU:         "", // Invalid empty SKU
			Name:        "Test Item",
			Description: "Test Description",
			Price:       99.99,
			Category:    "electronics",
		}

		result, err := useCase.CreateItem(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "SKU is required")
	})
}

func TestItemUseCase_GetItemByID(t *testing.T) {
	t.Run("successful retrieval", func(t *testing.T) {
		mockRepo := &MockItemRepository{}
		mockInventory := &MockInventoryService{}
		mockCategory := &MockCategoryService{}
		mockPricing := &MockPricingService{}

		useCase := NewItemUseCase(mockRepo, mockInventory, mockCategory, mockPricing)

		// Create test item
		testItem := createTestItem(t)
		itemID := testItem.ID()

		mockRepo.On("FindByID", mock.Anything, itemID).Return(testItem, nil)

		result, err := useCase.GetItemByID(context.Background(), itemID.String())

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, testItem.Name(), result.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("item not found", func(t *testing.T) {
		mockRepo := &MockItemRepository{}
		mockInventory := &MockInventoryService{}
		mockCategory := &MockCategoryService{}
		mockPricing := &MockPricingService{}

		useCase := NewItemUseCase(mockRepo, mockInventory, mockCategory, mockPricing)

		itemID := item.NewItemID()
		mockRepo.On("FindByID", mock.Anything, itemID).Return(nil, item.ItemNotFoundError(itemID))

		result, err := useCase.GetItemByID(context.Background(), itemID.String())

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid ID format", func(t *testing.T) {
		mockRepo := &MockItemRepository{}
		mockInventory := &MockInventoryService{}
		mockCategory := &MockCategoryService{}
		mockPricing := &MockPricingService{}

		useCase := NewItemUseCase(mockRepo, mockInventory, mockCategory, mockPricing)

		result, err := useCase.GetItemByID(context.Background(), "invalid-id")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid item ID format")
	})
}

func TestItemUseCase_UpdateInventory(t *testing.T) {
	t.Run("successful update", func(t *testing.T) {
		mockRepo := &MockItemRepository{}
		mockInventory := &MockInventoryService{}
		mockCategory := &MockCategoryService{}
		mockPricing := &MockPricingService{}

		useCase := NewItemUseCase(mockRepo, mockInventory, mockCategory, mockPricing)

		testItem := createTestItem(t)
		itemID := testItem.ID()

		req := &dto.UpdateInventoryRequest{
			Quantity: 50,
		}

		mockRepo.On("FindByID", mock.Anything, itemID).Return(testItem, nil)
		mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*item.Item")).Return(nil)

		result, err := useCase.UpdateInventory(context.Background(), itemID.String(), req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("item not found", func(t *testing.T) {
		mockRepo := &MockItemRepository{}
		mockInventory := &MockInventoryService{}
		mockCategory := &MockCategoryService{}
		mockPricing := &MockPricingService{}

		useCase := NewItemUseCase(mockRepo, mockInventory, mockCategory, mockPricing)

		itemID := item.NewItemID()
		req := &dto.UpdateInventoryRequest{
			Quantity: 50,
		}

		mockRepo.On("FindByID", mock.Anything, itemID).Return(nil, item.ItemNotFoundError(itemID))

		result, err := useCase.UpdateInventory(context.Background(), itemID.String(), req)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestItemUseCase_SearchItems(t *testing.T) {
	t.Run("search by query", func(t *testing.T) {
		mockRepo := &MockItemRepository{}
		mockInventory := &MockInventoryService{}
		mockCategory := &MockCategoryService{}
		mockPricing := &MockPricingService{}

		useCase := NewItemUseCase(mockRepo, mockInventory, mockCategory, mockPricing)

		testItem := createTestItem(t)
		req := &dto.SearchRequest{
			Query:    "test",
			Page:     1,
			PageSize: 10,
		}

		mockRepo.On("Search", mock.Anything, "test", 10, 0).Return([]*item.Item{testItem}, nil)

		result, err := useCase.SearchItems(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Items, 1)
		mockRepo.AssertExpectations(t)
	})

	t.Run("search by category", func(t *testing.T) {
		mockRepo := &MockItemRepository{}
		mockInventory := &MockInventoryService{}
		mockCategory := &MockCategoryService{}
		mockPricing := &MockPricingService{}

		useCase := NewItemUseCase(mockRepo, mockInventory, mockCategory, mockPricing)

		testItem := createTestItem(t)
		req := &dto.SearchRequest{
			Category: "electronics",
			Page:     1,
			PageSize: 10,
		}

		mockRepo.On("FindByCategory", mock.Anything, mock.AnythingOfType("item.Category"), 10, 0).Return([]*item.Item{testItem}, nil)

		result, err := useCase.SearchItems(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := &MockItemRepository{}
		mockInventory := &MockInventoryService{}
		mockCategory := &MockCategoryService{}
		mockPricing := &MockPricingService{}

		useCase := NewItemUseCase(mockRepo, mockInventory, mockCategory, mockPricing)

		req := &dto.SearchRequest{
			Query:    "test",
			Page:     1,
			PageSize: 10,
		}

		mockRepo.On("Search", mock.Anything, "test", 10, 0).Return(nil, assert.AnError)

		result, err := useCase.SearchItems(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestItemUseCase_DeleteItem(t *testing.T) {
	t.Run("successful deletion", func(t *testing.T) {
		mockRepo := &MockItemRepository{}
		mockInventory := &MockInventoryService{}
		mockCategory := &MockCategoryService{}
		mockPricing := &MockPricingService{}

		useCase := NewItemUseCase(mockRepo, mockInventory, mockCategory, mockPricing)

		itemID := item.NewItemID()
		mockRepo.On("Delete", mock.Anything, itemID).Return(nil)

		err := useCase.DeleteItem(context.Background(), itemID.String())

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("item not found", func(t *testing.T) {
		mockRepo := &MockItemRepository{}
		mockInventory := &MockInventoryService{}
		mockCategory := &MockCategoryService{}
		mockPricing := &MockPricingService{}

		useCase := NewItemUseCase(mockRepo, mockInventory, mockCategory, mockPricing)

		itemID := item.NewItemID()
		mockRepo.On("Delete", mock.Anything, itemID).Return(item.ItemNotFoundError(itemID))

		err := useCase.DeleteItem(context.Background(), itemID.String())

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

// Helper function to create a test item
func createTestItem(t *testing.T) *item.Item {
	t.Helper()

	sku, err := item.NewSKU("TEST-001")
	require.NoError(t, err)

	price, err := item.NewPrice(99.99, "USD")
	require.NoError(t, err)

	category, err := item.NewCategory("Electronics")
	require.NoError(t, err)

	testItem, err := item.NewItem(sku, "Test Item", "Test Description", price, category)
	require.NoError(t, err)

	return testItem
}
