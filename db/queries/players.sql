-- name: GetPlayersByClubID :many
SELECT * FROM players WHERE club_id = ? ORDER BY id;

-- name: GetPlayerByID :one
SELECT * FROM players WHERE id = ? LIMIT 1;

-- name: CreatePlayer :one
INSERT INTO players (club_id, name, quality)
VALUES (?, ?, ?)
RETURNING *;

-- name: DeletePlayer :exec
DELETE FROM players WHERE id = ?;
