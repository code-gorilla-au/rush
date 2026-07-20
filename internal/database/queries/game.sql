-- name: CreateGame :one
insert into games (name,
                   team_a,
                   team_b,
                   results_log,
                   status,
                   rounds,
                   current_round)
values (?,
        ?,
        ?,
        ?,
        'pending',
        ?,
        ?)
returning *;

-- name: GetGameByID :one
select *
from games
where id = ?;

-- name: StartGame :exec
update games
set status = 'running'
where id = ?;

-- name: UpdateGame :one
update games
set name          = ?,
    team_a        = ?,
    team_b        = ?,
    winner        = ?,
    status        = ?,
    results_log   = ?,
    rounds        = ?,
    current_round = ?
where id = ?
returning *;

-- name: ListCompletedGamesByTeam :many
select *
from games
where status = 'complete'
  and (team_a = ? or team_b = ?)
order by updated_at desc;

-- name: CreateTournament :one
insert into tournaments (name,
                         number_of_teams)
values (?,
        ?)
returning *;

-- name: CreateStage :one
insert into stages (name, tournament_id, status)
values (?, ?, ?)
returning *;

-- name: AllocateGameToStage :one
insert into stage_games (stage_id, game_id)
values (?, ?)
returning *;

-- name: UpdateStage :one
update stages
set name   = ?,
    status = ?
where id = ?
returning *;