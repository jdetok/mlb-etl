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

/*
create HTTP GET request from the passed endpoint and parameters, send the
request with an HTTP client and get the JSON response, unmarshal the JSON into
the struct passed as [T]
*/
func GetAndMakeDS[T any](endpt string, params []Param) (*T, error) {
	// create get request
	gr := HTTPGet{
		Base:     BASE,
		Endpoint: endpt,
		Params:   params,
	}

	// get JSON response
	js, err := gr.SendGetRequest()
	if err != nil {
		return nil, err
	}

	// create & return the data structure passed at [T] from JSON response
	ds, err := MakeDS[T](js)
	if err != nil {
		return nil, err
	}
	return ds, nil
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
	fmt.Println(schedule)

	// teams endpoint
	teams, err := GetAndMakeDS[RespTeams]("v1/teams", []Param{{Key: "158"}})
	if err != nil {
		ErrHndl(err)
	}
	fmt.Println(teams)
}
