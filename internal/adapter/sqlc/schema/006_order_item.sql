-- +goose Up
CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    item_id INTEGER NOT NULL REFERENCES menu_items(id) ON DELETE RESTRICT,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    unit_price DECIMAL(10,2) NOT NULL CHECK (unit_price >= 0),
    notes TEXT DEFAULT '',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_order_items_item_id ON order_items(item_id);
CREATE INDEX idx_order_items_created_at ON order_items(created_at);
CREATE INDEX idx_order_items_order_item ON order_items(order_id, item_id);

-- Unique constraint to prevent duplicate items in same order
CREATE UNIQUE INDEX idx_order_items_unique_order_item ON order_items(order_id, item_id);

-- +goose Down
DROP INDEX IF EXISTS idx_order_items_unique_order_item;
DROP INDEX IF EXISTS idx_order_items_order_item;
DROP INDEX IF EXISTS idx_order_items_created_at;
DROP INDEX IF EXISTS idx_order_items_item_id;
DROP INDEX IF EXISTS idx_order_items_order_id;
DROP TABLE IF EXISTS order_items;