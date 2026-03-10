-- name: CreateFeed :one
INSERT INTO feeds (id, name, url, user_id, created_at, updated_at)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetAllFeedsWithUserName :many
SELECT fd.name as "feedName", fd.url as "feedUrl", usr.name as "userName"
FROM feeds fd JOIN users usr ON fd.user_id = usr.id;


-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (
    id,
    feed_id,
    user_id,
    created_at,
    updated_at
    )
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
INNER JOIN feeds ON feeds.id = inserted_feed_follow.feed_id  
INNER JOIN users ON users.id = inserted_feed_follow.user_id;

-- name: GetFeedIdByUrl :one
SELECT id
FROM feeds
WHERE url = $1;

-- name: GetFeedFollowsForUser :many
SELECT
    feeds.name AS feed_name
FROM feed_follows
INNER JOIN feeds ON feeds.id = feed_follows.feed_id  
INNER JOIN users ON users.id = feed_follows.user_id
WHERE users.name = $1;

-- name: UnfollowFeed :exec
DELETE
FROM feed_follows
where user_id = $1
and feed_id = $2;