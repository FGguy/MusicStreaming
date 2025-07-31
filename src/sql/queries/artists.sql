-- name: CreateArtist :one
INSERT INTO Artists (name, cover_art, album_count)
VALUES ($1, $2, $3) RETURNING *;

-- name: GetArtist :one
SELECT * FROM Artists
WHERE artist_id = $1 LIMIT 1;

-- name: GetArtists :many
SELECT * FROM Artists;