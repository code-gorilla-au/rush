-- name: GetCoach :one
SELECT * FROM teams WHERE id = ?;

-- name: GetCoaches :many
SELECT * FROM coaches;

-- name: CreateCoach :exec
INSERT INTO coaches (name) VALUES (?) RETURNING *;

-- name: DeleteCoach :exec
DELETE FROM coaches WHERE id = ?;

-- name: GetTeams :many
SELECT * FROM teams;

-- name: GetTeam :one
SELECT * FROM teams WHERE id = ?;

-- name: CreateTeam :exec
INSERT INTO teams (name, is_default, coach_id) VALUES (?, ?, ?) RETURNING *;

-- name: SetDefaultTeam :exec
UPDATE teams SET is_default = false WHERE id = ?;

-- name: ClearDefaultTeam :exec
UPDATE teams SET is_default = false WHERE is_default = true;

-- name: DeleteTeam :exec
DELETE FROM teams WHERE id = ?;

-- name: GetTeamMembers :many
SELECT * FROM players WHERE team_id = ?;

-- name: CreatePlayer :exec
INSERT INTO players (name, team_id) VALUES (?,?) RETURNING *;

-- name: DeletePlayer :exec
DELETE FROM players WHERE id = ?;