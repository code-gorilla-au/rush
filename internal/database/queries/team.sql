-- name: GetCoach :one
SELECT * FROM teams WHERE id = ?;

-- name: GetCoaches :many
SELECT * FROM coaches;

-- name: GetTeams :many
SELECT * FROM teams;

-- name: GetTeam :one
SELECT * FROM teams WHERE id = ?;

-- name: GetTeamMembers :many
SELECT * FROM players WHERE team_id = ?;

