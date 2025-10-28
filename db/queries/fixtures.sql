-- name: GetAllFixtures :many
SELECT * FROM fixtures ORDER BY id;

-- name: GetFixtureByID :one
SELECT * FROM fixtures WHERE id = ? LIMIT 1;

-- name: GetFixturesByClubID :many
SELECT * FROM fixtures
WHERE home_team_id = ? OR away_team_id = ?
ORDER BY id;

-- name: GetUnplayedByClubID :many
SELECT f.id, f.gameweek, f.home_team_id, f.away_team_id, f.created_at
FROM fixtures f
LEFT JOIN matches m ON m.fixture_id = f.id AND m.is_completed = 1
WHERE (f.home_team_id = ?1 OR f.away_team_id = ?1)
  AND m.id IS NULL
ORDER BY f.gameweek, f.id;

-- name: CreateFixture :one
INSERT INTO fixtures (gameweek, home_team_id, away_team_id)
VALUES (?, ?, ?)
RETURNING *;

-- name: DeleteFixture :exec
DELETE FROM fixtures WHERE id = ?;
