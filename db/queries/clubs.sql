-- name: GetAllClubs :many
SELECT * FROM clubs ORDER BY name;

-- name: GetClubByID :one
SELECT * FROM clubs WHERE id = ? LIMIT 1;

-- name: GetClubByName :one
SELECT * FROM clubs WHERE name = ? LIMIT 1;

-- name: CreateClub :one
INSERT INTO clubs (name, strength, background_color, foreground_color)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: DeleteClub :exec
DELETE FROM clubs WHERE id = ?;
