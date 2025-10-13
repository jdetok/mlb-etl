package etl

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// loop through dates then games to get individual rows
func (rs *RespSchedule) SliceInsertRows() [][]any {
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
	return rows
}

// convert all TmpDateTime to DateTime in schedule response
func (rs *RespSchedule) CleanTempFields() error {
	for _, d := range rs.Dates {
		for _, g := range d.Games {
			if err := g.toDateTime(); err != nil {
				return err
			}
			if err := g.toFloats(); err != nil {
				return err
			}
			// fmt.Println(g.TmpDateTime, "|||", g.DateTime)
		}
	}
	return nil
}

// convert string TmpDateTime to time.Time DateTime
// RFC3339 format: 2025-10-04T18:08:00Z
func (g *MLBGame) toDateTime() error {
	dt, err := time.Parse(time.RFC3339, g.TmpDateTime)
	if err != nil {
		return err
	}
	g.DateTime = dt
	return nil
}

// - convert pct formatted by JSON as ".555" to a float
func (r *MLBSeriesRecord) pctToFloat() error {
	pct, err := strconv.ParseFloat(strings.TrimPrefix(r.PctStr, "."), 64)
	if err != nil {
		return err
	}
	r.Pct = pct / 1000 // convert "505" → 0.505
	return nil
}

func (g *MLBGame) toFloats() error {
	if err := g.Teams.Away.Record.pctToFloat(); err != nil {
		return err
	}
	if err := g.Teams.Home.Record.pctToFloat(); err != nil {
		return err
	}
	return nil
}
