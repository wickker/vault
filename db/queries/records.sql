-- name: DeleteRecords :exec
UPDATE records
SET deleted_at = NOW()
WHERE item_id = $1;