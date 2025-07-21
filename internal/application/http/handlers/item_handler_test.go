package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"item-pdp-service/internal/application/dto"
	"item-pdp-service/internal/domain/item"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockItemUseCase is a mock implementation of usecase.ItemUseCase
type MockItemUseCase struct {
	mock.Mock
}

func (m *MockItemUseCase) CreateItem(ctx context.Context, req *dto.CreateItemRequest) (*dto.ItemResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ItemResponse), args.Error(1)
}

func (m *MockItemUseCase) GetItemByID(ctx context.Context, id string) (*dto.ItemResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ItemResponse), args.Error(1)
}

func (m *MockItemUseCase) GetItemBySKU(ctx context.Context, sku string) (*dto.ItemResponse, error) {
	args := m.Called(ctx, sku)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ItemResponse), args.Error(1)
}

func (m *MockItemUseCase) UpdateItem(ctx context.Context, id string, req *dto.UpdateItemRequest) (*dto.ItemResponse, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ItemResponse), args.Error(1)
}

func (m *MockItemUseCase) UpdateInventory(ctx context.Context, id string, req *dto.UpdateInventoryRequest) (*dto.ItemResponse, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ItemResponse), args.Error(1)
}

func (m *MockItemUseCase) AddImage(ctx context.Context, id string, req *dto.AddImageRequest) (*dto.ItemResponse, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ItemResponse), args.Error(1)
}

func (m *MockItemUseCase) DeactivateItem(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockItemUseCase) ActivateItem(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockItemUseCase) DeleteItem(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockItemUseCase) SearchItems(ctx context.Context, req *dto.SearchRequest) (*dto.ItemListResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ItemListResponse), args.Error(1)
}

func (m *MockItemUseCase) GetItemsByCategory(ctx context.Context, category string, page, pageSize int) (*dto.ItemListResponse, error) {
	args := m.Called(ctx, category, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ItemListResponse), args.Error(1)
}

func (m *MockItemUseCase) GetAvailableItems(ctx context.Context, page, pageSize int) (*dto.ItemListResponse, error) {
	args := m.Called(ctx, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ItemListResponse), args.Error(1)
}

func TestItemHandler_CreateItem(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := new(MockItemUseCase)
	handler := NewItemHandler(mockUseCase)

	t.Run("successful creation", func(t *testing.T) {
		req := &dto.CreateItemRequest{
			SKU:         "TEST-001",
			Name:        "Test Item",
			Description: "Test Description",
			Price:       99.99,
			Currency:    "USD",
			Category:    "Electronics",
			Inventory:   10,
		}

		expectedResponse := &dto.ItemResponse{
			ID:          "550e8400-e29b-41d4-a716-446655440000",
			SKU:         req.SKU,
			Name:        req.Name,
			Description: req.Description,
			Price:       req.Price,
			Currency:    req.Currency,
		}

		mockUseCase.On("CreateItem", mock.Anything, mock.AnythingOfType("*dto.CreateItemRequest")).Return(expectedResponse, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body, _ := json.Marshal(req)
		c.Request = httptest.NewRequest("POST", "/items", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateItem(c)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response dto.ItemResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse.SKU, response.SKU)
		assert.Equal(t, expectedResponse.Name, response.Name)

		mockUseCase.AssertExpectations(t)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest("POST", "/items", bytes.NewBuffer([]byte("invalid json")))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateItem(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("use case error", func(t *testing.T) {
		req := &dto.CreateItemRequest{
			SKU:      "TEST-001",
			Name:     "Test Item",
			Price:    99.99,
			Currency: "USD",
			Category: "Electronics",
		}

		mockUseCase.On("CreateItem", mock.Anything, mock.AnythingOfType("*dto.CreateItemRequest")).Return(nil, errors.New("use case error")).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body, _ := json.Marshal(req)
		c.Request = httptest.NewRequest("POST", "/items", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateItem(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockUseCase.AssertExpectations(t)
	})
}

func TestItemHandler_GetItem(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := new(MockItemUseCase)
	handler := NewItemHandler(mockUseCase)

	t.Run("successful retrieval", func(t *testing.T) {
		itemID := "550e8400-e29b-41d4-a716-446655440000"
		expectedResponse := &dto.ItemResponse{
			ID:   itemID,
			SKU:  "TEST-001",
			Name: "Test Item",
		}

		mockUseCase.On("GetItemByID", mock.Anything, itemID).Return(expectedResponse, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: itemID}}
		c.Request = httptest.NewRequest("GET", "/items/"+itemID, nil)

		handler.GetItem(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dto.ItemResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse.ID, response.ID)

		mockUseCase.AssertExpectations(t)
	})

	t.Run("empty ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{}
		c.Request = httptest.NewRequest("GET", "/items/", nil)

		handler.GetItem(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("item not found", func(t *testing.T) {
		itemID := "550e8400-e29b-41d4-a716-446655440000"

		mockUseCase.On("GetItemByID", mock.Anything, itemID).Return(nil, item.ErrItemNotFound).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: itemID}}
		c.Request = httptest.NewRequest("GET", "/items/"+itemID, nil)

		handler.GetItem(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockUseCase.AssertExpectations(t)
	})
}

func TestItemHandler_UpdateInventory(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := new(MockItemUseCase)
	handler := NewItemHandler(mockUseCase)

	t.Run("successful update", func(t *testing.T) {
		itemID := "550e8400-e29b-41d4-a716-446655440000"
		req := &dto.UpdateInventoryRequest{Quantity: 50}

		expectedResponse := &dto.ItemResponse{
			ID: itemID,
			Inventory: dto.InventoryResponse{
				Quantity:    50,
				IsAvailable: true,
			},
		}

		mockUseCase.On("UpdateInventory", mock.Anything, itemID, mock.AnythingOfType("*dto.UpdateInventoryRequest")).Return(expectedResponse, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: itemID}}

		body, _ := json.Marshal(req)
		c.Request = httptest.NewRequest("PATCH", "/items/"+itemID+"/inventory", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateInventory(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dto.ItemResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 50, response.Inventory.Quantity)

		mockUseCase.AssertExpectations(t)
	})
}

func TestItemHandler_DeleteItem(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := new(MockItemUseCase)
	handler := NewItemHandler(mockUseCase)

	t.Run("successful deletion", func(t *testing.T) {
		itemID := "550e8400-e29b-41d4-a716-446655440000"

		mockUseCase.On("DeleteItem", mock.Anything, itemID).Return(nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: itemID}}
		c.Request = httptest.NewRequest("DELETE", "/items/"+itemID, nil)

		handler.DeleteItem(c)

		assert.Equal(t, http.StatusNoContent, w.Code)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("item not found", func(t *testing.T) {
		itemID := "550e8400-e29b-41d4-a716-446655440000"

		mockUseCase.On("DeleteItem", mock.Anything, itemID).Return(item.ErrItemNotFound).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: itemID}}
		c.Request = httptest.NewRequest("DELETE", "/items/"+itemID, nil)

		handler.DeleteItem(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockUseCase.AssertExpectations(t)
	})
}

func TestItemHandler_SearchItems(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := new(MockItemUseCase)
	handler := NewItemHandler(mockUseCase)

	t.Run("successful search", func(t *testing.T) {
		expectedResponse := &dto.ItemListResponse{
			Items: []dto.ItemResponse{
				{ID: "1", Name: "Item 1"},
				{ID: "2", Name: "Item 2"},
			},
			Total:      2,
			Page:       1,
			PageSize:   10,
			TotalPages: 1,
		}

		mockUseCase.On("SearchItems", mock.Anything, mock.AnythingOfType("*dto.SearchRequest")).Return(expectedResponse, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/items/search?query=test&page=1&page_size=10", nil)

		handler.SearchItems(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dto.ItemListResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Items, 2)

		mockUseCase.AssertExpectations(t)
	})

	t.Run("default pagination", func(t *testing.T) {
		expectedResponse := &dto.ItemListResponse{Items: []dto.ItemResponse{}}

		mockUseCase.On("SearchItems", mock.Anything, mock.MatchedBy(func(req *dto.SearchRequest) bool {
			return req.Page == 1 && req.PageSize == 10
		})).Return(expectedResponse, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/items/search", nil)

		handler.SearchItems(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockUseCase.AssertExpectations(t)
	})
}
