// entry point
package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jdetok/golib/pgresd"
)

// MLB stats api base url
const BASE string = "https://statsapi.mlb.com/api"

// super quick error handling for testing, replace later
func ErrHndl(err error) {
	fmt.Println("an error occured: killing program")
	log.Fatal(err)
}
func ConnectDB() (*sql.DB, error) {
	pg := pgresd.GetEnvPG()
	pg.MakeConnStr()
	db, err := pg.Conn()
	if err != nil {
		return nil, err
	}
	return db, nil
}
func main() {
	// database connection
	db, err := ConnectDB()
	if err != nil {
		ErrHndl(err)
	}
	// fmt.Println("database setup:", db.Stats().OpenConnections)
	// cols, err := pgresd.ColumnsInTable(db, "game_from_schedule")
	// if err != nil {
	// 	ErrHndl(err)
	// }
	// fmt.Println(cols)
	// schedule endpoint
	schedule, err := GetAndMakeDS[RespSchedule]("v1/schedule",
		[]Param{
			{Key: "sportId", Val: "1"},
			{Key: "season", Val: "2025"},
			{Key: "gameType", Val: "R"},
		},
	)
	if err != nil {
		ErrHndl(err)
	}
	schedule.CleanGamesData()

	// fmt.Println(schedule)
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
