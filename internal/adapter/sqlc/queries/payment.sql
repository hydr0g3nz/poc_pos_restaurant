-- name: CreatePayment :one
INSERT INTO payments (order_id, amount, method, reference)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetPaymentByID :one
SELECT * FROM payments
WHERE id = $1;

-- name: GetPaymentByOrderID :one
SELECT * FROM payments
WHERE order_id = $1;

-- name: UpdatePayment :one
UPDATE payments
SET order_id = $2, amount = $3, method = $4, reference = $5
WHERE id = $1
RETURNING *;

-- name: DeletePayment :exec
DELETE FROM payments
WHERE id = $1;

-- name: ListPayments :many
SELECT * FROM payments
ORDER BY paid_at DESC
LIMIT $1 OFFSET $2;

-- name: ListPaymentsByDateRange :many
SELECT * FROM payments
WHERE paid_at >= $1 AND paid_at < $2
ORDER BY paid_at DESC
LIMIT $3 OFFSET $4;

-- name: ListPaymentsByMethod :many
SELECT * FROM payments
WHERE method = $1
ORDER BY paid_at DESC
LIMIT $2 OFFSET $3;

-- name: GetPaymentsByOrderIDs :many
SELECT * FROM payments
WHERE order_id = ANY($1::int[])
ORDER BY paid_at DESC;

-- name: GetDailyPaymentSummary :many
SELECT 
    DATE(paid_at) as payment_date,
    method,
    COUNT(*) as payment_count,
    SUM(amount) as total_amount
FROM payments
WHERE paid_at >= $1 AND paid_at <= $2
GROUP BY DATE(paid_at), method
ORDER BY payment_date DESC, method;

-- name: GetMonthlyPaymentSummary :many
SELECT 
    DATE_TRUNC('month', paid_at) as payment_month,
    method,
    COUNT(*) as payment_count,
    SUM(amount) as total_amount
FROM payments
WHERE paid_at >= $1 AND paid_at <= $2
GROUP BY DATE_TRUNC('month', paid_at), method
ORDER BY payment_month DESC, method;

-- name: GetPaymentMethodStats :many
SELECT 
    method,
    COUNT(*) as payment_count,
    SUM(amount) as total_amount,
    AVG(amount) as average_amount
FROM payments
WHERE paid_at >= $1 AND paid_at <= $2
GROUP BY method
ORDER BY total_amount DESC;