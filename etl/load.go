// insert extracted/transformed data into postgres db
package main

import (
	"database/sql"
	"fmt"

	"github.com/jdetok/golib/pgresd"
)

// insert games in RespSchedule into post gres using pgresd package
func (rs *RespSchedule) InsertGames(db *sql.DB) error {
	tbl := "game_from_schedule"
	schema := "intake"
	pkey := "id"
	cols, err := pgresd.ColumnsInTable(db, tbl)
	fmt.Println(cols)
	if err != nil {
		return err
	}
	var rows [][]any
	// vals in order of DB insert
	for i := range rs.Dates {
		for _, g := range rs.Dates[i].Games {

			// CHECK SINGLE CHARACTER VARIABLES FOR LENGTH
			var single_chars = []string{g.IfNecessary, g.DayType,
				g.Status.AbstractCode, g.Status.StateCode, g.Status.Code, g.Type}
			for i, v := range single_chars {
				if str := checkLen(v); str != "" {
					fmt.Printf("FIELD AT POSITION %d GREATER THAN 1 CHARACTER | %s\n", i+1, v)
				}
			}
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
	}
	ins := pgresd.MakeInsert(schema, tbl, pkey, cols, rows)
	var rc int64 = 0
	return ins.InsertFast(db, &rc)
}

func checkLen(str string) string {
	if len(str) > 1 {
		return str
	}
	return ""
}
