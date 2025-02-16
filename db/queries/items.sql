-- name: CreateItem :one
INSERT INTO items (name, clerk_user_id)
VALUES ($1, $2)
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
SELECT id, name, created_at
FROM items
WHERE clerk_user_id = $1
AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: UpdateItem :one
UPDATE items
SET name = $1
WHERE id = $2
AND clerk_user_id = $3
RETURNING *;