package main

import (
	"database/sql"
	"fmt"
	"reflect"
	"time"

	"github.com/jdetok/golib/pgresd"
)

/* loop through fields in a struct
t := reflect.TypeOf(teams.Teams[0])
for i := range t.NumField() {
	f := t.Field(i)
	fmt.Println("field:", f)
}
*/

func (rs *RespSchedule) InsertGames(db *sql.DB) error {
	// printStructFields(reflect.TypeOf(rs.Dates[0].Games[0]))
	tbl := "game_from_schedule"
	schema := "intake"
	sctbl := schema + "." + tbl
	pkey := "id"
	cols, err := pgresd.ColumnsInTable(db, tbl)
	fmt.Println(cols)
	if err != nil {
		return err
	}
	var rows [][]any
	for _, g := range rs.Dates[0].Games {
		var vals = []any{
			g.GID, g.GUID, g.Type, g.Season, g.DateTime, g.DateStr,
			// STATUS STRUCT
			g.Status.AbstractState, g.Status.StateCode, g.Status.State, g.Status.Code,
			g.Status.StartTBD, g.Status.AbstractCode,
			// HOME/AWAY FIELDS
			g.Teams.Home.Win, g.Teams.Away.Win, g.Teams.Home.Score, g.Teams.Away.Score,
			g.Teams.Home.Detail.ID, g.Teams.Away.Detail.ID,
			g.Teams.Home.Detail.Name, g.Teams.Away.Detail.Name,
			g.Teams.Home.Record.Wins, g.Teams.Away.Record.Wins,
			g.Teams.Home.Record.Losses, g.Teams.Away.Record.Losses,
			g.Teams.Home.Record.Pct, g.Teams.Away.Record.Pct,
			g.Teams.Home.SeriesNum, g.Teams.Away.SeriesNum,
			g.Teams.Home.SplitSquad, g.Teams.Away.SplitSquad,
			g.Teams.Home.Detail.Link, g.Teams.Away.Detail.Link,
			// VENUE FIELDS
			g.Venue.ID, g.Venue.Name, g.Venue.Link,
			// OTHER GAME FIELDS
			g.IsTie, g.DayType, g.DayNight, g.Description, g.SeasonDisplay,
			g.SeriesDesc, g.IfNecessary, g.IfNecessaryDesc,
		}
		rows = append(rows, vals)
	}
	ins := pgresd.MakeInsert(sctbl, pkey, cols, rows)
	// fmt.Println(ins.BuildStmnt(rows))
	// fmt.Println(rows)
	var rc int64 = 0
	// ins.FlattenVals()
	// ins.ChunkVals()
	// fmt.Println(ins.Chunks)
	fmt.Println(ins.BuildStmnt(rows))
	return ins.InsertFast(db, &rc)
	// return nil
}

func printStructFields(t reflect.Type) {
	for i := range t.NumField() {
		f := t.Field(i)
		fmt.Println(f.Name)
		if f.Type.Kind() == reflect.Struct && f.Type != reflect.TypeOf(time.Time{}) {
			printStructFields(f.Type)
		}
	}
}
