package persistence

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"item-pdp-service/internal/domain/item"
	"item-pdp-service/internal/infrastructure/database"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostgresItemRepository_Save(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresItemRepository(&database.DB{DB: db})
	ctx := context.Background()

	// Create test item
	testItem := createTestItem(t)

	t.Run("successful save", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO items").
			WithArgs(
				testItem.ID().String(),
				testItem.SKU().String(),
				testItem.Name(),
				testItem.Description(),
				sqlmock.AnyArg(), // price amount in cents
				testItem.Price().Currency(),
				testItem.Category().Name(),
				testItem.Category().Slug(),
				testItem.Inventory().Quantity(),
				sqlmock.AnyArg(), // images JSON
				sqlmock.AnyArg(), // attributes JSON
				testItem.Status().String(),
				sqlmock.AnyArg(), // created_at
				sqlmock.AnyArg(), // updated_at
			).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Save(ctx, testItem)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO items").
			WillReturnError(sql.ErrConnDone)

		err := repo.Save(ctx, testItem)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to save item")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPostgresItemRepository_FindByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresItemRepository(&database.DB{DB: db})
	ctx := context.Background()

	testItem := createTestItem(t)
	id := testItem.ID()

	t.Run("successful find", func(t *testing.T) {
		// Mock the row data
		images, _ := json.Marshal([]map[string]interface{}{})
		attributes, _ := json.Marshal(map[string]string{})

		rows := sqlmock.NewRows([]string{
			"id", "sku", "name", "description", "price_amount", "price_currency",
			"category_name", "category_slug", "inventory_quantity", "images",
			"attributes", "status", "created_at", "updated_at",
		}).AddRow(
			testItem.ID().String(),
			testItem.SKU().String(),
			testItem.Name(),
			testItem.Description(),
			int64(testItem.Price().Amount()*100),
			testItem.Price().Currency(),
			testItem.Category().Name(),
			testItem.Category().Slug(),
			testItem.Inventory().Quantity(),
			images,
			attributes,
			testItem.Status().String(),
			testItem.CreatedAt(),
			testItem.UpdatedAt(),
		)

		mock.ExpectQuery("SELECT (.+) FROM items WHERE id = \\$1").
			WithArgs(id.String()).
			WillReturnRows(rows)

		result, err := repo.FindByID(ctx, id)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, testItem.SKU().String(), result.SKU().String())
		assert.Equal(t, testItem.Name(), result.Name())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("item not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM items WHERE id = \\$1").
			WithArgs(id.String()).
			WillReturnError(sql.ErrNoRows)

		result, err := repo.FindByID(ctx, id)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "not found")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM items WHERE id = \\$1").
			WithArgs(id.String()).
			WillReturnError(sql.ErrConnDone)

		result, err := repo.FindByID(ctx, id)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to find item by ID")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPostgresItemRepository_FindBySKU(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresItemRepository(&database.DB{DB: db})
	ctx := context.Background()

	testItem := createTestItem(t)
	sku := testItem.SKU()

	t.Run("successful find", func(t *testing.T) {
		images, _ := json.Marshal([]map[string]interface{}{})
		attributes, _ := json.Marshal(map[string]string{})

		rows := sqlmock.NewRows([]string{
			"id", "sku", "name", "description", "price_amount", "price_currency",
			"category_name", "category_slug", "inventory_quantity", "images",
			"attributes", "status", "created_at", "updated_at",
		}).AddRow(
			testItem.ID().String(),
			testItem.SKU().String(),
			testItem.Name(),
			testItem.Description(),
			int64(testItem.Price().Amount()*100),
			testItem.Price().Currency(),
			testItem.Category().Name(),
			testItem.Category().Slug(),
			testItem.Inventory().Quantity(),
			images,
			attributes,
			testItem.Status().String(),
			testItem.CreatedAt(),
			testItem.UpdatedAt(),
		)

		mock.ExpectQuery("SELECT (.+) FROM items WHERE sku = \\$1").
			WithArgs(sku.String()).
			WillReturnRows(rows)

		result, err := repo.FindBySKU(ctx, sku)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, testItem.SKU().String(), result.SKU().String())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("item not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM items WHERE sku = \\$1").
			WithArgs(sku.String()).
			WillReturnError(sql.ErrNoRows)

		result, err := repo.FindBySKU(ctx, sku)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "not found")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPostgresItemRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresItemRepository(&database.DB{DB: db})
	ctx := context.Background()

	testItem := createTestItem(t)

	t.Run("successful update", func(t *testing.T) {
		mock.ExpectExec("UPDATE items SET").
			WithArgs(
				testItem.ID().String(),
				testItem.Name(),
				testItem.Description(),
				sqlmock.AnyArg(), // price amount
				testItem.Price().Currency(),
				testItem.Category().Name(),
				testItem.Category().Slug(),
				testItem.Inventory().Quantity(),
				sqlmock.AnyArg(), // images JSON
				sqlmock.AnyArg(), // attributes JSON
				testItem.Status().String(),
				sqlmock.AnyArg(), // updated_at
			).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Update(ctx, testItem)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("item not found", func(t *testing.T) {
		mock.ExpectExec("UPDATE items SET").
			WillReturnResult(sqlmock.NewResult(1, 0)) // 0 rows affected

		err := repo.Update(ctx, testItem)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPostgresItemRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresItemRepository(&database.DB{DB: db})
	ctx := context.Background()

	id := item.NewItemID()

	t.Run("successful deletion", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM items WHERE id = \\$1").
			WithArgs(id.String()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Delete(ctx, id)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("item not found", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM items WHERE id = \\$1").
			WithArgs(id.String()).
			WillReturnResult(sqlmock.NewResult(1, 0))

		err := repo.Delete(ctx, id)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPostgresItemRepository_ExistsBySKU(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresItemRepository(&database.DB{DB: db})
	ctx := context.Background()

	sku, _ := item.NewSKU("TEST-001")

	t.Run("exists", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)
		mock.ExpectQuery("SELECT EXISTS").
			WithArgs(sku.String()).
			WillReturnRows(rows)

		exists, err := repo.ExistsBySKU(ctx, sku)

		assert.NoError(t, err)
		assert.True(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("does not exist", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"exists"}).AddRow(false)
		mock.ExpectQuery("SELECT EXISTS").
			WithArgs(sku.String()).
			WillReturnRows(rows)

		exists, err := repo.ExistsBySKU(ctx, sku)

		assert.NoError(t, err)
		assert.False(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPostgresItemRepository_Search(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresItemRepository(&database.DB{DB: db})
	ctx := context.Background()

	query := "test"
	limit := 10
	offset := 0

	t.Run("successful search", func(t *testing.T) {
		images, _ := json.Marshal([]map[string]interface{}{})
		attributes, _ := json.Marshal(map[string]string{})

		rows := sqlmock.NewRows([]string{
			"id", "sku", "name", "description", "price_amount", "price_currency",
			"category_name", "category_slug", "inventory_quantity", "images",
			"attributes", "status", "created_at", "updated_at",
		}).AddRow(
			"550e8400-e29b-41d4-a716-446655440000",
			"TEST-001",
			"Test Item",
			"Test Description",
			9999,
			"USD",
			"Electronics",
			"electronics",
			10,
			images,
			attributes,
			"active",
			time.Now(),
			time.Now(),
		)

		mock.ExpectQuery("SELECT (.+) FROM items WHERE \\(name ILIKE '%test%' OR description ILIKE '%test%' OR sku ILIKE '%test%'\\) ORDER BY created_at DESC LIMIT 10 OFFSET 0").
			WillReturnRows(rows)

		results, err := repo.Search(ctx, query, limit, offset)

		assert.NoError(t, err)
		assert.Len(t, results, 1)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("no results", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"id", "sku", "name", "description", "price_amount", "price_currency",
			"category_name", "category_slug", "inventory_quantity", "images",
			"attributes", "status", "created_at", "updated_at",
		})

		mock.ExpectQuery("SELECT (.+) FROM items WHERE \\(name ILIKE '%test%' OR description ILIKE '%test%' OR sku ILIKE '%test%'\\) ORDER BY created_at DESC LIMIT 10 OFFSET 0").
			WillReturnRows(rows)

		results, err := repo.Search(ctx, query, limit, offset)

		assert.NoError(t, err)
		assert.Len(t, results, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
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
