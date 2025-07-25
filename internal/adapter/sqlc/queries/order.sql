-- name: CreateOrder :one
INSERT INTO orders (table_id, status, notes, qrcode)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetOrderByID :one
SELECT * FROM orders
WHERE id = $1;

-- name: UpdateOrder :one
UPDATE orders
SET 
    table_id = $2,
    status = $3,
    notes = $4,
    closed_at = $5,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteOrder :exec
DELETE FROM orders
WHERE id = $1;

-- name: ListOrders :many
SELECT * FROM orders
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListOrdersByTable :many
SELECT * FROM orders
WHERE table_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetOpenOrderByTable :one
SELECT * FROM orders
WHERE table_id = $1 AND status = 'open'
ORDER BY created_at DESC
LIMIT 1;

-- name: ListOrdersByStatus :many
SELECT * FROM orders
WHERE status = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListOrdersByDateRange :many
SELECT * FROM orders
WHERE created_at >= $1 AND created_at <= $2
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

-- name: GetOrderByQRCode :one
SELECT * FROM orders
WHERE qrcode = $1;