-- name: CreatePlaybook :one
insert into playbooks (name,
                       description,
                       formations)
VALUES (?,
        ?,
        ?)
RETURNING *;

-- name: UpdatePlaybookFormations :one
update playbooks
set formations = ?
where id = ?
returning *;

-- name: GetPlaybooksByTeam :many
select * from playbooks where team_id = ?;
