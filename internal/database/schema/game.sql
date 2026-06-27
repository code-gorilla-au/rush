
create table if not exists tournaments (
    id integer primary key autoincrement,
    name text not null,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

create table if not exists games (
    id integer primary key autoincrement,
    name text not null,
    tournament_id integer references tournaments(id),
    team_a integer references teams(id),
    team_b integer references teams(id),
    winner integer references teams(id),
    status varchar(255) not null,
    rounds text not null,
    current_round integer not null default 0,
    results_log text not null,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

