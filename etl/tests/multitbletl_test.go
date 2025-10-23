package etl

import (
	"fmt"
	"testing"

	"github.com/jdetok/mlb-etl/etl"
	"github.com/jdetok/mlb-etl/pgresd"
)

func TestMultiTblETL(t *testing.T) {
	// season, cause it'll be needed where it's called
	season := "2025"
	gameId := "777933"
	db, err := pgresd.ConnectTestDB("../../.env")
	if err != nil {
		t.Errorf("failed to connect to database | %v\n", err)
	}

	metl := etl.MakeMultiTableETL(nil, &etl.RespBoxscore{},
		"v1/game", []etl.Param{{Key: gameId}, {Key: "boxscore"}},
		[]etl.PGTarget{
			{PGSchema: "intake", PGTable: "tbtg", PGPKey: "teamid, gameid"},
			{PGSchema: "intake", PGTable: "tptg", PGPKey: "teamid, gameid"},
			{PGSchema: "intake", PGTable: "tfdg", PGPKey: "teamid, gameid"},
			{PGSchema: "intake", PGTable: "pbtg", PGPKey: "plrid, gameid"},
			{PGSchema: "intake", PGTable: "pptg", PGPKey: "plrid, gameid"},
			{PGSchema: "intake", PGTable: "pfdg", PGPKey: "plrid, gameid"},
		},
	)
	if err := metl.ExtractData(); err != nil {
		t.Errorf("failed extracting data\n%v", err)
	}

	metl.Dataset.(*etl.RespBoxscore).SetSharedVals(season, gameId)
	metl.Dataset.CleanTempFields()
	tableSets := metl.Dataset.SliceInsertRows()[0]

	// DO NOT DELETE
	for i, pgt := range metl.PGTargets {
		fmt.Printf("%v:\n%v\n++++++++++++\n\n", pgt.PGTable, tableSets[i])
		cols, err := pgresd.ColumnsInTable(db, pgt.PGTable)
		if err != nil {
			t.Errorf("failed to make InSt | %v\n", err)
		}
		rows := tableSets[i].([][]any)
		fmt.Println(rows)

		metl.InSt = pgresd.MakeInsert(pgt.PGSchema, pgt.PGTable, pgt.PGPKey,
			cols, rows)

		if err := metl.InSt.InsertFast(db, &metl.RowCount); err != nil {
			t.Errorf("failed to insert %d rows into %s\n%v",
				len(rows), pgt.PGTable, err)
		}

	}

}
