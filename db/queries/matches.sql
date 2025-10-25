-- name: GetMatchByID :one
SELECT * FROM matches WHERE id = ? LIMIT 1;

-- name: GetMatchByFixtureID :one
SELECT * FROM matches WHERE fixture_id = ? LIMIT 1;

-- name: CreateMatch :one
INSERT INTO matches (
    fixture_id,
    current_minute,
    current_half,
    home_score,
    away_score,
    active_zone,
    home_attacking_direction
)
VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateMatch :exec
UPDATE matches
SET current_minute = ?,
    current_half = ?,
    home_score = ?,
    away_score = ?,
    active_zone = ?,
    home_attacking_direction = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: DeleteMatch :exec
DELETE FROM matches WHERE id = ?;
