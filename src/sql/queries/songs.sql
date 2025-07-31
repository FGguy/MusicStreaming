-- name: CreateSong :one
INSERT INTO Songs (album_id, title, album, artist, is_dir, cover_art, created, duration, bit_rate, size, suffix, content_type, is_video)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING *;

-- name: GetSong :one
SELECT * FROM Songs
WHERE song_id = $1 LIMIT 1;

-- name: GetSongs :many
SELECT * FROM Songs
WHERE album_id = $1;