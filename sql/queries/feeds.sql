-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: ReturnAllFeeds :many
SELECT feeds.name as feedname, feeds.url as url, users.name as username
FROM feeds
INNER JOIN users ON users.id = feeds.user_id;


-- name: GetFeed :one
SELECT * FROM feeds
WHERE url = $1;
