-- name: CreateMenuItem :one
INSERT INTO menu_items (category_id, name, description, price, is_active)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetMenuItemByID :one
SELECT * FROM menu_items
WHERE id = $1;

-- name: UpdateMenuItem :one
UPDATE menu_items
SET 
    category_id = $2,
    name = $3,
    description = $4,
    price = $5,
    is_active = $6,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteMenuItem :exec
DELETE FROM menu_items
WHERE id = $1;

-- name: ListMenuItems :many
SELECT * FROM menu_items
WHERE is_active = true
ORDER BY name ASC
LIMIT $1 OFFSET $2;

-- name: ListMenuItemsByCategory :many
SELECT * FROM menu_items
WHERE category_id = $1 AND is_active = true
ORDER BY name ASC
LIMIT $2 OFFSET $3;

-- name: SearchMenuItems :many
SELECT * FROM menu_items
WHERE (name ILIKE '%' || $1 || '%' OR description ILIKE '%' || $1 || '%')
  AND is_active = true
ORDER BY name ASC
LIMIT $2 OFFSET $3;

-- name: GetMenuItemsByIDs :many
SELECT * FROM menu_items
WHERE id = ANY($1::int[])
  AND is_active = true
ORDER BY name ASC;

-- name: UpdateMenuItemStatus :one
UPDATE menu_items
SET 
    is_active = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: CountMenuItemsByCategory :one
SELECT COUNT(*) FROM menu_items
WHERE category_id = $1 AND is_active = true;