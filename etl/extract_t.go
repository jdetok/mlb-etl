// extract_t.go contains data structures to unmarshal response JSON into
package etl

// TODO: teams/999/coaches/season=2020
// sports/1/players (this one will be by season)

import "time"

// used for team, venue, league, etc
type MLBObj struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`
}

// outer struct for schedule endpoint
type RespSchedule struct {
	Dates []MLBDate `json:"dates"`
}

// contains slice of games for a particular date
type MLBDate struct {
	DateStr      string    `json:"date"`
	TotalGames   uint8     `json:"totalGames"`
	TotalGamesIP uint8     `json:"totalGamesInProgress"` // use to skip date if > 0
	Games        []MLBGame `json:"games"`
}

// game data in schedule
type MLBGame struct {
	GID             uint64        `json:"gamePk"`
	GUID            string        `json:"gameGuid"`
	LiveFeed        string        `json:"link"`
	Type            string        `json:"gameType"`
	Season          string        `json:"season"`
	TmpDateTime     string        `json:"gameDate"`
	DateTime        time.Time     // convert TmpDateTime from api resp to go time
	DateStr         string        `json:"officialDate"`
	Status          MLBGameStatus `json:"status"`
	Teams           MLBGameTeams  `json:"teams"`
	Venue           MLBObj        `json:"venue"`
	IsTie           bool          `json:"isTie"`
	DayType         string        `json:"gameDayType"`
	DayNight        string        `json:"dayNight"`
	Description     string        `json:"description"` // only present in playoffs
	GamesInSeries   uint8         `json:"gamesInSeries"`
	GameNum         uint8         `json:"seriesGameNumber"`
	SeasonDisplay   string        `json:"seasonDisplay"`
	SeriesDesc      string        `json:"seriesDescription"`
	IfNecessary     string        `json:"ifNecessary"`
	IfNecessaryDesc string        `json:"ifNecessaryDescription"`
}

// game status in schedule - only care about final
type MLBGameStatus struct {
	AbstractState string `json:"abstractGameState"`
	StateCode     string `json:"codedGameState"`
	State         string `json:"detailedState"`
	Code          string `json:"statusCode"`
	StartTBD      bool   `json:"startTimeTBD"`
	AbstractCode  string `json:"abstractGameCode"`
}

// home and away teams (record, detail, etc) for schedule endpoint
type MLBGameTeams struct {
	Away MLBGameTeam `json:"away"`
	Home MLBGameTeam `json:"home"`
}

// team in schedule game
type MLBGameTeam struct {
	Record     MLBSeriesRecord `json:"leagueRecord"`
	Detail     MLBObj          `json:"team"`
	Win        bool            `json:"isWinner"`
	Score      uint8           `json:"score"`
	SplitSquad bool            `json:"splitSquad"`
	SeriesNum  uint8           `json:"seriesNumber"`
}

// team's record in the series a game is in
type MLBSeriesRecord struct {
	Wins   uint8  `json:"wins"`
	Losses uint8  `json:"losses"`
	PctStr string `json:"pct"`
	Pct    float64
}

// derived from teams endpoint, after RespSchedule struct exists
type RespTeams struct {
	CR    string       `json:"copyright"`
	Teams []TeamDetail `json:"teams"`
}

// contains details for team from calling teams/{teamId}
type TeamDetail struct {
	SpringLeague  MLBSpringLeague `json:"springLeague"`  //--
	AllStarSt     string          `json:"allStarStatus"` //--
	MLBObj                        // capture team id, name, link --
	Season        uint16          `json:"season"`          // --
	Venue         MLBObj          `json:"venue"`           //--
	SpringVenue   MLBObj          `json:"springVenue"`     //--
	Code          string          `json:"teamCode"`        // --
	FileCode      string          `json:"fileCode"`        //--
	Abbr          string          `json:"abbreviation"`    //--
	TeamName      string          `json:"teamName"`        //--
	Location      string          `json:"locationName"`    //--
	FirstYear     string          `json:"firstYearOfPlay"` //--
	League        MLBObj          `json:"league"`          //--
	Division      MLBObj          `json:"division"`        //--
	Sport         MLBObj          `json:"sport"`           //--
	ShortName     string          `json:"shortName"`       //--
	FranchiseName string          `json:"franchiseName"`   //--
	ClubName      string          `json:"clubName"`        //--
	Active        bool            `json:"active"`          //--
}

// spring league section of teams endpoint
type MLBSpringLeague struct {
	MLBObj
	Abbr string `json:"abbreviation"`
}

// PLAYER STRUCTS
type RespRoster struct {
	CR     string      `json:"copyright"`
	People []MLBPerson `json:"roster"`
}

type MLBPerson struct {
	Detail   MLBPlayerDetail `json:"person"`
	Jersey   string          `json:"jerseyNumber"`
	Position MLBPosition     `json:"position"`
	Status   MLBAttr         `json:"status"`
	TeamID   uint16          `json:"parentTeamId"`
}

// can't use MLBObj becaause json tag is fullName instead of name
// use this for most of /sports/player general stats, also can be used
// for only id fullname link
type MLBPlayerDetail struct {
	SPrID            uint64 // primary key, personidseasonid
	Season           string
	ID               uint64      `json:"id"`
	Name             string      `json:"fullName"`
	Link             string      `json:"link"`
	FName            string      `json:"firstName"`
	LName            string      `json:"lastName"`
	PrimNum          string      `json:"primaryNumber"`
	TmpBirthDay      string      `json:"birthDate"`
	BirthDay         time.Time   // convert TmpBirthday to dt
	Age              int16       `json:"currentAge"`
	BirthCity        string      `json:"birthCity"`
	BirthState       string      `json:"birthStateProvince"`
	BirthCountry     string      `json:"birthCountry"`
	Height           string      `json:"height"`
	Weight           uint16      `json:"weight"`
	Active           bool        `json:"active"`
	CurrentTeam      MLBObj      `json:"currentTeam"` // gets id and link, no team
	PrimPos          MLBPosition `json:"primaryPosition"`
	UseName          string      `json:"useName"`
	UseLName         string      `json:"useLastName"`
	MName            string      `json:"middleName"`
	BoxScoreName     string      `json:"boxscoreName"`
	Gender           string      `json:"gender"`
	IsPlayer         bool        `json:"isPlayer"`
	IsVerified       bool        `json:"isVerified"`
	DraftYear        uint16      `json:"draftYear"`
	TmpDebutDate     string      `json:"mlbDebutDate"`
	DebutDate        time.Time   // convert debut date
	BatSide          MLBAttr     `json:"batSide"`
	PitchHand        MLBAttr     `json:"pitchHand"`
	NameFL           string      `json:"nameFirstLast"`
	NameSlug         string      `json:"nameSlug"`
	FLName           string      `json:"firstLastName"`
	LFName           string      `json:"lastFirstName"`
	LIName           string      `json:"lastInitName"`
	FMLName          string      `json:"fullFMLName"`
	LMFName          string      `json:"fullLFMName"`
	StrikeZoneTop    float64     `json:"strikeZoneTop"`
	StrikeZoneBottom float64     `json:"strikeZoneBottom"`
}

type MLBPosition struct {
	Code string `json:"code"`
	Name string `json:"name"`
	Type string `json:"type"`
	Abbr string `json:"abbreviation"`
}

// convert from MLBPersonStatus to more general MLBAttr
type MLBAttr struct {
	Code string `json:"code"`
	Desc string `json:"description"`
}

// MLB PLAYERS /sports/1/players
type RespPlayers struct {
	Players []MLBPlayerDetail `json:"people"`
	Season  string            // use to make playerseason id for db
}

// box score stats for a specific game
// game/gameId/boxscore
type RespBoxscore struct {
	Teams MLBBoxScoreTeams `json:"teams"`
}

type MLBBoxScoreTeams struct {
	Away MLBTeamBoxScore `json:"away"`
	Home MLBTeamBoxScore `json:"home"`
}

type MLBTeamBoxScore struct {
	TeamDtl   TeamDetail                   `json:"team"`
	TeamStats MLBBoxStats                  `json:"teamStats"`
	Players   map[string]MLBPlayerBoxScore `json:"players"`
}

type MLBPlayerBoxScore struct {
	Person       MLBBoxPlayerDtl `json:"person"`
	Jersey       string          `json:"jerseyNumber"`
	Status       MLBAttr         `json:"status"`
	TeamID       uint16          `json:"parentTeamId"`
	BattingOrder string          `json:"battingOrder"`
	Stats        MLBBoxStats     `json:"stats"`
}

type MLBBoxPlayerDtl struct {
	ID      uint64 `json:"id"`
	Name    string `json:"fullName"`
	Link    string `json:"link"`
	BoxName string `json:"boxscoreName"`
}

// should work for team and player stats
type MLBBoxStats struct {
	Batting  MLBBattingStats  `json:"batting"`
	Pitching MLBPitchingStats `json:"pitching"`
	Fielding MLBFieldingStats `json:"fielding"`
}

type MLBBattingStats struct {
	Summary        string `json:"summary"`
	GP             uint16 `json:"gamesPlayed"`
	FlyOuts        uint16 `json:"flyOuts"`
	GndOuts        uint16 `json:"groundOuts"`
	Airouts        uint16 `json:"airOuts"`
	Runs           uint16 `json:"runs"`
	Doubles        uint16 `json:"doubles"`
	Triples        uint16 `json:"triples"`
	HomeRuns       uint16 `json:"homeRuns"`
	StrikeOuts     uint16 `json:"strikeOuts"`
	BaseOnBalls    uint16 `json:"baseOnBalls"`
	IntnWalks      uint16 `json:"intentionalWalks"`
	Hits           uint16 `json:"hits"`
	HitByPitch     uint16 `json:"hitByPitch"`
	AtBats         uint16 `json:"atBats"`
	CaughtStealing uint16 `json:"caughtStealing"`
	StolenBases    uint16 `json:"stolenBases"`
	StolenBasesPct string `json:"stolenBasePercentage"`
	GndIntoDP      uint16 `json:"groundIntoDoublePlay"`
	PlateApps      uint16 `json:"plateAppearances"`
	TotalBases     uint16 `json:"totalBases"`
	RBI            uint16 `json:"rbi"`
	LeftOnBase     uint16 `json:"leftOnBase"`
	SacBunts       uint16 `json:"sacBunts"`
	SacFlies       uint16 `json:"sacFlies"`
	CatchersIntr   uint16 `json:"catchersInterference"`
	Pickoffs       uint16 `json:"pickoffs"`
	AtBatPerHR     string `json:"atBatsPerHomeRun"`
	PopOuts        uint16 `json:"popOuts"`
	LineOuts       uint16 `json:"lineOuts"`
}

type MLBPitchingStats struct {
	MLBBattingStats          // a bunch of duplicate fields in the batting struct
	GamesStarted      uint16 `json:"gamesStarted"`
	NumPitches        uint16 `json:"numberOfPitches"`
	InningsPitched    uint16 `json:"inningsPitched"`
	Wins              uint16 `json:"wins"`
	Losses            uint16 `json:"losses"`
	Saves             uint16 `json:"saves"`
	SaveOpps          uint16 `json:"saveOpportunities"`
	Holds             uint16 `json:"holds"`
	BlownSaves        uint16 `json:"blownSaves"`
	EarnedRuns        uint16 `json:"earnedRuns"`
	BattersFaced      uint16 `json:"battersFaced"`
	Outs              uint16 `json:"outs"`
	GamesPitched      uint16 `json:"gamesPitched"`
	CompleteGames     uint16 `json:"completeGames"`
	Shutouts          uint16 `json:"shutouts"`
	PitchesThrown     uint16 `json:"pitchesThrown"`
	Balls             uint16 `json:"balls"`
	Strikes           uint16 `json:"strikes"`
	StrikePct         string `json:"strikePercentage"`
	HitBatsmen        uint16 `json:"hitBatsmen"`
	Balks             uint16 `json:"balks"`
	WildPitches       uint16 `json:"wildPitches"`
	GamesFinished     uint16 `json:"gamesFinished"`
	RunsPer9          string `json:"runsScoredPer9"`
	HRPer9            string `json:"homeRunsPer9"`
	InhrRunners       uint16 `json:"inheritedRunners"`
	InhrRunnersScored uint16 `json:"inheritedRunnersScored"`
	PassedBall        uint16 `json:"passedBall"`
}

type MLBFieldingStats struct {
	GamesStarted      uint16 `json:"gamesStarted"`
	CaughtStealing    uint16 `json:"caughtStealing"`
	StolenBases       uint16 `json:"stolenBases"`
	StolenBasePct     string `json:"stolenBasePercentage"`
	CaughtStealingPct string `json:"caughtStealingPercentage"`
	Assists           uint16 `json:"assists"`
	PutOuts           uint16 `json:"putOuts"`
	Errors            uint16 `json:"errors"`
	Chances           uint16 `json:"chances"`
	Fielding          string `json:"fielding"`
	PassedBall        uint16 `json:"passedBall"`
	Pickoffs          uint16 `json:"pickoffs"`
}
