-- name: CreateAlbum :one
INSERT INTO Albums (artist_id, name, cover_art, song_count, created, duration, artist)
VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: GetAlbum :one
SELECT * FROM Albums
WHERE album_id = $1 LIMIT 1;

-- name: GetAlbums :many
SELECT * FROM Albums
WHERE artist_id = $1;