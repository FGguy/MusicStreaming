-- name: CreateCover :one
INSERT INTO Covers (cover_id, path)
VALUES ($1, $2) RETURNING *;

-- name: GetCover :one
SELECT * FROM Covers
WHERE cover_id = $1 LIMIT 1;