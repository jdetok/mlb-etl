package etl

import (
	"fmt"
	"strconv"
	"sync"
	"testing"

	"github.com/jdetok/golib/pgresd"
	"github.com/jdetok/mlb-etl/etl"
)

func TestETL(t *testing.T) {
	endpt := "v1/sports"
	// sch := "intake"
	// table := "person"
	// pkey := "id"
	params := []etl.Param{{Key: "1"}, {Key: "players"}}
	ds, err := etl.GetAndMakeDS[etl.RespPlayers](endpt, params)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(ds)
}

// test full etl process for schedule endpoint
func TestScheduleETL(t *testing.T) {
	schedule, err := etl.GetAndMakeDS[etl.RespSchedule]("v1/schedule",
		[]etl.Param{
			{Key: "sportId", Val: "1"},
			{Key: "season", Val: "1980"},
			{Key: "gameType", Val: "R"},
		},
	)
	if err != nil {
		t.Errorf("failed to get schedule response | %v\n", err)
	}
	if err := schedule.CleanTempFields(); err != nil {
		t.Errorf("failed to clean games data | %v\n", err)
	}
	// create database connection
	db, err := pgresd.ConnectTestDB("../../.env")
	if err != nil {
		t.Errorf("failed to connect to database | %v\n", err)
	}
	if err := schedule.InsertGames(db); err != nil {
		t.Errorf("failed to insert records in database | %v\n", err)
	}
}

// attempt above etl process using new interface
func TestETLInterface(t *testing.T) {
	// create database connection
	db, err := pgresd.ConnectTestDB("../../.env")
	if err != nil {
		t.Errorf("failed to connect to database | %v\n", err)
	}

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
				t.Error(err)
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
				t.Error(err)
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
				t.Error(err)
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
		t.Error(err)
	}
	rc += pe.RowCount

	fmt.Println("FINAL COUNT: ", rc)
}

func TestPlayersETL(t *testing.T) {
	db, err := pgresd.ConnectTestDB("../../.env")
	if err != nil {
		t.Errorf("failed to connect to database | %v\n", err)
	}
	// pass already made struct to record season
	// enables prsid primary key

	var startYr int = 2004
	var endYr int = 2025
	// var dir int = 1

	for i := range endYr - startYr {
		// sznStr := strconv.Itoa(endYr - i)
		// direction by whivh is bigger
		go func() {
			var pl etl.RespPlayers
			pl.Season = strconv.Itoa(endYr - i)
			ple := etl.MakeETL(
				&pl,
				"intake",
				"splayer",
				"sprid",
				"v1/sports", // sports/1/players?season=2025
				[]etl.Param{{Key: "1"}, {Key: "players"}, {Key: "season", Val: pl.Season}})

			if err := ple.RunFullETL(db); err != nil {
				t.Error(err)
			}

			// fmt.Println(pl.Players)
			fmt.Println(pl.Players[0].SPrID)
		}()
	}
}
