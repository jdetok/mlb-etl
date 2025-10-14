// extract_t.go contains data structures to unmarshal response JSON into
package etl

import "time"

// used for team, venue, league, etc
type MLBObj struct {
	ID   uint16 `json:"id"`
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
