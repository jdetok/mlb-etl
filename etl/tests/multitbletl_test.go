package etl

import (
	"fmt"
	"testing"

	"github.com/jdetok/mlb-etl/etl"
	"github.com/jdetok/mlb-etl/logd"
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

// start with TestMultiTblETL, but call all games for a season
func TestSeasonBoxETL(t *testing.T) {
	lg := logd.Logder{Prj: "mlb-etl-test"}
	lg.Log("starting", nil, nil)

	// database connection
	db, err := pgresd.ConnectTestDB("../../.env")
	if err != nil {
		t.Errorf("failed to connect to database | %v\n", err)
	}
	lg.DB = db
	start := 2024
	end := 2025

	// max number of goroutines
	maxcon := 20
	total_rows := int64(0)

	betl := etl.BatchETL{
		StartSzn:  start,
		EndSzn:    end,
		MaxGoRtns: maxcon,
		RowCount:  0,
		TotalRC:   &total_rows,
		Log:       &lg,
	}

	betl.LoadManyBoxScoreETL(db, &lg)
}

func TestGetGs(t *testing.T) {
	db, err := pgresd.ConnectTestDB("../../.env")
	if err != nil {
		t.Errorf("failed to connect to database | %v\n", err)
	}
	start := 2024
	end := 2025

	// max number of goroutines
	maxcon := 20
	total_rows := int64(0)

	b := etl.BatchETL{
		StartSzn:  start,
		EndSzn:    end,
		MaxGoRtns: maxcon,
		RowCount:  0,
		TotalRC:   &total_rows,
		// Log:       &lg,
	}

	lg := logd.Logder{Prj: "mlb-etl-test"}

	if err := b.GetGameIDsSeasons(db, &lg); err != nil {
		// lg.Log(fmt.Sprintf("failed to get game ids:\n%v\n", err), err, &b.RowCount)
		t.Error(err)
	}
}

// if err := betl.LoadManyBoxScoreETL(db, &lg); err != nil {
// 	t.Errorf("failed to load many:\n%v", err)
// }

// for i := range end - start {
// 	if err := betl.GetGameIdSzn(db, end-i); err != nil {
// 		t.Errorf("failed to get game ids:\n%v\n", err)
// 	}
// }
// betl.ChunkGameIDs()
// fmt.Println(betl.ChunkedGameIDs)
// for i, chunk := range betl.ChunkedGameIDs {
// 	fmt.Printf("chunk %d/%d\n%v\n", i+1, len(betl.ChunkedGameIDs), chunk)
// }
// fmt.Println(betl.ChunkedGameIDs[0])
// for _, gmap := range betl.ChunkedGameIDs[0] {
// 	var gameId string
// 	var season string
// 	for gId, szn := range gmap {
// 		gameId = strconv.FormatUint(gId, 10)
// 		season = szn
// 	}
// 	metl := etl.MakeMultiTableETL(nil, &etl.RespBoxscore{},
// 		"v1/game", []etl.Param{
// 			{Key: gameId}, {Key: "boxscore"},
// 		}, []etl.PGTarget{
// 			{PGSchema: "intake", PGTable: "tbtg", PGPKey: "teamid, gameid"},
// 			{PGSchema: "intake", PGTable: "tptg", PGPKey: "teamid, gameid"},
// 			{PGSchema: "intake", PGTable: "tfdg", PGPKey: "teamid, gameid"},
// 			{PGSchema: "intake", PGTable: "pbtg", PGPKey: "plrid, gameid"},
// 			{PGSchema: "intake", PGTable: "pptg", PGPKey: "plrid, gameid"},
// 			{PGSchema: "intake", PGTable: "pfdg", PGPKey: "plrid, gameid"},
// 		},
// 	)
// 	if err := metl.ExtractData(); err != nil {
// 		t.Errorf("failed extracting data\n%v", err)
// 	}

// 	metl.Dataset.(*etl.RespBoxscore).SetSharedVals(season, gameId)
// 	metl.Dataset.CleanTempFields()
// 	tableSets := metl.Dataset.SliceInsertRows()[0]

// 	// DO NOT DELETE
// 	for i, pgt := range metl.PGTargets {
// 		fmt.Printf("%v:\n%v\n++++++++++++\n\n", pgt.PGTable, tableSets[i])
// 		cols, err := pgresd.ColumnsInTable(db, pgt.PGTable)
// 		if err != nil {
// 			t.Errorf("failed to make InSt | %v\n", err)
// 		}
// 		rows := tableSets[i].([][]any)
// 		fmt.Println(rows)

// 		metl.InSt = pgresd.MakeInsert(pgt.PGSchema, pgt.PGTable, pgt.PGPKey,
// 			cols, rows)

// 		if err := metl.InSt.InsertFast(db, &metl.RowCount); err != nil {
// 			t.Errorf("failed to insert %d rows into %s\n%v",
// 				len(rows), pgt.PGTable, err)
// 		}

// 	}
// }

// }
