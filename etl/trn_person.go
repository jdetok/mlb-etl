package etl

import "fmt"

// ROSTER ENDPOINT
func (rp *RespRoster) CleanTempFields() error {
	return nil
}

func (rp *RespRoster) SliceInsertRows() [][]any {
	var rows [][]any
	for _, p := range rp.People {
		var vals = []any{
			p.Detail.ID, p.Detail.Name, p.Detail.Link, p.Jersey,
			p.Position.Code, p.Position.Name, p.Position.Type, p.Position.Abbr,
			p.Status.Code, p.Status.Desc, p.TeamID,
		}
		rows = append(rows, vals)
	}
	return rows
}

// PLAYERS ENDPOINT
func (pl *RespPlayers) CleanTempFields() error {
	// convert debut date and birthdate to dt
	for _, p := range pl.Players {
		ddErr := StrToDT(&p.TmpDebutDate, &p.DebutDate, BASIC_DATE_STR)
		bdErr := StrToDT(&p.TmpBirthDay, &p.BirthDay, BASIC_DATE_STR)
		if ddErr != nil || bdErr != nil {
			return fmt.Errorf(`
			failed converting debut date OR birthdate to date time
			** debute date error: %w
			** birthday error: %w
			`, ddErr, bdErr)
		}
	}
	return nil
}

// might want to add a playeridseason field
// how do i access season though?
func (pl *RespPlayers) SliceInsertRows() [][]any {
	var rows [][]any
	for _, p := range pl.Players {
		var vals = []any{
			p.ID, p.Name,
		}
		rows = append(rows, vals)
	}
	return rows
}
