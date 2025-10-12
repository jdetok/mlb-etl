package main

import (
	"encoding/json"
	"fmt"
	"time"
)

/* TODO
- convert pct formatted by JSON as ".555" to a float
*/

/*
GENERIC JSON TO GO STRUCT UNMARSHALER
creates a variable of the desired type, attempts to unmarshal the passed js
slice of bytes into that variable. returns pointer to the variable if successful
*/
func MakeDS[T any](js []byte) (*T, error) {
	var v T
	if err := json.Unmarshal(js, &v); err != nil {
		fmt.Println(err)
		return &v, err
	}
	return &v, nil
}

// convert all TmpDateTime to DateTime in schedule response
func (rs *RespSchedule) GameDatesToDT() error {
	for _, d := range rs.Dates {
		for _, g := range d.Games {
			if err := g.toDateTime(); err != nil {
				return err
			}
			fmt.Println(g.TmpDateTime, "|||", g.DateTime)
		}
	}
	return nil
}

/*
convert string TmpDateTime to time.Time DateTime
RFC3339 format: 2025-10-04T18:08:00Z
*/
func (g *MLBGame) toDateTime() error {
	dt, err := time.Parse(time.RFC3339, g.TmpDateTime)
	if err != nil {
		return err
	}
	g.DateTime = dt
	return nil
}
