// entry point
package main

import (
	"fmt"
	"log"

	"github.com/jdetok/golib/pgresd"
	"github.com/jdetok/mlb-etl/etl"
)

// super quick error handling for testing, replace later
func ErrHndl(err error) {
	fmt.Println("an error occured: killing program")
	log.Fatal(err)
}

func main() {
	db, err := pgresd.ConnectDB()
	if err != nil {
		ErrHndl(err)
	}
	schedule, err := etl.GetAndMakeDS[etl.RespSchedule]("v1/schedule",
		[]etl.Param{
			{Key: "sportId", Val: "1"},
			{Key: "season", Val: "1971"},
			{Key: "gameType", Val: "R"},
		},
	)
	if err != nil {
		ErrHndl(err)
	}
	schedule.CleanTempFields()

	if err := schedule.InsertGames(db); err != nil {
		ErrHndl(err)
	}

	// teams endpoint
	// teams, err := GetAndMakeDS[RespTeams]("v1/teams", []Param{{Key: "158"}})
	// if err != nil {
	// 	ErrHndl(err)
	// }
	// fmt.Println(teams)
}
