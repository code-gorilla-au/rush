-- name: CreatePlaybook :one
insert into playbooks (name,
                       description,
                       formations,
                       team_id)
VALUES (?,
        ?,
        ?,
        ?)
RETURNING *;

-- name: UpdatePlaybook :one
update playbooks
set name = ?,
    description = ?,
    formations = ?
where id = ?
returning *;

-- name: DeletePlaybook :exec
delete from playbooks where id = ?;

-- name: GetPlaybooksByTeam :many
select * from playbooks where team_id = ?;
