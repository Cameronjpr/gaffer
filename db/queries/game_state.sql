-- name: GetAllGameState :many
SELECT * FROM game_states ORDER BY id;

-- name: GetMostRecentGameState :one
SELECT * FROM game_states ORDER BY updated_at DESC LIMIT 1;

-- name: CreateGameState :one
INSERT INTO game_states (selected_club_id, manager_name, created_at)
VALUES (?, ?, NOW())
RETURNING *;

-- name: UpdateGameState :one
UPDATE game_states SET selected_club_id = ?, manager_name = ?, updated_at = NOW()
WHERE id = ?
RETURNING *;

-- name: DeleteGameState :exec
DELETE FROM game_states WHERE id = ?;

-- name: DeleteAllGameStates :exec
DELETE FROM game_states;
