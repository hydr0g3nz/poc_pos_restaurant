-- name: CreateTable :one
INSERT INTO tables (table_number, qr_code, seating, is_active)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetTableByID :one
SELECT * FROM tables
WHERE id = $1;

-- name: GetTableByNumber :one
SELECT * FROM tables
WHERE table_number = $1;

-- name: GetTableByQRCode :one
SELECT * FROM tables
WHERE qr_code = $1;

-- name: UpdateTable :one
UPDATE tables
SET 
    table_number = $2,
    qr_code = $3,
    seating = $4,
    is_active = $5
WHERE id = $1
RETURNING *;

-- name: DeleteTable :exec
DELETE FROM tables
WHERE id = $1;

-- name: ListTables :many
SELECT * FROM tables
ORDER BY table_number ASC;

-- name: CheckTableHasOrders :one
SELECT EXISTS(
    SELECT 1 FROM orders 
    WHERE table_id = $1
) as has_orders;