
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
    status varchar(255),
    results_log string not null,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

