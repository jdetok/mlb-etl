# JSON examples for public MLB API
## /schedule endpoint
- EXAMPLE: 2025 divisional search
    - `https://statsapi.mlb.com/api/v1/schedule?sportId=1&season=2025&gameType=D`
```json
{
    "copyright": "Copyright 2025 MLB Advanced Media, L.P.  Use of any content on this page acknowledges agreement to the terms posted here http://gdx.mlb.com/components/copyright.txt",
    "totalItems": 18,
    "totalEvents": 0,
    "totalGames": 18,
    "totalGamesInProgress": 0,
    "dates": [
        {
            "date": "2025-10-04",
            "totalItems": 4,
            "totalEvents": 0,
            "totalGames": 4,
            "totalGamesInProgress": 0,
            "games": [
                {
                    "gamePk": 813047,
                    "gameGuid": "390e6924-36dc-4464-ac57-772ad2c94731",
                    "link": "/api/v1.1/game/813047/feed/live",
                    "gameType": "D",
                    "season": "2025",
                    "gameDate": "2025-10-04T18:08:00Z",
                    "officialDate": "2025-10-04",
                    "status": {
                        "abstractGameState": "Final",
                        "codedGameState": "F",
                        "detailedState": "Final",
                        "statusCode": "F",
                        "startTimeTBD": false,
                        "abstractGameCode": "F"
                    },
                    "teams": {
                        "away": {
                            "leagueRecord": {
                                "wins": 0,
                                "losses": 1,
                                "pct": ".000"
                            },
                            "score": 3,
                            "team": {
                                "id": 112,
                                "name": "Chicago Cubs",
                                "link": "/api/v1/teams/112"
                            },
                            "isWinner": false,
                            "splitSquad": false,
                            "seriesNumber": 3
                        },
                        "home": {
                            "leagueRecord": {
                                "wins": 1,
                                "losses": 0,
                                "pct": "1.000"
                            },
                            "score": 9,
                            "team": {
                                "id": 158,
                                "name": "Milwaukee Brewers",
                                "link": "/api/v1/teams/158"
                            },
                            "isWinner": true,
                            "splitSquad": false,
                            "seriesNumber": 3
                        }
                    },
                    "venue": {
                        "id": 32,
                        "name": "American Family Field",
                        "link": "/api/v1/venues/32"
                    },
                    "content": {
                        "link": "/api/v1/game/813047/content"
                    },
                    "isTie": false,
                    "gameNumber": 1,
                    "publicFacing": true,
                    "doubleHeader": "N",
                    "gamedayType": "P",
                    "tiebreaker": "N",
                    "calendarEventID": "14-813047-2025-10-04",
                    "seasonDisplay": "2025",
                    "dayNight": "day",
                    "description": "NLDS 'A' Game 1",
                    "scheduledInnings": 9,
                    "reverseHomeAwayStatus": false,
                    "inningBreakLength": 175,
                    "gamesInSeries": 5,
                    "seriesGameNumber": 1,
                    "seriesDescription": "Division Series",
                    "recordSource": "S",
                    "ifNecessary": "N",
                    "ifNecessaryDescription": "Normal Game"
                },
            ]
        }
    ]
}

```