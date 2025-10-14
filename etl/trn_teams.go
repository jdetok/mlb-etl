package etl

func (rt *RespTeams) CleanTempFields() error {
	return nil
}

func (rt *RespTeams) SliceInsertRows() [][]any {
	var rows [][]any
	for _, t := range rt.Teams {
		var vals = []any{
			t.ID, t.Name, t.Link, t.Season, t.Abbr, t.Code, t.TeamName, t.Location,
			t.League.ID, t.League.Name, t.League.Link, t.Division.ID, t.Division.Name,
			t.Division.Link, t.Sport.ID, t.Sport.Name, t.Sport.Link, t.ShortName,
			t.FranchiseName, t.ClubName, t.FirstYear, t.FileCode, t.AllStarSt,
			t.Active, t.Venue.ID, t.Venue.Name, t.Venue.Link, t.SpringVenue.ID,
			t.SpringVenue.Name, t.SpringVenue.Link, t.SpringLeague.ID, t.SpringLeague.Name,
			t.SpringLeague.Link, t.SpringLeague.Abbr,
		}
		rows = append(rows, vals)
	}
	return rows
}
