package etl

import (
	"database/sql"

	"github.com/jdetok/mlb-etl/logd"
)

func (b *BatchETL) LoadManyBoxScoreETL(db *sql.DB, lg *logd.Logder) error {

	return nil
}

func (b *BatchETL) GetGameIDs(db *sql.DB, season string) error {
	q := `
select gameid
from intake.game_from_schedule
where season = $1
order by gdate
`
	rows, err := db.Query(q, season)
	if err != nil {
		return err
	}
	var ids []uint64
	for rows.Next() {
		var id uint64
		if err := rows.Scan(&id); err != nil {
			return err
		}
		ids = append(ids, id)
	}
	b.GameIDs = ids
	return nil
}
