-- +goose Up
CREATE TABLE tables (
    id SERIAL PRIMARY KEY,
    table_number INTEGER UNIQUE NOT NULL CHECK (table_number > 0),
    qr_code VARCHAR(255) UNIQUE NOT NULL,
    seating INTEGER NOT NULL DEFAULT 4 CHECK (seating >= 0),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_tables_table_number ON tables(table_number);
CREATE INDEX idx_tables_qr_code ON tables(qr_code);
CREATE INDEX idx_tables_is_active ON tables(is_active);

-- Insert sample tables
INSERT INTO tables (table_number, qr_code, seating) VALUES
(1, '/order?table=1', 4),
(2, '/order?table=2', 2),
(3, '/order?table=3', 6),
(4, '/order?table=4', 4),
(5, '/order?table=5', 8),
(6, '/order?table=6', 4),
(7, '/order?table=7', 2),
(8, '/order?table=8', 4),
(9, '/order?table=9', 6),
(10, '/order?table=10', 4);

-- +goose Down
DROP INDEX IF EXISTS idx_tables_is_active;
DROP INDEX IF EXISTS idx_tables_qr_code;
DROP INDEX IF EXISTS idx_tables_table_number;
DROP TABLE IF EXISTS tables;