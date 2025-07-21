-- Drop trigger
DROP TRIGGER IF EXISTS update_items_updated_at ON items;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_items_search;
DROP INDEX IF EXISTS idx_items_available;
DROP INDEX IF EXISTS idx_items_updated_at;
DROP INDEX IF EXISTS idx_items_created_at;
DROP INDEX IF EXISTS idx_items_inventory_quantity;
DROP INDEX IF EXISTS idx_items_status;
DROP INDEX IF EXISTS idx_items_category_slug;
DROP INDEX IF EXISTS idx_items_sku;

-- Drop table
DROP TABLE IF EXISTS items; 