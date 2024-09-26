-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING *
) 
SELECT 
    inserted_feed_follow.*,
    feeds.name AS feed_name,
    users.name AS user_name
FROM inserted_feed_follow
INNER JOIN users 
ON users.id = inserted_feed_follow.user_id
INNER JOIN feeds 
ON feeds.id = inserted_feed_follow.feed_id;


-- name: GetFeedFollowsForUser :many
SELECT f.Name as feed_name, u.Name as user_name
FROM feed_follows ff
INNER JOIN users u
ON u.id = ff.user_id
INNER JOIN feeds f
ON f.id = ff.feed_id
WHERE u.name = $1;


-- name: Unfollow :exec
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
);