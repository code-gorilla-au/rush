-- name: GetCoach :one
SELECT * FROM teams WHERE id = ?;

-- name: GetDefaultCoach :one
SELECT * FROM coaches WHERE is_default = true LIMIT 1;

-- name: SetDefaultCoach :exec
UPDATE coaches SET is_default = true WHERE id = ?;

-- name: ClearDefaultCoach :exec
UPDATE coaches SET is_default = false WHERE is_default = true;

-- name: GetCoaches :many
SELECT * FROM coaches;

-- name: CreateCoach :one
INSERT INTO coaches (name, is_default) VALUES (?, ?) RETURNING *;

-- name: DeleteCoach :exec
DELETE FROM coaches WHERE id = ?;

-- name: GetTeams :many
SELECT * FROM teams;

-- name: GetTeam :one
SELECT * FROM teams WHERE id = ?;

-- name: GetTeamByCoachID :one
SELECT * FROM teams WHERE coach_id = ? AND is_default = true LIMIT 1;

-- name: CreateTeam :exec
INSERT INTO teams (name, is_default, coach_id) VALUES (?, ?, ?) RETURNING *;

-- name: SetDefaultTeam :exec
UPDATE teams SET is_default = true WHERE id = ?;

-- name: ClearDefaultTeam :exec
UPDATE teams SET is_default = false WHERE is_default = true;

-- name: DeleteTeam :exec
DELETE FROM teams WHERE id = ?;

-- name: GetTeamMembers :many
SELECT * FROM players WHERE team_id = ?;

-- name: CreatePlayer :one
INSERT INTO players (name, team_id) VALUES (?,?) RETURNING *;

-- name: DeletePlayer :exec
DELETE FROM players WHERE id = ?;