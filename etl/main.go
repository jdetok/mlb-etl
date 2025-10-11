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
	err := GetAndMakeDS[RespSchedule]("v1/schedule",
		[]Param{
			{Key: "sportId", Val: "1"},
			{Key: "season", Val: "2025"},
			{Key: "gameType", Val: "D"},
		},
	)
	if err != nil {
		ErrHndl(err)
	}
	err = GetAndMakeDS[RespTeams]("v1/teams", []Param{{Key: "158"}})
	if err != nil {
		ErrHndl(err)
	}
}

func GetAndMakeDS[T any](endpt string, params []Param) error {
	gr := HTTPGet{
		Base:     BASE,
		Endpoint: endpt,
		Params:   params,
	}

	js, err := gr.SendGetRequest()
	if err != nil {
		return err
	}
	rs, err := MakeDS[T](js)
	if err != nil {
		return err
	}
	fmt.Println(rs)
	return nil
}
