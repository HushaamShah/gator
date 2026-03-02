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
