package etl

import (
	"database/sql"
	"fmt"

	"github.com/jdetok/mlb-etl/logd"
)

func (b *BatchETL) LoadManyBoxScoreETL(db *sql.DB, lg *logd.Logder) error {

	return nil
}

func (b *BatchETL) GetGameIDs(db *sql.DB, season int) error {
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
	if len(b.GameIDs) == 0 {
		b.GameIDs = ids
	} else {
		b.GameIDs = append(b.GameIDs, ids...)
	}
	return nil
}

func (b *BatchETL) ChunkGameIDs() {
	chunkSize := 10
	oglen := len(b.GameIDs)
	var chunks [][]uint64
	for i := 0; i < oglen; i += chunkSize {
		end := i + chunkSize
		chunks = append(chunks, b.GameIDs[i:end])
	}
	fmt.Println("total vals: ", oglen, "chunk size:", chunkSize, " | numChunks: ", len(chunks))
	b.ChunkedGameIDs = chunks
}
