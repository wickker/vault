-- name: DeleteRecords :many
UPDATE records
SET deleted_at = NOW()
WHERE item_id = $1
RETURNING *;

-- name: ListRecordsByItemId :many
SELECT id, name, value
FROM records
where deleted_at IS NULL
AND item_id = $1;

-- name: CreateRecord :one
INSERT INTO records (name, value, item_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateRecord :one
UPDATE records
SET name = $1,
value = $2
WHERE id = $3
RETURNING *;

-- name: DeleteRecord :one
UPDATE records
SET deleted_at = NOW()
WHERE id = $1
RETURNING *;

-- name: GetRecordUserID :one
SELECT items.clerk_user_id
FROM records
INNER JOIN items on items.id = records.item_id
WHERE records.id = $1
AND records.deleted_at IS NULL;
