-- name: CreateItem :one
INSERT INTO items (name, clerk_user_id, category_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: DeleteItem :one
UPDATE items
SET deleted_at = NOW()
WHERE id = $1
AND clerk_user_id = $2
RETURNING *;

-- name: GetItem :one
SELECT id, name, clerk_user_id
FROM items
WHERE id = $1
AND deleted_at IS NULL;

-- name: ListItemsByUser :many
SELECT id, name, created_at, category_id
FROM items
WHERE clerk_user_id = $1
AND (name ILIKE sqlc.narg('name') OR sqlc.narg('name') IS NULL)
AND (category_id = sqlc.narg('category_id') OR sqlc.narg('category_id') IS NULL)
AND deleted_at IS NULL
ORDER BY
CASE
WHEN @order_by::text = 'name_desc' THEN name
END DESC,
CASE
WHEN @order_by::text = 'created_at_desc' THEN created_at
END DESC,
CASE
WHEN @order_by::text = 'name_asc' THEN name
END ASC,
CASE
WHEN @order_by::text = 'created_at_asc' THEN created_at
END ASC;

-- name: UpdateItem :one
UPDATE items
SET name = $1,
category_id = $2
WHERE id = $3
AND clerk_user_id = $4
RETURNING *;