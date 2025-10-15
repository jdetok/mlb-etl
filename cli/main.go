// entry point
package main

import (
	"fmt"
	"log"
	"strconv"
	"sync"

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
	db.SetMaxOpenConns(200)
	db.SetMaxIdleConns(50)

	var startYr int = 1970
	var endYr int = 2025
	var rc int64 = 0
	var mu sync.Mutex
	var wg sync.WaitGroup
	for i := 0; i < endYr-startYr; i++ {
		wg.Add(1)
		go func(i int, rc *int64) {
			defer wg.Done()
			szn := strconv.Itoa(endYr - i)

			// SCHEDULE ENDPOINT TEST
			// schema | table | primary key | endpoint | endpoint parameters
			e := etl.MakeETL(&etl.RespSchedule{},
				"intake", "game_from_schedule", "id", "v1/schedule",
				[]etl.Param{
					{Key: "sportId", Val: "1"},
					{Key: "season", Val: szn},
					{Key: "gameType", Val: "R"},
				},
			)
			if err := e.RunFullETL(db); err != nil {
				fmt.Println(err)
			}

			fmt.Println("schedule rows:", e.RowCount)

			mu.Lock()
			*rc += e.RowCount
			mu.Unlock()

			var pl etl.RespPlayers
			// sports/1/players?season=2025
			pl.Season = szn
			ple := etl.MakeETL(&pl, "intake", "splayer", "sprid", "v1/sports",
				[]etl.Param{
					{Key: "1"},
					{Key: "players"},
					{Key: "season", Val: pl.Season}})

			if err := ple.RunFullETL(db); err != nil {
				fmt.Println(err)
			}

			mu.Lock()
			*rc += ple.RowCount
			mu.Unlock()

			fmt.Println("player rows:", ple.RowCount)
			// TEAMS ETL
			te := etl.MakeETL(&etl.RespTeams{},
				"intake", "team_detail", "id", "v1/teams",
				[]etl.Param{{Key: "season", Val: szn}})

			if err := te.RunFullETL(db); err != nil {
				fmt.Println(err)
			}

			mu.Lock()
			*rc += te.RowCount
			mu.Unlock()

			fmt.Println("team rows:", te.RowCount)

		}(i, &rc)
	}
	wg.Wait()
	// TODO: FINISH PLAYER ETL
	pe := etl.MakeETL(&etl.RespRoster{},
		"intake", "person", "id", "v1/teams", []etl.Param{{Key: "138"}, {Key: "roster"}})

	if err := pe.RunFullETL(db); err != nil {
		fmt.Println(err)
	}
	rc += pe.RowCount

	fmt.Println("FINAL COUNT: ", rc)
}
