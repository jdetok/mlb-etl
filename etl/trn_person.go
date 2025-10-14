package etl

func (rp *RespPeople) CleanTempFields() error {
	return nil
}

func (rp *RespPeople) SliceInsertRows() [][]any {
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
