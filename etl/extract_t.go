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
	ID               uint64      `json:"id"`
	Name             string      `json:"fullName"`
	Link             string      `json:"link"`
	FName            string      `json:"firstName"`
	LName            string      `json:"lastName"`
	PrimNum          string      `json:"primaryNumber"`
	TmpBirthDay      string      `json:"birthDate"`
	BirthDay         time.Time   // convert TmpBirthday to dt
	Age              uint16      `json:"currentAge"`
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
