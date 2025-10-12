# Public MLB API Documentation
## gameType parameter
- using in /schedule endpoint
- R for regular season
- S for spring training
- documentation says P for post season, NOT accurate. have to speciy playoff round
    - F for fist round (wild card)
    - D for divisional champ
    - L for league champ
    - W for world series
## sportId parameter
- 1 for MLB

### seriesNumber field (within home/away objects)
- seems to be the count of series that team has played in so far that year

# /game endpoint
- /game/{gamePk}/linescore seems to give number of innings and final score infos