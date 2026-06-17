create table if not exists coaches (
	id serial primary key,
	name varchar(255) not null
);

create table if not exists teams (
    id serial primary key,
    name varchar(255) not null,
    coach_id integer references coach(id)
)

create table if not exists players (
	id serial primary key,
	name varchar(255) not null,
	team_id integer references team(id)
);