-- name: Insert :one
INSERT INTO movies (title, year, runtime, genres)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: Get :one
SELECT *
FROM movies
WHERE id = $1;
