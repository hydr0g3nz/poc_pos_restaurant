-- +goose Up
CREATE TYPE payment_method AS ENUM ('cash', 'credit_card', 'wallet');

CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES orders(id) ON DELETE RESTRICT,
    amount DECIMAL(10,2) NOT NULL CHECK (amount > 0),
    method payment_method NOT NULL,
    reference TEXT DEFAULT '',
    paid_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_payments_order_id ON payments(order_id);
CREATE INDEX idx_payments_method ON payments(method);
CREATE INDEX idx_payments_paid_at ON payments(paid_at);
CREATE INDEX idx_payments_method_paid_at ON payments(method, paid_at);

-- Unique constraint to prevent duplicate payments for same order
CREATE UNIQUE INDEX idx_payments_unique_order ON payments(order_id);

-- +goose Down
DROP INDEX IF EXISTS idx_payments_unique_order;
DROP INDEX IF EXISTS idx_payments_method_paid_at;
DROP INDEX IF EXISTS idx_payments_paid_at;
DROP INDEX IF EXISTS idx_payments_method;
DROP INDEX IF EXISTS idx_payments_order_id;
DROP TABLE IF EXISTS payments;
DROP TYPE IF EXISTS payment_method;