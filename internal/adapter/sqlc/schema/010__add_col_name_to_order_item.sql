-- +goose Up
ALTER TABLE order_items ADD COLUMN name VARCHAR(255) NULL;
-- +goose Down
ALTER TABLE order_items DROP COLUMN name;