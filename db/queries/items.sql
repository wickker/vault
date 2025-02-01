-- name: CreateItem :one
INSERT INTO items (name, clerk_user_id)
VALUES ($1, $2)
RETURNING *;

-- name: ListItemsByUser :many
SELECT id, name
FROM items
WHERE clerk_user_id = $1
  AND deleted_at IS NULL
ORDER BY created_at DESC;