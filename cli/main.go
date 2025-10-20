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

	// database connection
	db, err := pgresd.ConnectDB()
	if err != nil {
		lg.Log("couldn't connect to db", err, 0)
	}
	db.SetMaxOpenConns(200)
	db.SetMaxIdleConns(200)
	lg.DB = db

	// ETL years
	var startYr int = 1990
	var endYr int = 2025

	// mutex and waitgroup for safe concurrency
	var mu sync.Mutex
	var wg sync.WaitGroup

	// global row count
	var rc int64 = 0

	// only allow maxcon goroutines at once (semaphore)
	maxcon := 10
	sem := make(chan struct{}, maxcon)

	// COLLECT ERRS
	var allErrs []error
	var errMu sync.Mutex

	// loop through each season, run ETLs
	for i := 0; i < endYr-startYr; i++ {
		sem <- struct{}{}
		wg.Add(1)
		go func(i int, rc *int64) {
			defer wg.Done()
			defer func() { <-sem }() // clear one spot in sem
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
				errMu.Lock()
				allErrs = append(allErrs, err)
				errMu.Unlock()
			}

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
				errMu.Lock()
				allErrs = append(allErrs, err)
				errMu.Unlock()
			}

			lg.Log(fmt.Sprintf("done with players etl for %s", pl.Season), nil, *rc)

			mu.Lock()
			*rc += ple.RowCount
			mu.Unlock()

			// TEAMS ETL
			te := etl.MakeETL(&etl.RespTeams{},
				"intake", "team_detail", "id", "v1/teams",
				[]etl.Param{{Key: "season", Val: szn}}, &lg)

			if err := te.RunFullETL(db); err != nil {
				lg.Log("error with teams endpoint", err, *rc)
				errMu.Lock()
				allErrs = append(allErrs, err)
				errMu.Unlock()

			}

			lg.Log(fmt.Sprintf("done with teams etl for %s", szn), nil, *rc)

			mu.Lock()
			*rc += te.RowCount
			mu.Unlock()

		}(i, &rc)
	}
	wg.Wait()

	if len(allErrs) > 0 {
		fmt.Printf("errors occured: %d\n", len(allErrs))
	} else {
		fmt.Println("no errors")
	}

	lg.Log(fmt.Sprintf("FINAL ROW COUNT: %d", rc), nil, rc)
}
