package main

import (
	"fmt"
	"log"
)

// MLB stats api base url
const BASE string = "https://statsapi.mlb.com/api"

// super quick error handling for testing, replace later
func ErrHndl(err error) {
	fmt.Println("an error occured: killing program")
	log.Fatal(err)
}

func main() {
	// schedule endpoint
	schedule, err := GetAndMakeDS[RespSchedule]("v1/schedule",
		[]Param{
			{Key: "sportId", Val: "1"},
			{Key: "season", Val: "2025"},
			{Key: "gameType", Val: "D"},
		},
	)
	if err != nil {
		ErrHndl(err)
	}
	schedule.GameDatesToDT()
	// fmt.Println(schedule)

	// teams endpoint
	teams, err := GetAndMakeDS[RespTeams]("v1/teams", []Param{{Key: "158"}})
	if err != nil {
		ErrHndl(err)
	}
	fmt.Println(teams)
}
