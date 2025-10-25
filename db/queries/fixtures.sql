-- name: GetAllFixtures :many
SELECT * FROM fixtures ORDER BY id;

-- name: GetFixtureByID :one
SELECT * FROM fixtures WHERE id = ? LIMIT 1;

-- name: GetFixturesByClubID :many
SELECT * FROM fixtures
WHERE home_team_id = ? OR away_team_id = ?
ORDER BY id;

-- name: CreateFixture :one
INSERT INTO fixtures (home_team_id, away_team_id)
VALUES (?, ?)
RETURNING *;

-- name: DeleteFixture :exec
DELETE FROM fixtures WHERE id = ?;
