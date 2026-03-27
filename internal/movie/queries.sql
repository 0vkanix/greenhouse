-- name: Insert :one
INSERT INTO movies (title, year, runtime, genres)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at, version;

-- name: Get :one
SELECT *
FROM movies
WHERE id = $1;

-- name: Update :one
UPDATE movies
SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
WHERE id = $5
RETURNING version;

-- name: Delete :execresult
DELETE FROM movies WHERE id = $1;
