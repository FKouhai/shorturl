// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package db

import (
	"context"
)

const createUrl = `-- name: CreateUrl :one
INSERT INTO urls (
  name
) VALUES (
  ?
)
RETURNING id, name
`

func (q *Queries) CreateUrl(ctx context.Context, name string) (Url, error) {
	row := q.db.QueryRowContext(ctx, createUrl, name)
	var i Url
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const deleteUrl = `-- name: DeleteUrl :exec
DELETE FROM urls
WHERE id = ?
`

func (q *Queries) DeleteUrl(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteUrl, id)
	return err
}

const getUrlData = `-- name: GetUrlData :one
SELECT id, name FROM urls
WHERE id = ? LIMIT 1
`

func (q *Queries) GetUrlData(ctx context.Context, id int64) (Url, error) {
	row := q.db.QueryRowContext(ctx, getUrlData, id)
	var i Url
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const getUrlId = `-- name: GetUrlId :one
SELECT id FROM urls
WHERE name = ? LIMIT 1
`

func (q *Queries) GetUrlId(ctx context.Context, name string) (int64, error) {
	row := q.db.QueryRowContext(ctx, getUrlId, name)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const getUrls = `-- name: GetUrls :many
SELECT id, name FROM urls
ORDER BY name
`

func (q *Queries) GetUrls(ctx context.Context) ([]Url, error) {
	rows, err := q.db.QueryContext(ctx, getUrls)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Url
	for rows.Next() {
		var i Url
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateUrl = `-- name: UpdateUrl :exec
UPDATE urls
set name = ?
WHERE id = ?
`

type UpdateUrlParams struct {
	Name string
	ID   int64
}

func (q *Queries) UpdateUrl(ctx context.Context, arg UpdateUrlParams) error {
	_, err := q.db.ExecContext(ctx, updateUrl, arg.Name, arg.ID)
	return err
}
