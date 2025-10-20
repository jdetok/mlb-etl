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
    lgid integer references lg.lgdtl(lgid),
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
    teamid integer references lg.tmdtl(teamid),
    sznid varchar(4) references lg.szndtl(sznid),
    lgid integer references lg.lgdtl(lgid),
    dvid integer references lg.dvdtl(dvid),
    primary key (teamid, sznid)
);

-- player detail, nothing team related
create table if not exists lg.plrdtl (
    plrid integer primary key,
    name varchar(255),
    fname varchar(255),
    lname varchar(255),
    birhdate date,
    debut_date date,
    draft_year varchar(4)
);

-- season team roster, season id, team id, player id for each row
create table if not exists lg.tmroster (
    sznid varchar(4) references lg.szndtl(sznid),
    teamid integer references lg.tmdtl(teamid),
    plrid integer references lg.plrdtl(plrid),
    dvid integer references lg.dvdtl(dvid),
    lgid integer references lg.lgdtl(lgid),
    primary key (sznid, teamid, plrid)
);