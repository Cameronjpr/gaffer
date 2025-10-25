-- name: GetEventsByMatchID :many
SELECT * FROM match_events
WHERE match_id = ?
ORDER BY minute, id;

-- name: CreateMatchEvent :one
INSERT INTO match_events (
    match_id,
    event_type,
    minute,
    team_side,
    player_name
)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: DeleteMatchEvents :exec
DELETE FROM match_events WHERE match_id = ?;
