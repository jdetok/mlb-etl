-- league schema (base player, team, lg, data) 

-- detail for american and national leagues
create table if not exists lg.lgdtl (
    lgid integer primary key,
    lg varchar(255)
);

-- season detail, include season id, team ids for ws winner, alcs, nlcs, etc
create table if not exists lg.szndtl (
    sznid varchar(4) primary key,
    szn varchar(255)
);

-- division detail
create table if not exists lg.dvdtl (
    dvid integer primary key,
    lgid integer references lg.lgdtl,
    dvn varchar(255)
);

-- team detail, no player data
create table if not exists lg.tmdtl (
    teamid integer primary key,
    teamcde varchar(10),
    abbr varchar(10),
    name varchar(255),
    loc varchar(255)
);

-- row by team by season
create table if not exists lg.tmszn (
    teamid integer references lg.tmdtl,
    sznid varchar(4) references lg.szndtl,
    lgid integer references lg.lgdtl,
    dvid integer references lg.dvdtl,
    primary key (teamid, sznid)
);

create table if not exists lg.psndtl (
    psnid integer primary key,
    psn varchar(255),
    psn_type varchar(255),
    abbr varchar(50)
);

-- player detail, nothing team related
create table if not exists lg.plrdtl (
    plrid integer primary key,
    name varchar(255),
    fname varchar(255),
    lname varchar(255),
    birhdate date,
    debut_date date,
    draft_year varchar(4),
    bat_hand varchar(10),
    pitch_hand varchar(10),
    strkz_top numeric(5, 3),
    strkz_btm numeric(5, 3)
);

-- season team roster, season id, team id, player id for each row
create table if not exists lg.tmroster (
    sznid varchar(4) references lg.szndtl,
    teamid integer references lg.tmdtl,
    dvid integer references lg.dvdtl,
    lgid integer references lg.lgdtl,
    plrid integer references lg.plrdtl,
    psnid integer references lg.psndtl, 
    jersey integer,
    primary key (sznid, teamid, plrid)
);