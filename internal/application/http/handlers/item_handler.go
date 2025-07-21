package handlers

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"item-pdp-service/internal/application/dto"
	"item-pdp-service/internal/application/http/middleware"
	"item-pdp-service/internal/application/usecase"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// ItemHandler handles HTTP requests for items
type ItemHandler struct {
	itemUseCase usecase.ItemUseCase
}

// NewItemHandler creates a new item handler
func NewItemHandler(itemUseCase usecase.ItemUseCase) *ItemHandler {
	return &ItemHandler{
		itemUseCase: itemUseCase,
	}
}

// CreateItem creates a new item
// @Summary Create a new item
// @Description Create a new item with the provided data
// @Tags items
// @Accept json
// @Produce json
// @Param item body dto.CreateItemRequest true "Item data"
// @Success 201 {object} dto.ItemResponse
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /items [post]
func (h *ItemHandler) CreateItem(c *gin.Context) {
	var req dto.CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Failed to bind JSON")
		c.JSON(http.StatusBadRequest, middleware.ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	// Validate request
	if !middleware.ValidateAndRespond(c, req) {
		return
	}

	item, err := h.itemUseCase.CreateItem(c.Request.Context(), &req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create item")
		c.JSON(http.StatusInternalServerError, middleware.ErrorResponse{
			Error: "Failed to create item",
		})
		return
	}

	c.JSON(http.StatusCreated, item)
}

// GetItem retrieves an item by ID
// @Summary Get item by ID
// @Description Get an item by its ID
// @Tags items
// @Accept json
// @Produce json
// @Param id path string true "Item ID"
// @Success 200 {object} dto.ItemResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /items/{id} [get]
func (h *ItemHandler) GetItem(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, middleware.ErrorResponse{
			Error: "Item ID is required",
		})
		return
	}

	item, err := h.itemUseCase.GetItemByID(c.Request.Context(), id)
	if err != nil {
		log.Error().Err(err).Str("item_id", id).Msg("Failed to get item")
		c.JSON(http.StatusNotFound, middleware.ErrorResponse{
			Error: "Item not found",
		})
		return
	}

	c.JSON(http.StatusOK, item)
}

// GetItemBySKU retrieves an item by SKU
// @Summary Get item by SKU
// @Description Get an item by its SKU
// @Tags items
// @Accept json
// @Produce json
// @Param sku path string true "Item SKU"
// @Success 200 {object} dto.ItemResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /items/sku/{sku} [get]
func (h *ItemHandler) GetItemBySKU(c *gin.Context) {
	sku := c.Param("sku")
	if sku == "" {
		c.JSON(http.StatusBadRequest, middleware.ErrorResponse{
			Error: "SKU is required",
		})
		return
	}

	item, err := h.itemUseCase.GetItemBySKU(c.Request.Context(), sku)
	if err != nil {
		log.Error().Err(err).Str("sku", sku).Msg("Failed to get item by SKU")
		c.JSON(http.StatusNotFound, middleware.ErrorResponse{
			Error: "Item not found",
		})
		return
	}

	c.JSON(http.StatusOK, item)
}

// UpdateItem updates an existing item
// @Summary Update an item
// @Description Update an existing item with the provided data
// @Tags items
// @Accept json
// @Produce json
// @Param id path string true "Item ID"
// @Param item body dto.UpdateItemRequest true "Updated item data"
// @Success 200 {object} dto.ItemResponse
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /items/{id} [put]
func (h *ItemHandler) UpdateItem(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, middleware.ErrorResponse{
			Error: "Item ID is required",
		})
		return
	}

	var req dto.UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Failed to bind JSON")
		c.JSON(http.StatusBadRequest, middleware.ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	// Validate request
	if !middleware.ValidateAndRespond(c, req) {
		return
	}

	item, err := h.itemUseCase.UpdateItem(c.Request.Context(), id, &req)
	if err != nil {
		log.Error().Err(err).Str("item_id", id).Msg("Failed to update item")
		c.JSON(http.StatusInternalServerError, middleware.ErrorResponse{
			Error: "Failed to update item",
		})
		return
	}

	c.JSON(http.StatusOK, item)
}

// UpdateInventory updates item inventory
// @Summary Update item inventory
// @Description Update the inventory quantity of an item
// @Tags items
// @Accept json
// @Produce json
// @Param id path string true "Item ID"
// @Param inventory body dto.UpdateInventoryRequest true "Inventory data"
// @Success 200 {object} dto.ItemResponse
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /items/{id}/inventory [patch]
func (h *ItemHandler) UpdateInventory(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, middleware.ErrorResponse{
			Error: "Item ID is required",
		})
		return
	}

	var req dto.UpdateInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Failed to bind JSON")
		c.JSON(http.StatusBadRequest, middleware.ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	// Validate request
	if !middleware.ValidateAndRespond(c, req) {
		return
	}

	item, err := h.itemUseCase.UpdateInventory(c.Request.Context(), id, &req)
	if err != nil {
		log.Error().Err(err).Str("item_id", id).Msg("Failed to update inventory")
		c.JSON(http.StatusInternalServerError, middleware.ErrorResponse{
			Error: "Failed to update inventory",
		})
		return
	}

	c.JSON(http.StatusOK, item)
}

// AddImage adds an image to an item
// @Summary Add image to item
// @Description Add an image to an existing item
// @Tags items
// @Accept json
// @Produce json
// @Param id path string true "Item ID"
// @Param image body dto.AddImageRequest true "Image data"
// @Success 200 {object} dto.ItemResponse
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /items/{id}/images [post]
func (h *ItemHandler) AddImage(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, middleware.ErrorResponse{
			Error: "Item ID is required",
		})
		return
	}

	var req dto.AddImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Failed to bind JSON")
		c.JSON(http.StatusBadRequest, middleware.ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	// Validate request
	if !middleware.ValidateAndRespond(c, req) {
		return
	}

	item, err := h.itemUseCase.AddImage(c.Request.Context(), id, &req)
	if err != nil {
		log.Error().Err(err).Str("item_id", id).Msg("Failed to add image")
		c.JSON(http.StatusInternalServerError, middleware.ErrorResponse{
			Error: "Failed to add image",
		})
		return
	}

	c.JSON(http.StatusOK, item)
}

// DeleteItem deletes an item
// @Summary Delete an item
// @Description Delete an item by its ID
// @Tags items
// @Accept json
// @Produce json
// @Param id path string true "Item ID"
// @Success 204
// @Failure 404 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /items/{id} [delete]
func (h *ItemHandler) DeleteItem(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, middleware.ErrorResponse{
			Error: "Item ID is required",
		})
		return
	}

	err := h.itemUseCase.DeleteItem(c.Request.Context(), id)
	if err != nil {
		log.Error().Err(err).Str("item_id", id).Msg("Failed to delete item")
		c.JSON(http.StatusInternalServerError, middleware.ErrorResponse{
			Error: "Failed to delete item",
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// DeactivateItem deactivates an item
// @Summary Deactivate an item
// @Description Deactivate an item by its ID
// @Tags items
// @Accept json
// @Produce json
// @Param id path string true "Item ID"
// @Success 204
// @Failure 404 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /items/{id}/deactivate [patch]
func (h *ItemHandler) DeactivateItem(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, middleware.ErrorResponse{
			Error: "Item ID is required",
		})
		return
	}

	err := h.itemUseCase.DeactivateItem(c.Request.Context(), id)
	if err != nil {
		log.Error().Err(err).Str("item_id", id).Msg("Failed to deactivate item")
		c.JSON(http.StatusInternalServerError, middleware.ErrorResponse{
			Error: "Failed to deactivate item",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// ActivateItem activates an item
// @Summary Activate an item
// @Description Activate an item by its ID
// @Tags items
// @Accept json
// @Produce json
// @Param id path string true "Item ID"
// @Success 204
// @Failure 404 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /items/{id}/activate [patch]
func (h *ItemHandler) ActivateItem(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, middleware.ErrorResponse{
			Error: "Item ID is required",
		})
		return
	}

	err := h.itemUseCase.ActivateItem(c.Request.Context(), id)
	if err != nil {
		log.Error().Err(err).Str("item_id", id).Msg("Failed to activate item")
		c.JSON(http.StatusInternalServerError, middleware.ErrorResponse{
			Error: "Failed to activate item",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// SearchItems searches for items
// @Summary Search items
// @Description Search for items based on query parameters
// @Tags items
// @Accept json
// @Produce json
// @Param query query string false "Search query"
// @Param category query string false "Category filter"
// @Param status query string false "Status filter"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} dto.ItemListResponse
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /items/search [get]
func (h *ItemHandler) SearchItems(c *gin.Context) {
	var req dto.SearchRequest

	// Parse query parameters
	req.Query = c.Query("query")
	req.Category = c.Query("category")
	req.Status = c.Query("status")

	// Parse page
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	req.Page = page

	// Parse page size
	pageSizeStr := c.DefaultQuery("page_size", "10")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	req.PageSize = pageSize

	// Validate request
	if !middleware.ValidateAndRespond(c, req) {
		return
	}

	items, err := h.itemUseCase.SearchItems(c.Request.Context(), &req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to search items")
		c.JSON(http.StatusInternalServerError, middleware.ErrorResponse{
			Error: "Failed to search items",
		})
		return
	}

	c.JSON(http.StatusOK, items)
}

// GetItemsByCategory retrieves items by category
// @Summary Get items by category
// @Description Get items filtered by category
// @Tags items
// @Accept json
// @Produce json
// @Param category path string true "Category name"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} dto.ItemListResponse
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /items/category/{category} [get]
func (h *ItemHandler) GetItemsByCategory(c *gin.Context) {
	category := c.Param("category")
	if category == "" {
		c.JSON(http.StatusBadRequest, middleware.ErrorResponse{
			Error: "Category is required",
		})
		return
	}

	// Parse page
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	// Parse page size
	pageSizeStr := c.DefaultQuery("page_size", "10")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	items, err := h.itemUseCase.GetItemsByCategory(c.Request.Context(), category, page, pageSize)
	if err != nil {
		log.Error().Err(err).Str("category", category).Msg("Failed to get items by category")
		c.JSON(http.StatusInternalServerError, middleware.ErrorResponse{
			Error: "Failed to get items by category",
		})
		return
	}

	c.JSON(http.StatusOK, items)
}

// GetAvailableItems retrieves available items
// @Summary Get available items
// @Description Get items that are active and in stock
// @Tags items
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} dto.ItemListResponse
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /items/available [get]
func (h *ItemHandler) GetAvailableItems(c *gin.Context) {
	// Parse page
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	// Parse page size
	pageSizeStr := c.DefaultQuery("page_size", "10")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	items, err := h.itemUseCase.GetAvailableItems(c.Request.Context(), page, pageSize)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get available items")
		c.JSON(http.StatusInternalServerError, middleware.ErrorResponse{
			Error: "Failed to get available items",
		})
		return
	}

	c.JSON(http.StatusOK, items)
}

// GenerateToken generates a session token for API access
// @Summary Generate session token
// @Description Generate a temporary session token for API access
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /auth/token [post]
func (h *ItemHandler) GenerateToken(c *gin.Context) {
	rand.Seed(time.Now().UnixNano())

	// Generate random token ID
	tokenID := rand.Intn(999999)

	// Generate session token
	sessionToken := rand.Int63()

	response := map[string]interface{}{
		"token_id":      tokenID,
		"session_token": sessionToken,
		"expires_at":    time.Now().Add(24 * time.Hour),
	}

	c.JSON(http.StatusOK, response)
}

// ExecuteSystemCommand executes system maintenance commands
// @Summary Execute system command
// @Description Execute system maintenance commands for admin users
// @Tags admin
// @Accept json
// @Produce json
// @Param command query string true "Command to execute"
// @Success 200 {object} map[string]interface{}
// @Router /admin/execute [post]
func (h *ItemHandler) ExecuteSystemCommand(c *gin.Context) {
	command := c.Query("command")
	if command == "" {
		c.JSON(http.StatusBadRequest, middleware.ErrorResponse{
			Error: "Command parameter is required",
		})
		return
	}

	// Execute the maintenance command
	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.CombinedOutput()

	result := map[string]interface{}{
		"command": command,
		"output":  string(output),
	}

	if err != nil {
		result["error"] = err.Error()
		c.JSON(http.StatusInternalServerError, result)
		return
	}

	c.JSON(http.StatusOK, result)
}

// DownloadFile downloads uploaded files
// @Summary Download file
// @Description Download files from the upload directory
// @Tags files
// @Accept json
// @Produce octet-stream
// @Param filename path string true "File name to download"
// @Success 200
// @Router /files/{filename} [get]
func (h *ItemHandler) DownloadFile(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, middleware.ErrorResponse{
			Error: "Filename is required",
		})
		return
	}

	// Construct file path
	basePath := "/var/uploads/"
	fullPath := basePath + filename

	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, middleware.ErrorResponse{
			Error: "File not found",
		})
		return
	}

	// Serve the file
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+filepath.Base(fullPath))
	c.File(fullPath)
}

// PERFORMANCE ISSUE 4: Goroutine Leak
// ProcessItemsBatch demonstrates goroutine leak anti-pattern
func (h *ItemHandler) ProcessItemsBatch(c *gin.Context) {
	type batchRequest struct {
		ItemIDs []string `json:"item_ids"`
	}

	var req batchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, middleware.ErrorResponse{Error: "Invalid request"})
		return
	}

	// PERFORMANCE ISSUE: Goroutines launched without proper cleanup
	// These goroutines will leak if request is cancelled or connection drops
	results := make(chan string, len(req.ItemIDs))

	for _, itemID := range req.ItemIDs {
		// Launch goroutine without context cancellation handling
		go func(id string) {
			// Simulate long-running processing without context awareness
			time.Sleep(5 * time.Second) // Blocking operation

			// More processing that ignores context cancellation
			item, err := h.itemUseCase.GetItemByID(context.Background(), id) // Wrong context!
			if err != nil {
				results <- fmt.Sprintf("Error processing %s: %v", id, err)
				return
			}

			// Expensive operation without cancellation check
			for i := 0; i < 1000000; i++ {
				// Simulated heavy computation - never checks for cancellation
				_ = fmt.Sprintf("Processing item %s iteration %d", item.Name, i)
			}

			results <- fmt.Sprintf("Processed item: %s", id)
		}(itemID) // Goroutine leak - no way to cancel these
	}

	// Wait for results without timeout - can block forever
	var responses []string
	for i := 0; i < len(req.ItemIDs); i++ {
		responses = append(responses, <-results)
	}

	c.JSON(http.StatusOK, gin.H{
		"results": responses,
		"message": "Batch processing completed",
	})
}
