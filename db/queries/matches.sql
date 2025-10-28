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

-- name: CompleteMatch :exec
UPDATE matches
SET home_score = ?,
    away_score = ?,
    is_completed = 1,
    completed_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE fixture_id = ?;

-- name: GetCompletedMatches :many
SELECT * FROM matches WHERE is_completed = 1 ORDER BY completed_at DESC;
