-- name: CreateOrderItem :one
INSERT INTO order_items (order_id, item_id, quantity, unit_price, notes, name)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetOrderItemByID :one
SELECT * FROM order_items
WHERE id = $1;

-- name: UpdateOrderItem :one
UPDATE order_items
SET 
    order_id = $2,
    item_id = $3,
    quantity = $4,
    unit_price = $5,
    notes = $6,
    name = $7,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteOrderItem :exec
DELETE FROM order_items
WHERE id = $1;

-- name: GetOrderItemsByOrderID :many
SELECT * FROM order_items
WHERE order_id = $1
ORDER BY created_at ASC;

-- name: DeleteOrderItemsByOrderID :exec
DELETE FROM order_items
WHERE order_id = $1;

-- name: GetOrderItemByOrderAndItem :one
SELECT * FROM order_items
WHERE order_id = $1 AND item_id = $2
LIMIT 1;