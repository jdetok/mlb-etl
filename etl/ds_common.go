package main

import (
	"encoding/json"
	"fmt"
)

/*
GENERIC JSON TO GO STRUCT UNMARSHALER
creates a variable of the desired type, attempts to unmarshal the passed js
slice of bytes into that variable. returns pointer to the variable if successful
*/
func MakeDS[T any](js []byte) (*T, error) {
	var v T
	if err := json.Unmarshal(js, &v); err != nil {
		fmt.Println(err)
		return &v, err
	}
	return &v, nil
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
