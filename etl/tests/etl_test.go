package etl

import (
	"testing"

	"github.com/jdetok/golib/pgresd"
	"github.com/jdetok/mlb-etl/etl"
)

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

	// SCHEDULE ENDPOINT TEST
	// schema | table | primary key | endpoint | endpoint parameters
	e := etl.MakeETL(&etl.RespSchedule{},
		"intake", "game_from_schedule", "id", "v1/schedule",
		[]etl.Param{
			{Key: "sportId", Val: "1"},
			{Key: "season", Val: "1995"},
			{Key: "gameType", Val: "R"},
		},
	)

	if err := e.RunFullETL(db); err != nil {
		t.Error(err)
	}

	// TEAMS ETL
	te := etl.MakeETL(&etl.RespTeams{},
		"intake", "team_detail", "id", "v1/teams", []etl.Param{{}})

	if err := te.RunFullETL(db); err != nil {
		t.Error(err)
	}

	// TODO: FINISH PLAYER ETL
	pe := etl.MakeETL(&etl.RespPeople{},
		"intake", "person", "id", "v1/teams", []etl.Param{{Key: "138"}, {Key: "roster"}})

	if err := pe.RunFullETL(db); err != nil {
		t.Error(err)
	}

}
