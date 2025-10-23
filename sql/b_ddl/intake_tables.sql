-- INTAKE TABLES CREATION FOR MLB POSTGRES DB - 10/12/2025
-- these tables should accept the data fetched directly from the MLB API

create table if not exists log.log (
    logid uuid primary key default gen_random_uuid(),
    prj varchar(255),
    msg text,
    ltime timestamptz,
    ltstr varchar(20),
    caller varchar(255),
    err text,
    rc bigint
);

-- FROM SCHEDULE ENDPOINT
create table if not exists intake.game_from_schedule (
    gameid integer primary key,
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


create index idx_gfs_season on intake.game_from_schedule(season);
create index idx_gfs_gdate on intake.game_from_schedule(gdate);
create index idx_gfs_gtype on intake.game_from_schedule(gtype);
create index idx_gfs_home on intake.game_from_schedule(home_tm);
create index idx_gfs_away on intake.game_from_schedule(away_tm);

create table if not exists intake.team_detail (
    teamid integer primary key,
    name varchar(255),
    api_link varchar(255),
    season char(4),
    abbr varchar(10),
    team_cde varchar(10),
    team_name varchar(255),
    loc varchar(255),
    lgid integer,
    league varchar(255),
    league_api_link varchar(255),
    dvid integer,
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
    all_star_status varchar(10), -- verify type
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
    slg_abbr varchar(10)
);
create index idx_team_detail_season on intake.team_detail(season);
create index idx_team_detail_active on intake.team_detail(active);
create index idx_team_detail_abbr on intake.team_detail(abbr);
create index idx_team_detail_dvid on intake.team_detail(dvid);
create index idx_team_detail_lgid on intake.team_detail(lgid);

-- needs to be playerid seasonid
-- person season id, use with coaches too
-- splayer for season player
create table if not exists intake.splayer (
    sprid bigint primary key,
    season varchar(4),
    plrid integer,
    name varchar(255),
    api_link varchar(255),
    fname varchar(255),
    lname varchar(255),
    prim_jnum varchar(10),
    birth_date date,
    age smallint,
    birth_city varchar(255),
    birth_state varchar(255),
    birth_country varchar(255),
    height varchar(50),
    weight smallint,
    active boolean,
    teamid integer,
    team_api_link varchar(255),
    posn_cde varchar(10),
    posn varchar(50),
    posn_type varchar(50),
    posn_abbr varchar(10),
    use_name varchar(255),
    use_lname varchar(255),
    mname varchar(255),
    box_name varchar(255),
    gender varchar(50),
    is_player boolean,
    is_verif boolean,
    draft_year integer,
    debut_date date,
    bat_cde varchar(10),
    bat_desc varchar(255),
    pitch_cde varchar(10),
    pitch_desc varchar(255),
    namefl varchar(255),
    name_slug varchar(255),
    flname varchar(255),
    lfname varchar(255),
    liname varchar(255),
    fmlname varchar(255),
    lmfname varchar(255),
    strike_zone_top numeric(5, 3),
    strike_zone_btm numeric(5, 3)
); 
create index idx_splayer_season on intake.splayer(season);
create index idx_splayer_teamid on intake.splayer(teamid);
create index idx_splayer_active on intake.splayer(active);
create index idx_splayer_posn_id on intake.splayer(posn_cde);
-- team batting stats
create table if not exists intake.tbtg (
    teamid integer not null,
    gameid integer not null,
    season varchar(4),
    flyouts smallint,
    groundouts smallint,
    airouts smallint,
    doubles smallint,
    triples smallint,
    homeruns smallint,
    strkouts smallint,
    base_on_balls smallint,
    intnwalks smallint,
    hits smallint,
    hit_by_pitch smallint,
    avg numeric(5, 3),
    atbats smallint,
    obp numeric(5, 3),
    slg numeric(5, 3),
    ops numeric(5, 3),
    caught_stl smallint,
    stl_bases smallint,
    stl_base_pct numeric(5, 3),
    gnd_into_dp smallint, 
    gnd_into_tp smallint,
    plate_appr smallint,
    total_bases smallint,
    rbi smallint,
    left_on_base smallint,
    sacbunts smallint,
    sacflies smallint,
    catcherintf smallint,
    pickoffs smallint,
    ab_per_hr numeric(5, 3),
    popouts smallint,
    lineouts smallint,
    primary key (teamid, gameid)
);

create index idx_tbtg_season on intake.tbtg(season);
create index idx_tbtg_teamid on intake.tbtg(teamid);

-- team pitching stats
create table if not exists intake.tptg (
    teamid integer not null,
    gameid integer not null,
    season varchar(4),
    flyouts smallint,
    groundouts smallint,
    airouts smallint,
    doubles smallint,
    triples smallint,
    homeruns smallint,
    strkouts smallint,
    base_on_balls smallint,
    intnwalks smallint,
    hits smallint,
    hit_by_pitch smallint,
    avg numeric(5, 3),
    atbats smallint,
    obp numeric(5, 3),
    caught_stl smallint,
    stl_bases smallint,
    stl_base_pct numeric(5, 3),
    caught_stl_pct numeric(5, 3),
    num_pitches smallint,
    era numeric(5, 3),
    innings numeric(5, 3),
    sav_opps smallint,
    earned_runs smallint,
    whip numeric(5, 3),
    batters_faced smallint,
    complete_games smallint,
    shutouts smallint,
    pitches_thrown smallint,
    balls smallint,
    strikes smallint,
    strike_pct numeric(5, 3),
    hit_batter smallint,
    balks smallint,
    wild_pitches smallint,
    pickoffs smallint,
    groundouts_to_airouts numeric(5, 3),
    rbi smallint,
    pitches_per_inning numeric(5, 3),
    runs_scored_per9 numeric(5, 3),
    home_runs_per9 numeric(5, 3),
    inht_runners smallint,
    inht_runners_scored smallint,
    catchers_intf smallint,
    sacbunts smallint, 
    sacflies smallint,
    passed_ball smallint,
    popouts smallint,
    lineouts smallint,
    primary key (teamid, gameid)
);

create index idx_tptg_season on intake.tptg(season);
create index idx_tptg_teamid on intake.tptg(teamid);

-- team fielding stats
create table if not exists intake.tfdg (
    teamid integer not null,
    gameid integer not null,
    season varchar(4),
    caught_stl smallint,
    stl_bases smallint,
    stl_base_pct numeric(5, 3), 
    caught_stl_pct numeric(5, 3),
    assists smallint,
    putouts smallint,
    erors smallint,
    chances smallint,
    passed_ball smallint,
    pickoffs smallint,
    primary key (teamid, gameid)
);

create index idx_tfdg_season on intake.tfdg(season);
create index idx_tfdg_teamid on intake.tfdg(teamid);

-- player batting stats
create table if not exists intake.pbtg (
    plrid integer not null,
    teamid integer not null,
    gameid integer not null,
    season varchar(4),
    summary varchar(255),
    gp smallint,
    flyouts smallint,
    groundouts smallint,
    airouts smallint,
    doubles smallint,
    triples smallint,
    homeruns smallint,
    strkouts smallint,
    base_on_balls smallint,
    intnwalks smallint,
    hits smallint,
    hit_by_pitch smallint,
    atbats smallint,
    caught_stl smallint,
    stl_bases smallint,
    stl_base_pct numeric(5, 3),
    gnd_into_dp smallint, 
    gnd_into_tp smallint,
    plate_appr smallint,
    total_bases smallint,
    rbi smallint,
    left_on_base smallint,
    sacbunts smallint,
    sacflies smallint,
    catcherintf smallint,
    pickoffs smallint,
    ab_per_hr numeric(5, 3),
    popouts smallint,
    lineouts smallint,
    primary key (plrid, gameid)
);

create index idx_pbtg_season on intake.pbtg(season);
create index idx_pbtg_teamid on intake.pbtg(teamid);

-- player pitching stats
create table if not exists intake.pptg (
    plrid integer not null,
    teamid integer not null,
    gameid integer not null,
    season varchar(4),
    summary varchar(255),
    gp smallint,
    flyouts smallint,
    groundouts smallint,
    airouts smallint,
    doubles smallint,
    triples smallint,
    homeruns smallint,
    strkouts smallint,
    base_on_balls smallint,
    intnwalks smallint,
    hits smallint,
    hit_by_pitch smallint,
    avg numeric(5, 3),
    atbats smallint,
    obp numeric(5, 3),
    caught_stl smallint,
    stl_bases smallint,
    stl_base_pct numeric(5, 3),
    caught_stl_pct numeric(5, 3),
    num_pitches smallint,
    era numeric(5, 3),
    innings numeric(5, 3),
    sav_opps smallint,
    holds smallint,
    blown_saves smallint,
    earned_runs smallint,
    whip numeric(5, 3),
    batters_faced smallint,
    outs smallint,
    complete_games smallint,
    shutouts smallint,
    pitches_thrown smallint,
    balls smallint,
    strikes smallint,
    strike_pct numeric(5, 3),
    hit_batter smallint,
    balks smallint,
    wild_pitches smallint,
    pickoffs smallint,
    groundouts_to_airouts numeric(5, 3),
    rbi smallint,
    winpct numeric(5, 3),
    pitches_per_inning numeric(5, 3),
    games_finished smallint,
    so_walk_ratio numeric(5, 3),
    so_per9 numeric(5, 3),
    walks_per9 numeric(5, 3),
    hits_per9 numeric(5, 3),
    runs_scored_per9 numeric(5, 3),
    home_runs_per9 numeric(5, 3),
    inht_runners smallint,
    inht_runners_scored smallint,
    catchers_intf smallint,
    sacbunts smallint, 
    sacflies smallint,
    passed_ball smallint,
    popouts smallint,
    lineouts smallint,
    primary key (plrid, gameid)
);

create index idx_pptg_season on intake.pptg(season);
create index idx_pptg_teamid on intake.pptg(teamid);

-- player fielding stats
create table if not exists intake.pfdg (
    plrid integer not null,
    teamid integer not null,
    gameid integer not null,
    season varchar(4),
    caught_stl smallint,
    stl_bases smallint,
    stl_base_pct numeric(5, 3), 
    caught_stl_pct numeric(5, 3),
    assists smallint,
    putouts smallint,
    erors smallint,
    chances smallint,
    fielding numeric(5, 3),
    passed_ball smallint,
    pickoffs smallint,
    primary key (plrid, gameid)
);

create index idx_pfdg_season on intake.pfdg(season);
create index idx_pfdg_teamid on intake.pfdg(teamid);