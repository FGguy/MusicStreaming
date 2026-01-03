-- name: CreateAdminUser :one
INSERT INTO Users (username, password, email, adminRole)
VALUES ($1, $2, $3, TRUE) 
ON CONFLICT (username) DO UPDATE SET username = EXCLUDED.username
RETURNING *;

-- name: CreateDefaultUser :one
INSERT INTO Users (username, password, email)
VALUES ($1, $2, $3) 
ON CONFLICT (username) DO UPDATE SET username = EXCLUDED.username
RETURNING *;

-- name: DeleteUser :one
DELETE FROM Users 
WHERE username = $1 RETURNING *;

-- name: GetUserByUsername :one
SELECT * FROM Users
WHERE username = $1 LIMIT 1;

-- name: GetUsers :many
SELECT * FROM Users;

-- name: ChangeUserPassword :one
UPDATE Users SET password = $2
WHERE username = $1 RETURNING *;

-- name: UpdateUser :one
UPDATE Users SET
    password = $2,
    email = $3,
    scrobblingEnabled = $4,
    ldapAuthenticated = $5,
    adminRole = $6,
    settingsRole = $7,
    streamRole = $8,
    jukeboxRole = $9,
    downloadRole = $10,
    uploadRole = $11,
    playlistRole = $12,
    coverArtRole = $13,
    commentRole = $14,
    podcastRole = $15,
    shareRole = $16,
    videoConversionRole = $17,
    musicFolderId = $18,
    maxBitRate = $19
WHERE username = $1 RETURNING *;