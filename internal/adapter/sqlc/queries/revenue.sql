-- internal/adapter/sqlc/queries/revenue.sql

-- name: GetDailyRevenue :one
SELECT 
    COALESCE(SUM(amount), 0)::decimal as total_revenue
FROM payments
WHERE paid_at >= $1 AND paid_at < $2;

-- name: GetMonthlyRevenue :one
SELECT 
    COALESCE(SUM(amount), 0)::decimal as total_revenue
FROM payments
WHERE paid_at >= $1 AND paid_at < $2;

-- name: GetDailyRevenueRange :many
SELECT 
    DATE(paid_at) as revenue_date,
    COALESCE(SUM(amount), 0)::decimal as total_revenue
FROM payments
WHERE paid_at >= $1 AND paid_at <= $2
GROUP BY DATE(paid_at)
ORDER BY revenue_date ASC;

-- name: GetMonthlyRevenueRange :many
SELECT 
    DATE_TRUNC('month', paid_at)::timestamp as revenue_month,
    COALESCE(SUM(amount), 0)::decimal as total_revenue
FROM payments
WHERE paid_at >= $1 AND paid_at <= $2
GROUP BY DATE_TRUNC('month', paid_at)
ORDER BY revenue_month ASC;

-- name: GetTotalRevenue :one
SELECT 
    COALESCE(SUM(amount), 0)::decimal as total_revenue
FROM payments
WHERE paid_at >= $1 AND paid_at <= $2;

-- name: GetRevenueByPaymentMethod :many
SELECT 
    method,
    COALESCE(SUM(amount), 0)::decimal as total_revenue,
    COUNT(*)::bigint as transaction_count
FROM payments
WHERE paid_at >= $1 AND paid_at <= $2
GROUP BY method
ORDER BY total_revenue DESC;

-- name: GetTopSellingItems :many
SELECT 
    mi.name as item_name,
    SUM(oi.quantity)::bigint as total_quantity,
    SUM(oi.quantity * oi.unit_price)::decimal as total_revenue
FROM order_items oi
JOIN menu_items mi ON oi.item_id = mi.id
JOIN orders o ON oi.order_id = o.id
JOIN payments p ON o.id = p.order_id
WHERE p.paid_at >= $1 AND p.paid_at <= $2
GROUP BY mi.id, mi.name
ORDER BY total_revenue DESC
LIMIT $3;

-- name: GetHourlyRevenue :many
SELECT 
    EXTRACT(HOUR FROM paid_at)::integer as hour,
    COALESCE(SUM(amount), 0)::decimal as total_revenue,
    COUNT(*)::bigint as transaction_count
FROM payments
WHERE DATE(paid_at) = DATE($1)
GROUP BY EXTRACT(HOUR FROM paid_at)
ORDER BY hour ASC;

-- name: GetRevenueStats :one
SELECT 
    COUNT(*)::bigint as total_transactions,
    COALESCE(SUM(amount), 0)::decimal as total_revenue,
    COALESCE(AVG(amount), 0)::decimal as average_transaction,
    COALESCE(MIN(amount), 0)::decimal as min_transaction,
    COALESCE(MAX(amount), 0)::decimal as max_transaction
FROM payments
WHERE paid_at >= $1 AND paid_at <= $2;