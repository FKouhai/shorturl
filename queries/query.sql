-- name: GetUrl :one
SELECT * FROM urls
WHERE id = ? LIMIT 1;

-- name: GetUrls :many
SELECT * FROM urls
ORDER BY name;

-- name: CreateUrl :one
INSERT INTO urls (
  name
) VALUES (
  ?
)
RETURNING *;

-- name: UpdateUrl :exec
UPDATE urls
set name = ?
WHERE id = ?;

-- name: DeleteUrl :exec
DELETE FROM urls
WHERE id = ?;
