create table if not exists playbooks (
	id serial primary key,
	name varchar(255) not null,
    team_id integer references teams(id),
	description text,
    formations string not null,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);