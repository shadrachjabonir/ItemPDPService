-- Create extension for UUID generation if not exists
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create items table
CREATE TABLE items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sku VARCHAR(20) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price_amount BIGINT NOT NULL, -- stored in cents
    price_currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    category_name VARCHAR(100) NOT NULL,
    category_slug VARCHAR(100) NOT NULL,
    inventory_quantity INTEGER NOT NULL DEFAULT 0,
    images JSONB DEFAULT '[]'::jsonb,
    attributes JSONB DEFAULT '{}'::jsonb,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes for performance
CREATE INDEX idx_items_sku ON items(sku);
CREATE INDEX idx_items_category_slug ON items(category_slug);
CREATE INDEX idx_items_status ON items(status);
CREATE INDEX idx_items_inventory_quantity ON items(inventory_quantity);
CREATE INDEX idx_items_created_at ON items(created_at);
CREATE INDEX idx_items_updated_at ON items(updated_at);

-- Create partial index for available items
CREATE INDEX idx_items_available ON items(status, inventory_quantity) 
WHERE status = 'active' AND inventory_quantity > 0;

-- Create GIN index for searching in name and description
CREATE INDEX idx_items_search ON items USING gin(to_tsvector('english', name || ' ' || description));

-- Create trigger to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_items_updated_at 
    BEFORE UPDATE ON items 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Add check constraints
ALTER TABLE items ADD CONSTRAINT chk_price_amount_positive CHECK (price_amount >= 0);
ALTER TABLE items ADD CONSTRAINT chk_inventory_quantity_non_negative CHECK (inventory_quantity >= 0);
ALTER TABLE items ADD CONSTRAINT chk_status_valid CHECK (status IN ('active', 'inactive', 'draft', 'archived'));
ALTER TABLE items ADD CONSTRAINT chk_currency_length CHECK (length(price_currency) = 3); 