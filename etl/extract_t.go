package main

/* extract_t.go
this file should contain data structures to unmarshal response JSON into
*/

import "time"

type RespSchedule struct {
	Dates []MLBDate `json:"dates"`
}

type MLBDate struct {
	DateStr      string    `json:"date"`
	TotalGames   uint8     `json:"totalGames"`
	TotalGamesIP uint8     `json:"totalGamesInProgress"` // use to skip date if > 0
	Games        []MLBGame `json:"games"`
}

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
	Description     string        `json:"description"`
	NumInSeries     uint8         `json:"gameInSeries"`
	SeasonDisplay   string        `json:"seasonDisplay"`
	SeriesDesc      string        `json:"seriesDescription"`
	IfNecessary     string        `json:"ifNecessary"`
	IfNecessaryDesc string        `json:"ifNecessaryDescription"`
}

type MLBGameStatus struct {
	AbstractState string `json:"abstractGameState"`
	StateCode     string `json:"codedGameState"`
	State         string `json:"detailedState"`
	Code          string `json:"statusCode"`
	StartTBD      bool   `json:"startTimeTBD"`
	AbstractCode  string `json:"abstractGameCode"`
}

type MLBGameTeams struct {
	Away MLBGameTeam `json:"away"`
	Home MLBGameTeam `json:"home"`
}

type MLBGameTeam struct {
	Record     MLBSeriesRecord `json:"leagueRecord"`
	Detail     MLBObj          `json:"team"`
	SplitSquad bool            `json:"splitSquad"`
	SeriesNum  uint8           `json:"seriesNumber"`
}

type MLBSeriesRecord struct {
	Wins   uint8  `json:"wins"`
	Losses uint8  `json:"losses"`
	Pct    string `json:"pct"`
}

// derived from schedule endpoint
type MLBObj struct {
	ID   uint16 `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`
}

// derived from teams endpoint, after RespSchedule struct exists
type RespTeams struct {
	CR    string       `json:"copyright"`
	Teams []TeamDetail `json:"teams"`
}

type TeamDetail struct {
	SpringLeague  MLBSpringLeague `json:"springLeague"`
	AllStarSt     string          `json:"alStarStatus"`
	MLBObj                        // capture team id, name, link
	Season        uint16          `json:"season"`
	Venue         MLBObj          `json:"venue"`
	SpringVenue   MLBObj          `json:"springVenue"`
	Code          string          `json:"teamCode"`
	FileCode      string          `json:"fileCode"`
	Abbr          string          `json:"abbreviation"`
	TeamName      string          `json:"teamName"`
	Location      string          `json:"locationName"`
	FirstYear     string          `json:"firstYearOfPlay"`
	League        MLBObj          `json:"league"`
	Division      MLBObj          `json:"division"`
	Sport         MLBObj          `json:"sport"`
	ShortName     string          `json:"shortName"`
	FranchiseName string          `json:"franchiseName"`
	ClubName      string          `json:"clubName"`
	Active        bool            `json:"active"`
}

type MLBSpringLeague struct {
	MLBObj
	Abbr string `json:"abbreviation"`
}
