// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: items.sql

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createItem = `-- name: CreateItem :one
INSERT INTO items (name, clerk_user_id, category_id)
VALUES ($1, $2, $3)
RETURNING id, name, clerk_user_id, created_at, updated_at, deleted_at, category_id
`

type CreateItemParams struct {
	Name        string
	ClerkUserID string
	CategoryID  pgtype.Int4
}

func (q *Queries) CreateItem(ctx context.Context, arg CreateItemParams) (Item, error) {
	row := q.db.QueryRow(ctx, createItem, arg.Name, arg.ClerkUserID, arg.CategoryID)
	var i Item
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.ClerkUserID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.CategoryID,
	)
	return i, err
}

const deleteItem = `-- name: DeleteItem :one
UPDATE items
SET deleted_at = NOW()
WHERE id = $1
AND clerk_user_id = $2
RETURNING id, name, clerk_user_id, created_at, updated_at, deleted_at, category_id
`

type DeleteItemParams struct {
	ID          int32
	ClerkUserID string
}

func (q *Queries) DeleteItem(ctx context.Context, arg DeleteItemParams) (Item, error) {
	row := q.db.QueryRow(ctx, deleteItem, arg.ID, arg.ClerkUserID)
	var i Item
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.ClerkUserID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.CategoryID,
	)
	return i, err
}

const getItem = `-- name: GetItem :one
SELECT id, name, clerk_user_id
FROM items
WHERE id = $1
AND deleted_at IS NULL
`

type GetItemRow struct {
	ID          int32
	Name        string
	ClerkUserID string
}

func (q *Queries) GetItem(ctx context.Context, id int32) (GetItemRow, error) {
	row := q.db.QueryRow(ctx, getItem, id)
	var i GetItemRow
	err := row.Scan(&i.ID, &i.Name, &i.ClerkUserID)
	return i, err
}

const listItemsByCategory = `-- name: ListItemsByCategory :many
SELECT id
FROM items
WHERE category_id = $1
AND clerk_user_id = $2
AND deleted_at IS NULL
`

type ListItemsByCategoryParams struct {
	CategoryID  pgtype.Int4
	ClerkUserID string
}

func (q *Queries) ListItemsByCategory(ctx context.Context, arg ListItemsByCategoryParams) ([]int32, error) {
	rows, err := q.db.Query(ctx, listItemsByCategory, arg.CategoryID, arg.ClerkUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int32
	for rows.Next() {
		var id int32
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listItemsByUser = `-- name: ListItemsByUser :many
SELECT id, name, created_at, category_id
FROM items
WHERE clerk_user_id = $1
AND (name ILIKE $2 OR $2 IS NULL)
AND (category_id = $3 OR $3 IS NULL)
AND deleted_at IS NULL
ORDER BY
CASE
WHEN $4::text = 'name_desc' THEN name
END DESC,
CASE
WHEN $4::text = 'created_at_desc' THEN created_at
END DESC,
CASE
WHEN $4::text = 'name_asc' THEN name
END ASC,
CASE
WHEN $4::text = 'created_at_asc' THEN created_at
END ASC
`

type ListItemsByUserParams struct {
	ClerkUserID string
	Name        pgtype.Text
	CategoryID  pgtype.Int4
	OrderBy     string
}

type ListItemsByUserRow struct {
	ID         int32
	Name       string
	CreatedAt  pgtype.Timestamp
	CategoryID pgtype.Int4
}

func (q *Queries) ListItemsByUser(ctx context.Context, arg ListItemsByUserParams) ([]ListItemsByUserRow, error) {
	rows, err := q.db.Query(ctx, listItemsByUser,
		arg.ClerkUserID,
		arg.Name,
		arg.CategoryID,
		arg.OrderBy,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListItemsByUserRow
	for rows.Next() {
		var i ListItemsByUserRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.CreatedAt,
			&i.CategoryID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateItem = `-- name: UpdateItem :one
UPDATE items
SET name = $1,
category_id = $2
WHERE id = $3
AND clerk_user_id = $4
RETURNING id, name, clerk_user_id, created_at, updated_at, deleted_at, category_id
`

type UpdateItemParams struct {
	Name        string
	CategoryID  pgtype.Int4
	ID          int32
	ClerkUserID string
}

func (q *Queries) UpdateItem(ctx context.Context, arg UpdateItemParams) (Item, error) {
	row := q.db.QueryRow(ctx, updateItem,
		arg.Name,
		arg.CategoryID,
		arg.ID,
		arg.ClerkUserID,
	)
	var i Item
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.ClerkUserID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.CategoryID,
	)
	return i, err
}
