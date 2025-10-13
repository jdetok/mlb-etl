-- INTAKE TABLES CREATION FOR MLB POSTGRES DB - 10/12/2025
-- these tables should accept the data fetched directly from the MLB API

-- FROM SCHEDULE ENDPOINT
create table if not exists intake.game_from_schedule (
    id integer primary key,
    guid uuid,
    gtype char(1),
    season char(4),
    start_time timestamptz,
    gdate date,
    state_abstract varchar(255),
    state char(1),
    state_desc varchar(255),
    status varchar(2), -- originally char(1)
    start_time_tbd boolean,
    gcode_abstract char(1),
    -- home/away fields begin here
    home_win boolean,
    away_win boolean,
    home_score smallint,
    away_score smallint,
    home_tm integer, -- home team ID
    away_tm integer, -- away team ID
    home_tm_name varchar(255),
    away_tm_name varchar(255),
    home_wins smallint,
    away_wins smallint,
    home_losses smallint,
    away_losses smallint,
    home_pct numeric(5,3),
    away_pct numeric(5,3),
    home_series_num smallint,
    away_series_num smallint,
    home_split_squad boolean,
    away_split_squad boolean,
    home_api_link varchar(255),
    away_api_link varchar(255),
    -- venue fields
    venue_id integer,
    venue varchar(255),
    venue_api_link varchar(255),
    -- other game fields
    tie boolean,
    day_type char(1),
    day_night varchar(10),
    description varchar(255), -- only present in playoffs
    season_display char(4),
    series_desc varchar(255),
    if_necessary char(1),
    if_necessary_desc varchar(255)
);

create table if not exists intake.team_detail (
    id integer primary key,
    name varchar(255),
    api_link varchar(255),
    season char(4),
    abbr varchar(5),
    team_cde varchar(5),
    team_name varchar(255),
    loc varchar(255),
    league_id integer,
    league varchar(255),
    league_api_link varchar(255),
    div_id integer,
    division varchar(255),
    div_api_link varchar(255),
    sport_id integer,
    sport varchar(255),
    sport_api_link varchar(255),
    short_name varchar(255),
    franchise varchar(255),
    club varchar(255),
    first_year char(4),
    file_cde varchar(10),
    all_star_status varchar(2555), -- verify type
    active boolean,
    ven_id integer,
    venue varchar(255),
    ven_api_link varchar(255),
    sven_id integer,
    spring_venue varchar(255),
    sven_api_link varchar(255),
    slg_id integer,
    spring_league varchar(255),
    slg_api_link varchar(255),
    slg_abbr varchar(5)
);