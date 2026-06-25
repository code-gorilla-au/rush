-- name: CreateGame :one
insert into games (name,
                   team_a,
                   team_b,
                   tournament_id,
                   results_log,
                   status,
                   rounds,
                   current_round)
values (?,
        ?,
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

-- name: UpdateGame :one
update games
set name          = ?,
    team_a        = ?,
    team_b        = ?,
    winner        = ?,
    status        = ?,
    results_log   = ?,
    rounds        = ?,
    current_round = ?,
    tournament_id = ?
where id = ?
returning *;