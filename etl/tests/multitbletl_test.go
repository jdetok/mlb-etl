package etl

import (
	"fmt"
	"testing"

	"github.com/jdetok/mlb-etl/etl"
)

func TestMultiTblETL(t *testing.T) {
	// db, err := pgresd.ConnectTestDB("../../.env")
	// if err != nil {
	// 	t.Errorf("failed to connect to database | %v\n", err)
	// }

	metl := etl.MakeMultiTableETL(nil, &etl.RespBoxscore{},
		"v1/game", []etl.Param{{Key: "777933"}, {Key: "boxscore"}},
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
		t.Errorf("error now\n%v", err)
	}

	// fmt.Println(metl.Dataset.SliceInsertRows())
	fmt.Println(len(metl.Dataset.SliceInsertRows()))
	tableSets := metl.Dataset.SliceInsertRows()[0]
	for _, ts := range tableSets {
		fmt.Printf("ts: %v\n++++++++++++\n\n", ts)
	}
	// for _, p := range rows[0][1].([][]any) {
	// 	for _, val := range p {
	// 		fmt.Printf("%v | ", val)
	// 	}
	// 	fmt.Println()
	// }
	// // {metl.Dataset}
	// s, ok := metl.Dataset.(*etl.RespBoxscore)
	// if !ok {
	// 	t.Errorf("type assertion failure")
	// }

	// for _, t := range []etl.MLBTeamBoxScore{s.Teams.Away, s.Teams.Home} {
	// 	fmt.Println("TEAM +++++++++++++++++++++++++++++++++++++++++++++++")
	// 	fmt.Println(t.TeamStats.Batting.Avg)
	// 	// for _, p := range t.Players {
	// 	// 	fmt.Println(p)
	// 	// 	fmt.Println("================================================")
	// 	// }
	// 	fmt.Println("done with", t.TeamDtl.Abbr, "players")

	// }

}
