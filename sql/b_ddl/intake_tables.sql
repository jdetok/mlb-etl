-- INTAKE TABLES CREATION FOR MLB POSTGRES DB - 10/12/2025
-- these tables should accept the data fetched directly from the MLB API

create table intake.game_from_schedule (
    id integer primary key,
    guid uuid,
    gtype char(1),
    season char(4),
    start_time timestamptz,
    gdate date,
    state_abstract varchar(255),
    state char(1),
    state_desc varchar(255),
    status char(1),
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
)