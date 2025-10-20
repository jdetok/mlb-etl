// entry point
package main

import (
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/jdetok/golib/pgresd"
	"github.com/jdetok/mlb-etl/etl"
	"github.com/jdetok/mlb-etl/logd"
)

// super quick error handling for testing, replace later
func ErrHndl(err error) {
	fmt.Println("an error occured: killing program")
	log.Fatal(err)
}

func main() {
	lg := logd.Logder{Prj: "mlb-etl"}
	lg.Log("starting", nil, 0)
	db, err := pgresd.ConnectDB()
	if err != nil {
		lg.Log("couldn't connect to db", err, 0)
	}
	db.SetMaxOpenConns(500)
	db.SetMaxIdleConns(50)
	lg.DB = db
	var startYr int = 1950
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
				}, &lg,
			)
			if err := e.RunFullETL(db); err != nil {
				lg.Log("schedule endpoint failed", err, *rc)
			}

			// fmt.Println("schedule rows:", e.RowCount)

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
					{Key: "season", Val: pl.Season}}, &lg)

			if err := ple.RunFullETL(db); err != nil {
				lg.Log("error with players endpoint", err, *rc)
				// fmt.Println(err)
			}

			lg.Log(fmt.Sprintf("done with players etl for %s", pl.Season), nil, *rc)

			mu.Lock()
			*rc += ple.RowCount
			mu.Unlock()

			// fmt.Println("player rows:", ple.RowCount)
			// TEAMS ETL
			te := etl.MakeETL(&etl.RespTeams{},
				"intake", "team_detail", "id", "v1/teams",
				[]etl.Param{{Key: "season", Val: szn}}, &lg)

			if err := te.RunFullETL(db); err != nil {
				// fmt.Println(err)
				lg.Log("error with teams endpoint", err, *rc)
			}

			lg.Log(fmt.Sprintf("done with teams etl for %s", szn), nil, *rc)

			mu.Lock()
			*rc += te.RowCount
			mu.Unlock()

			// fmt.Println("team rows:", te.RowCount)

		}(i, &rc)
	}
	wg.Wait()
	// TODO: FINISH PLAYER ETL
	pe := etl.MakeETL(&etl.RespRoster{},
		"intake", "person", "id", "v1/teams",
		[]etl.Param{{Key: "138"}, {Key: "roster"}}, &lg)

	if err := pe.RunFullETL(db); err != nil {
		// fmt.Println(err)
		lg.Log("error with roster endpoint", err, rc)
	}
	rc += pe.RowCount

	lg.Log(fmt.Sprintf("FINAL ROW COUNT: %d", rc), nil, rc)
	// fmt.Println("FINAL COUNT: ", rc)
}
