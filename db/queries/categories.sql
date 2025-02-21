-- name: ListCategoriesByUser :many
SELECT id, name, color
FROM categories
WHERE clerk_user_id = $1
AND deleted_at IS NULL
ORDER BY name;

-- name: CreateCategory :one
INSERT INTO categories (name, color, clerk_user_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateCategory :one
UPDATE categories
SET name = $1,
color = $2
WHERE id = $3
AND clerk_user_id = $4
RETURNING *;

-- name: DeleteCategory :one
UPDATE categories
SET deleted_at = NOW()
WHERE id = $1
AND clerk_user_id = $2
RETURNING *;