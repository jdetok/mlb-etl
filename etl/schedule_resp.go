package main

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
