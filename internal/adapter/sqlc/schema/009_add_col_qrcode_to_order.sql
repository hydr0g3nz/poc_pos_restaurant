-- +goose Up
ALTER TABLE orders ADD COLUMN qrcode VARCHAR(255) NULL;
-- +goose Down
ALTER TABLE orders DROP COLUMN qrcode;