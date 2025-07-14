-- name: InsertUser :one
INSERT INTO Users (name, password)
VALUES ($1, $2) 
ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
RETURNING *;

-- name: DeleteUser :one
DELETE FROM Users 
WHERE name = $1 RETURNING *;

-- name: GetUserByName :one
SELECT * FROM Users
WHERE name = $1 LIMIT 1;

-- name: UpdateUserPassword :one
UPDATE Users SET password = $2
WHERE name = $1 RETURNING *;