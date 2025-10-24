package etl

import (
	"database/sql"
	"fmt"
	"strconv"
	"sync"

	"github.com/jdetok/mlb-etl/logd"
	"github.com/jdetok/mlb-etl/pgresd"
)

// LAUNCH GOROUTINES TO GET BOXSCORES FOR MANY GAMEIDS AT ONCE
// gets all game ids for the seasons attached to BatchETL, splits into several
// smaller chunks & starts the full ETL for each chunk
func (b *BatchETL) LoadManyBoxScoreETL(db *sql.DB, lg *logd.Logder,

// start, end, maxcon int,
) error {
	// query database for all game ids for specified seasons
	if err := b.GetGameIDsSeasons(db, lg); err != nil {
		lg.Log(fmt.Sprintf("failed to get game ids and seasons:\n%v\n", err), err, &b.RowCount)
		return err
	}

	// chunk all game ids
	b.ChunkGameIDs(20)

	// gsem := make(chan struct{}, b.MaxGoRtns)
	var gmu sync.Mutex
	var gwg sync.WaitGroup
	// start a goroutine that launches ETL for each chunk of game ids
	for i, chunk := range b.ChunkedGameIDs {
		// gsem <- struct{}{}
		gwg.Add(1)

		go func(lg *logd.Logder, ggmu *sync.Mutex, i int,
			chunk []map[uint64]string, rc *int64,
		) {
			defer gwg.Done()
			// defer func() { <-gsem }() // clear one spot in sem
			lg.Log(fmt.Sprintf("chunk %d/%d\n%v\n", i+1, len(b.ChunkedGameIDs),
				chunk), nil, rc)
			// mutex and waitgroup for safe concurrency
			var mu sync.Mutex
			var wg sync.WaitGroup

			sem := make(chan struct{}, b.MaxGoRtns)

			// get the map of each gameid/season and launch goroutine
			for _, gmap := range chunk {
				sem <- struct{}{}
				wg.Add(1)
				go b.LoadChunk(db, lg, gmap, &sem, &wg, &mu)
			}
			wg.Wait()
		}(b.Log, &gmu, i, chunk, b.TotalRC)
	}
	gwg.Wait()

	return nil
}

func (b *BatchETL) LoadChunk(db *sql.DB, lg *logd.Logder, gmap map[uint64]string,
	sem *chan struct{}, wg *sync.WaitGroup, mu *sync.Mutex,
) error {
	defer wg.Done()
	defer func() { <-*sem }() // clear one spot in sem

	var gameId string
	var season string
	for gId, szn := range gmap {
		gameId = strconv.FormatUint(gId, 10)
		season = szn
	}
	metl := MakeMultiTableETL(nil, &RespBoxscore{},
		"v1/game", []Param{
			{Key: gameId}, {Key: "boxscore"},
		}, []PGTarget{
			{PGSchema: "intake", PGTable: "tbtg", PGPKey: "teamid, gameid"},
			{PGSchema: "intake", PGTable: "tptg", PGPKey: "teamid, gameid"},
			{PGSchema: "intake", PGTable: "tfdg", PGPKey: "teamid, gameid"},
			{PGSchema: "intake", PGTable: "pbtg", PGPKey: "plrid, gameid"},
			{PGSchema: "intake", PGTable: "pptg", PGPKey: "plrid, gameid"},
			{PGSchema: "intake", PGTable: "pfdg", PGPKey: "plrid, gameid"},
		},
	)
	metl.Dataset.(*RespBoxscore).SetSharedVals(season, gameId)
	if err := metl.ExtractData(); err != nil {
		lg.Log(fmt.Sprintf("failed extracting data\n%v", err), err, &b.RowCount)
		return err
	}

	metl.Dataset.CleanTempFields()
	tableSets := metl.Dataset.SliceInsertRows()[0]
	for i, pgt := range metl.PGTargets {
		rows := tableSets[i].([][]any)
		if err := metl.BuildAndInsert(db, lg, &pgt, rows); err != nil {
			lg.Log(fmt.Sprintf("failed to do insert | %v\n", err), err, &b.RowCount)
			return err
		}
	}
	return nil
}

func (e *ETL) BuildAndInsert(db *sql.DB, lg *logd.Logder, pgt *PGTarget,
	rows [][]any,
) error {
	lg.Log(fmt.Sprintf("INSERT INTO %v:\nNUMBER OF ROWS: %v\n++++++++++++\n\n",
		pgt.PGTable, len(rows)), nil, &e.RowCount)
	cols, err := pgresd.ColumnsInTable(db, pgt.PGTable)
	if err != nil {
		lg.Log(fmt.Sprintf("failed to make InSt | %v\n", err), err, &e.RowCount)
		return err
	}

	e.InSt = pgresd.MakeInsert(pgt.PGSchema, pgt.PGTable, pgt.PGPKey,
		cols, rows)

	if err := e.InSt.InsertFast(db, &e.RowCount); err != nil {
		lg.Log(fmt.Sprintf("failed to insert %d rows into %s\n%v",
			len(rows), pgt.PGTable, err), err, &e.RowCount)
		return err
	}
	lg.Log(fmt.Sprintf("inserted %d rows into %s\n",
		len(rows), pgt.PGTable), nil, &e.RowCount)
	return nil
}

func (b *BatchETL) GetGameIDsSeasons(db *sql.DB, lg *logd.Logder) error {
	for szn := b.StartSzn; szn <= b.EndSzn; szn++ {
		if err := b.GetGameIdSzn(db, szn); err != nil {
			lg.Log(fmt.Sprintf("failed to get game ids:\n%v\n", err), err, &b.RowCount)
			return err
		}
		lg.Log(fmt.Sprintf("got %d game ids\n", len(b.GameIDs)), nil, &b.RowCount)
	}
	return nil
}

func (b *BatchETL) GetGameIdSzn(db *sql.DB, season int) error {
	q := `
select gameid, season
from intake.game_from_schedule
where season = $1
order by gdate
`
	rows, err := db.Query(q, season)
	if err != nil {
		return err
	}
	var idSznPairs []map[uint64]string
	for rows.Next() {
		var id uint64
		var szn string
		pair := make(map[uint64]string)
		if err := rows.Scan(&id, &szn); err != nil {
			return err
		}
		pair[id] = szn
		idSznPairs = append(idSznPairs, pair)
	}
	if len(b.GameIDs) == 0 {
		b.GameIDs = idSznPairs
	} else {
		b.GameIDs = append(b.GameIDs, idSznPairs...)
	}
	return nil
}

func (b *BatchETL) ChunkGameIDs(chunkSize int) {
	oglen := len(b.GameIDs)
	var chunks [][]map[uint64]string
	for i := 0; i < oglen; i += chunkSize {
		end := i + chunkSize
		chunks = append(chunks, b.GameIDs[i:end])
	}
	b.Log.Log(fmt.Sprintf(
		"total vals: %d | chunk size: %d | numChunks: %d |",
		oglen, chunkSize, len(chunks)), nil, b.TotalRC)
	b.ChunkedGameIDs = chunks
}
