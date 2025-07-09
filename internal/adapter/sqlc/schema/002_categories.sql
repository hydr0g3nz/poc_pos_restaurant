-- +goose Up
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_categories_name ON categories(name);

-- Insert default categories
INSERT INTO categories (name) VALUES
('ของคาว'),
('ของหวาน'),
('ของทานเล่น'),
('โรตี');

-- +goose Down
DROP INDEX IF EXISTS idx_categories_name;
DROP TABLE IF EXISTS categories;