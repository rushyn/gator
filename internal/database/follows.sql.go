// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: follows.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createFeedFollow = `-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING id, created_at, updated_at, user_id, feed_id
) 
SELECT 
    inserted_feed_follow.id, inserted_feed_follow.created_at, inserted_feed_follow.updated_at, inserted_feed_follow.user_id, inserted_feed_follow.feed_id,
    feeds.name AS feed_name,
    users.name AS user_name
FROM inserted_feed_follow
INNER JOIN users 
ON users.id = inserted_feed_follow.user_id
INNER JOIN feeds 
ON feeds.id = inserted_feed_follow.feed_id
`

type CreateFeedFollowParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	FeedID    uuid.UUID
}

type CreateFeedFollowRow struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	FeedID    uuid.UUID
	FeedName  string
	UserName  string
}

func (q *Queries) CreateFeedFollow(ctx context.Context, arg CreateFeedFollowParams) (CreateFeedFollowRow, error) {
	row := q.db.QueryRowContext(ctx, createFeedFollow,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.UserID,
		arg.FeedID,
	)
	var i CreateFeedFollowRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.FeedID,
		&i.FeedName,
		&i.UserName,
	)
	return i, err
}

const getFeedFollowsForUser = `-- name: GetFeedFollowsForUser :many
SELECT f.Name as feed_name, u.Name as user_name
FROM feed_follows ff
INNER JOIN users u
ON u.id = ff.user_id
INNER JOIN feeds f
ON f.id = ff.feed_id
WHERE u.name = $1
`

type GetFeedFollowsForUserRow struct {
	FeedName string
	UserName string
}

func (q *Queries) GetFeedFollowsForUser(ctx context.Context, name string) ([]GetFeedFollowsForUserRow, error) {
	rows, err := q.db.QueryContext(ctx, getFeedFollowsForUser, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFeedFollowsForUserRow
	for rows.Next() {
		var i GetFeedFollowsForUserRow
		if err := rows.Scan(&i.FeedName, &i.UserName); err != nil {
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

const unfollow = `-- name: Unfollow :exec
DELETE 
FROM feed_follows ff
WHERE ff.id = (
        SELECT ff.id
        FROM feed_follows ff
        INNER JOIN users u
        ON u.id = ff.user_id
        INNER JOIN feeds f
        ON f.id = ff.feed_id
        WHERE u.name = $1 and f.url = $2
)
`

type UnfollowParams struct {
	Name string
	Url  string
}

func (q *Queries) Unfollow(ctx context.Context, arg UnfollowParams) error {
	_, err := q.db.ExecContext(ctx, unfollow, arg.Name, arg.Url)
	return err
}
