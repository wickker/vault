-- name: ListCategoriesByUser :many
SELECT id, name, color
FROM categories
WHERE clerk_user_id = $1
AND deleted_at IS NULL;
