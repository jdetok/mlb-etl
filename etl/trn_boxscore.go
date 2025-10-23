package etl

func (rb *RespBoxscore) CleanTempFields() error {
	return nil
}

func (rb *RespBoxscore) SliceInsertRows() [][]any {
	var rows [][]any
	var tbtgRows, tptgRows, tfdgRows, pbtgRows, pptgRows, pfdgRows [][]any

	for _, tbs := range []MLBTeamBoxScore{rb.Teams.Away, rb.Teams.Home} {
		// team, game. season
		var tmMetaVals = []any{tbs.TeamDtl.ID, rb.GameID, rb.Season}

		tb := tbs.TeamStats.Batting
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
			tpb.CaughtStealing, tpb.CaughtStealing, tpb.StolenBases,
			tpb.StolenBasesPct, tp.CaughtStealingPct, tp.NumPitches, tp.ERA,
			tp.InningsPitched, tp.SaveOpps, tp.EarnedRuns, tp.Whip,
			tp.BattersFaced, tp.CompleteGames, tp.Shutouts, tp.PitchesThrown,
			tp.Balls, tp.Strikes, tp.StrikePct, tp.HitBatsmen, tp.Balks,
			tp.WildPitches, tpb.Pickoffs, tp.GndToAir, tpb.RBI, tp.PitchPerInning,
			tp.RunsPer9, tp.HRPer9, tp.InhrRunners, tp.InhrRunnersScored,
			tpb.CatchersIntr, tpb.SacBunts, tpb.SacFlies, tpb.PopOuts, tpb.LineOuts,
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
			prMetaVals = append(prMetaVals, p.Person.ID)
			prMetaVals = append(prMetaVals, tmMetaVals...)

			// player batting
			pb := p.Stats.Batting
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
			pp := p.Stats.Pitching
			ppb := pp.BattingFields
			var pptgVals []any
			pptgVals = append(pptgVals, prMetaVals...)
			pptgVals = append(pptgVals, []any{
				ppb.Summary, ppb.GP, ppb.FlyOuts, ppb.GndOuts, ppb.Airouts,
				ppb.Doubles, ppb.Triples, ppb.HomeRuns, ppb.StrikeOuts, ppb.BaseOnBalls,
				ppb.IntnWalks, ppb.Hits, ppb.HitByPitch, ppb.AtBats, ppb.CaughtStealing,
				ppb.StolenBases, ppb.StolenBasesPct, pp.NumPitches, pp.ERA,
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
			pptgRows = append(pptgRows, pptgVals)

			// player fielding
			pf := p.Stats.Fielding
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
