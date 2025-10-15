package etl

import (
	"fmt"
)

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
	var ddErr, bdErr error
	// convert debut date and birthdate to dt
	for i := range pl.Players {
		if pl.Players[i].TmpDebutDate != "" {
			ddErr = StrToDT(&pl.Players[i].TmpDebutDate,
				&pl.Players[i].DebutDate, BASIC_DATE_STR)
		}

		if pl.Players[i].TmpBirthDay != "" {
			bdErr = StrToDT(&pl.Players[i].TmpBirthDay, &pl.Players[i].BirthDay,
				BASIC_DATE_STR)
		}

		SPrID, prErr := MakeSPrID(&pl.Players[i].ID, &pl.Season)
		if prErr != nil {
			// NEEDS TO BE LOGGED
			fmt.Printf(`
			** can't make player season primary key 
			* player id: %d | name: %s | season: %s
			** %v
			*** CONTINUNG | TODO ADD THIS TO LOGGER 
			`, pl.Players[i].ID, pl.Players[i].Name, pl.Season, prErr)
		}
		pl.Players[i].SPrID = *SPrID
		fmt.Println(pl.Players[i].SPrID)

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
			p.SPrID, p.ID, p.Name, p.Link, p.FName, p.LName, p.PrimNum,
			p.BirthDay, p.Age, p.BirthCity, p.BirthState, p.BirthCountry,
			p.Height, p.Weight, p.Active, p.CurrentTeam.ID, p.CurrentTeam.Link,
			p.PrimPos.Code, p.PrimPos.Name, p.PrimPos.Type, p.PrimPos.Abbr,
			p.UseName, p.UseLName, p.MName, p.BoxScoreName, p.Gender,
			p.IsPlayer, p.IsVerified, p.DraftYear, p.DebutDate,
			p.BatSide.Code, p.BatSide.Desc, p.PitchHand.Code, p.PitchHand.Desc,
			p.NameFL, p.NameSlug, p.FLName, p.LFName, p.LIName,
			p.FMLName, p.LFName, p.StrikeZoneTop, p.StrikeZoneBottom,
		}
		rows = append(rows, vals)
	}
	return rows
}
