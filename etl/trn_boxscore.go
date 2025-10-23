package etl

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Convert interface {
	ConvertStructFields()
}

// pass season & set
func (rb *RespBoxscore) SetSharedVals(season, gameId string) {
	rb.Season = season
	gId, _ := strconv.ParseUint(gameId, 10, 64)
	rb.GameID = gId
}

// call cleanStructRecursive, uses reflect to find fields with
// `convert:"true"` tag - finds strings that need to be numbers
func (rb *RespBoxscore) CleanTempFields() error {
	count := cleanStructRecursive(reflect.ValueOf(rb))
	fmt.Println("total vals cleaned:", count)
	return nil
}

// recursively clean tagged string fields
func cleanStructRecursive(v reflect.Value) int {
	if !v.IsValid() {
		return 0
	}

	// If pointer, dereference it
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return 0
		}
		v = v.Elem()
	}

	count := 0

	switch v.Kind() {
	case reflect.Struct:
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			fld := v.Field(i)
			sf := t.Field(i)

			// Recurse for nested structs or maps
			if fld.Kind() == reflect.Struct {
				if fld.CanAddr() {
					count += cleanStructRecursive(fld.Addr())
				} else {
					count += cleanStructRecursive(fld)
				}
				continue
			}
			if fld.Kind() == reflect.Map {
				for _, key := range fld.MapKeys() {
					count += cleanStructRecursive(fld.MapIndex(key))
				}
				continue
			}

			// Clean tagged string fields
			tag := sf.Tag.Get("convert")
			if tag == "true" && fld.Kind() == reflect.String && fld.CanSet() {
				str := fld.String()
				if strings.Contains(str, ".-") || strings.Contains(str, "") {
					fmt.Printf("updating %v | old str: %s | new str %s |\n",
						fld.Kind(), str, "0")

					fld.SetString("0")

					count++
				}
			}
		}

	case reflect.Map:
		for _, key := range v.MapKeys() {
			count += cleanStructRecursive(v.MapIndex(key))
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			count += cleanStructRecursive(v.Index(i))
		}
	}

	return count
}

// slices the data into six separate [][]any, caller shares an idx with the
// PGTargets slice (each table gets its own [][]any to insert)
// returns a [][]any for each target table
func (rb *RespBoxscore) SliceInsertRows() [][]any {
	var rows [][]any
	var tbtgRows, tptgRows, tfdgRows, pbtgRows, pptgRows, pfdgRows [][]any

	for _, tbs := range []MLBTeamBoxScore{rb.Teams.Away, rb.Teams.Home} {
		// team, game. season
		var tmMetaVals = []any{tbs.TeamDtl.ID, rb.GameID, rb.Season}
		ts := tbs.TeamStats
		// ts.DashesToBlank()
		tb := ts.Batting
		var tbtgVals []any
		tbtgVals = append(tbtgVals, tmMetaVals...)
		tbtgVals = append(tbtgVals, []any{
			tb.FlyOuts, tb.GndOuts, tb.Airouts, tb.Doubles, tb.Triples,
			tb.HomeRuns, tb.StrikeOuts, tb.BaseOnBalls, tb.IntnWalks, tb.Hits,
			tb.HitByPitch, tb.Avg, tb.AtBats, tb.OBP, tb.SLG, tb.OPS,
			tb.CaughtStealing, tb.StolenBases, tb.StolenBasesPct, tb.GndIntoDP,
			tb.GndIntoTP, tb.PlateApps, tb.TotalBases, tb.RBI, tb.LeftOnBase,
			tb.SacBunts, tb.SacFlies, tb.CatchersIntr, tb.Pickoffs,
			tb.AtBatPerHR, tb.PopOuts, tb.LineOuts,
		}...)
		tbtgRows = append(tbtgRows, tbtgVals)

		tp := tbs.TeamStats.Pitching
		tpb := tp.BattingFields
		var tptgVals []any
		tptgVals = append(tptgVals, tmMetaVals...)
		tptgVals = append(tptgVals, []any{
			tpb.FlyOuts, tpb.GndOuts, tpb.Airouts, tpb.Doubles, tpb.Triples,
			tpb.HomeRuns, tpb.StrikeOuts, tpb.BaseOnBalls, tpb.IntnWalks,
			tpb.Hits, tpb.HitByPitch, tpb.Avg, tpb.AtBats, tpb.OBP,
			tpb.CaughtStealing, tpb.StolenBases,
			tpb.StolenBasesPct, tp.CaughtStealingPct, tp.NumPitches, tp.ERA,
			tp.InningsPitched, tp.SaveOpps, tp.EarnedRuns, tp.Whip,
			tp.BattersFaced, tp.CompleteGames, tp.Shutouts, tp.PitchesThrown,
			tp.Balls, tp.Strikes, tp.StrikePct, tp.HitBatsmen, tp.Balks,
			tp.WildPitches, tpb.Pickoffs, tp.GndToAir, tpb.RBI, tp.PitchPerInning,
			tp.RunsPer9, tp.HRPer9, tp.InhrRunners, tp.InhrRunnersScored,
			tpb.CatchersIntr, tpb.SacBunts, tpb.SacFlies, tp.PassedBall,
			tpb.PopOuts, tpb.LineOuts,
		}...)
		tptgRows = append(tptgRows, tptgVals)

		tf := tbs.TeamStats.Fielding
		var tfdgVals []any
		tfdgVals = append(tfdgVals, tmMetaVals...)
		tfdgVals = append(tfdgVals, []any{
			tf.CaughtStealing, tf.StolenBases, tf.StolenBasePct,
			tf.CaughtStealingPct, tf.Assists, tf.PutOuts, tf.Errors, tf.Chances,
			tf.PassedBall, tf.Pickoffs,
		}...)
		tfdgRows = append(tfdgRows, tfdgVals)

		for _, p := range tbs.Players {
			// player, team, game, season
			var prMetaVals []any
			ps := p.Stats
			// ps.DashesToBlank()
			prMetaVals = append(prMetaVals, p.Person.ID)
			prMetaVals = append(prMetaVals, tmMetaVals...)

			// player batting
			pb := ps.Batting
			var pbtgVals []any
			pbtgVals = append(pbtgVals, prMetaVals...)
			pbtgVals = append(pbtgVals, []any{
				pb.Summary, pb.GP, pb.FlyOuts, pb.GndOuts, pb.Airouts,
				pb.Doubles, pb.Triples, pb.HomeRuns, pb.StrikeOuts, pb.BaseOnBalls,
				pb.IntnWalks, pb.Hits, pb.HitByPitch, pb.AtBats, pb.CaughtStealing,
				pb.StolenBases, pb.StolenBasesPct, pb.GndIntoDP, pb.GndIntoTP,
				pb.PlateApps, pb.TotalBases, pb.RBI, pb.LeftOnBase, pb.SacBunts,
				pb.SacFlies, pb.CatchersIntr, pb.Pickoffs, pb.AtBatPerHR,
				pb.PopOuts, pb.LineOuts,
			}...)
			pbtgRows = append(pbtgRows, pbtgVals)

			// player pitching
			pp := ps.Pitching
			ppb := pp.BattingFields
			var pptgVals []any
			pptgVals = append(pptgVals, prMetaVals...)
			pptgVals = append(pptgVals, []any{
				ppb.Summary, ppb.GP, ppb.FlyOuts, ppb.GndOuts, ppb.Airouts,
				ppb.Doubles, ppb.Triples, ppb.HomeRuns, ppb.StrikeOuts, ppb.BaseOnBalls,
				ppb.IntnWalks, ppb.Hits, ppb.HitByPitch, ppb.Avg, ppb.AtBats,
				ppb.OBP, ppb.CaughtStealing,
				ppb.StolenBases, ppb.StolenBasesPct, pp.CaughtStealingPct, pp.NumPitches, pp.ERA,
				pp.InningsPitched, pp.SaveOpps, pp.Holds, pp.BlownSaves,
				pp.EarnedRuns, pp.Whip, pp.BattersFaced, pp.Outs,
				pp.CompleteGames, pp.Shutouts, pp.PitchesThrown, pp.Balls,
				pp.Strikes, pp.StrikePct, pp.HitBatsmen, pp.Balks,
				pp.WildPitches, ppb.Pickoffs, pp.GndToAir, ppb.RBI,
				pp.WinPct, pp.PitchPerInning, pp.GamesFinished, pp.SOWalkRatio,
				pp.SOPer9, pp.WalksPer9, pp.HitsPer9, pp.RunsPer9, pp.HRPer9,
				pp.InhrRunners, pp.InhrRunnersScored, ppb.CatchersIntr,
				ppb.SacBunts, ppb.SacFlies, pp.PassedBall, ppb.PopOuts,
				ppb.LineOuts,
			}...)
			// fmt.Println("player insert fields:", len(pptgVals))
			pptgRows = append(pptgRows, pptgVals)

			// player fielding
			pf := ps.Fielding
			var pfdgVals []any
			pfdgVals = append(pfdgVals, prMetaVals...)
			pfdgVals = append(pfdgVals, []any{
				pf.CaughtStealing, pf.StolenBases, pf.StolenBasePct,
				pf.CaughtStealingPct, pf.Assists, pf.PutOuts, pf.Errors,
				pf.Chances, pf.Fielding, pf.PassedBall, pf.Pickoffs,
			}...)
			pfdgRows = append(pfdgRows, pfdgVals)
		}
	}
	// append each slice of table rows ([]any type) to rows
	rows = append(rows, []any{
		tbtgRows, tptgRows, tfdgRows, pbtgRows, pptgRows, pfdgRows,
	})
	return rows
}
