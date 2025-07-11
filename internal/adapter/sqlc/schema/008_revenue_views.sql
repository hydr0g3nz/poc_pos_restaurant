-- internal/adapter/sqlc/schema/008_revenue_views.sql
-- +goose Up

-- Create view for daily revenue summary
CREATE VIEW daily_revenue_summary AS
SELECT 
    DATE(paid_at) as revenue_date,
    COUNT(*) as transaction_count,
    SUM(amount) as total_revenue,
    AVG(amount) as average_transaction,
    MIN(amount) as min_transaction,
    MAX(amount) as max_transaction
FROM payments
GROUP BY DATE(paid_at)
ORDER BY revenue_date DESC;

-- Create view for monthly revenue summary
CREATE VIEW monthly_revenue_summary AS
SELECT 
    DATE_TRUNC('month', paid_at) as revenue_month,
    COUNT(*) as transaction_count,
    SUM(amount) as total_revenue,
    AVG(amount) as average_transaction,
    COUNT(DISTINCT order_id) as unique_orders
FROM payments
GROUP BY DATE_TRUNC('month', paid_at)
ORDER BY revenue_month DESC;

-- Create view for payment method breakdown
CREATE VIEW payment_method_summary AS
SELECT 
    method,
    DATE(paid_at) as payment_date,
    COUNT(*) as transaction_count,
    SUM(amount) as total_amount,
    AVG(amount) as average_amount
FROM payments
GROUP BY method, DATE(paid_at)
ORDER BY payment_date DESC, total_amount DESC;

-- Create view for hourly revenue (useful for understanding peak hours)
CREATE VIEW hourly_revenue_summary AS
SELECT 
    DATE(paid_at) as revenue_date,
    EXTRACT(HOUR FROM paid_at) as hour,
    COUNT(*) as transaction_count,
    SUM(amount) as total_revenue,
    AVG(amount) as average_transaction
FROM payments
GROUP BY DATE(paid_at), EXTRACT(HOUR FROM paid_at)
ORDER BY revenue_date DESC, hour ASC;

-- Create view for top selling items with revenue
CREATE VIEW top_selling_items AS
SELECT 
    mi.id as item_id,
    mi.name as item_name,
    c.name as category_name,
    SUM(oi.quantity) as total_quantity_sold,
    SUM(oi.quantity * oi.unit_price) as total_revenue,
    AVG(oi.unit_price) as average_price,
    COUNT(DISTINCT oi.order_id) as unique_orders
FROM order_items oi
JOIN menu_items mi ON oi.item_id = mi.id
JOIN categories c ON mi.category_id = c.id
JOIN orders o ON oi.order_id = o.id
JOIN payments p ON o.id = p.order_id
GROUP BY mi.id, mi.name, c.name
ORDER BY total_revenue DESC;

-- Create indexes for better performance
CREATE INDEX idx_payments_paid_at_date ON payments(DATE(paid_at));
CREATE INDEX idx_payments_paid_at_month ON payments(DATE_TRUNC('month', paid_at));
CREATE INDEX idx_payments_method_date ON payments(method, DATE(paid_at));

-- +goose Down
DROP INDEX IF EXISTS idx_payments_method_date;
DROP INDEX IF EXISTS idx_payments_paid_at_month;
DROP INDEX IF EXISTS idx_payments_paid_at_date;

DROP VIEW IF EXISTS top_selling_items;
DROP VIEW IF EXISTS hourly_revenue_summary;
DROP VIEW IF EXISTS payment_method_summary;
DROP VIEW IF EXISTS monthly_revenue_summary;
DROP VIEW IF EXISTS daily_revenue_summary;