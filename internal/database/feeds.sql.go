// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: feeds.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createFeed = `-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING id, created_at, updated_at, name, url, user_id
`

type CreateFeedParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Url       string
	UserID    uuid.UUID
}

func (q *Queries) CreateFeed(ctx context.Context, arg CreateFeedParams) (Feed, error) {
	row := q.db.QueryRowContext(ctx, createFeed,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Name,
		arg.Url,
		arg.UserID,
	)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Url,
		&i.UserID,
	)
	return i, err
}

const returnAllFeeds = `-- name: ReturnAllFeeds :many
SELECT feeds.name as feedname, feeds.url as url, users.name as username
FROM feeds
INNER JOIN users ON users.id = feeds.user_id
`

type ReturnAllFeedsRow struct {
	Feedname string
	Url      string
	Username string
}

func (q *Queries) ReturnAllFeeds(ctx context.Context) ([]ReturnAllFeedsRow, error) {
	rows, err := q.db.QueryContext(ctx, returnAllFeeds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ReturnAllFeedsRow
	for rows.Next() {
		var i ReturnAllFeedsRow
		if err := rows.Scan(&i.Feedname, &i.Url, &i.Username); err != nil {
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
