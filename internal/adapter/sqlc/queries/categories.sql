-- name: CreateCategory :one
INSERT INTO categories (name)
VALUES ($1)
RETURNING *;

-- name: GetCategoryByID :one
SELECT * FROM categories
WHERE id = $1;

-- name: GetCategoryByName :one
SELECT * FROM categories
WHERE name = $1;

-- name: UpdateCategory :one
UPDATE categories
SET 
    name = $2
WHERE id = $1
RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM categories
WHERE id = $1;

-- name: ListCategories :many
SELECT * FROM categories
ORDER BY name ASC;

-- name: CheckCategoryHasMenuItems :one
SELECT EXISTS(
    SELECT 1 FROM menu_items 
    WHERE category_id = $1
) as has_items;