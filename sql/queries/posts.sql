-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, desription, published_at, feed_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
RETURNING *;


-- name: GetPostsForUser :many
SELECT * FROM posts
WHERE feed_id in (
    SELECT feed_id
    FROM feed_follows
    WHERE user_id = (
        SELECT id
        FROM users
        WHERE name = $1
    )
)
ORDER BY published_at DESC
LIMIT $2;

