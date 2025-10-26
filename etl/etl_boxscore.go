package etl

import (
	"database/sql"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/jdetok/mlb-etl/logd"
)

type ETLChans struct {
	done    chan string
	gnerrs  chan string
	exerrs  chan string
	dberrs  chan string
	success chan string
	sem     chan struct{}
}

func (b *BatchETL) BoxScoreETL(db *sql.DB, lg *logd.Logder,
	start, end, grLim, chunkSize int, delay time.Duration) {
	// set start and end years
	b.StartSzn = start
	b.EndSzn = end

	// global wait group
	var mwg sync.WaitGroup

	// channels to pass through goroutines
	var ch ETLChans
	ch.done = make(chan string)
	ch.gnerrs = make(chan string, 200)
	ch.exerrs = make(chan string, 200)
	ch.dberrs = make(chan string, 200)
	ch.success = make(chan string, 200)
	ch.sem = make(chan struct{}, grLim)

	// target postgres tables (moved to high level caller instead of LoadChunk)
	tables := []PGTarget{
		{PGSchema: "intake", PGTable: "tbtg", PGPKey: "teamid, gameid"},
		{PGSchema: "intake", PGTable: "tptg", PGPKey: "teamid, gameid"},
		{PGSchema: "intake", PGTable: "tfdg", PGPKey: "teamid, gameid"},
		{PGSchema: "intake", PGTable: "pbtg", PGPKey: "plrid, gameid"},
		{PGSchema: "intake", PGTable: "pptg", PGPKey: "plrid, gameid"},
		{PGSchema: "intake", PGTable: "pfdg", PGPKey: "plrid, gameid"},
	}

	// call to a func here that gets the seasons/game ids
	if err := b.GetGameIDsSeasons(db, lg); err != nil {
		ch.gnerrs <- fmt.Sprintf(
			"failed getting gameid:season maps from db\n**%v", err)
	}
	ch.success <- "Successfully retrieved game ids from db"

	// split into chunks
	b.ChunkGameIDs(chunkSize)

	// TODO: rewrite iteration through chunks
	// LOOP THROUGH CHUNKS
	for ic, chunk := range b.ChunkedGameIDs {
		// chunk is a []map[uint64]string, loop through each & launch goroutine
		// waitgroup & mutex should be created here for each individual chunkss
		var cwg sync.WaitGroup

		var gameId string
		var season string
		// iterate through each game map in the chunk, launch goroutine for each
		for _, g := range chunk {
			ch.sem <- struct{}{} // write to sem chan
			cwg.Add(1)

			go func() {
				defer cwg.Done()
				defer func() { <-ch.sem }()

				// game wait group (for pg goroutines)
				var gwg sync.WaitGroup

				for gId, szn := range g {
					gameId = strconv.FormatUint(gId, 10)
					season = szn
				}
				metl := MakeMultiTableETL(nil, &RespBoxscore{},
					"v1/game", []Param{{Key: gameId}, {Key: "boxscore"}}, tables)

				metl.Dataset.(*RespBoxscore).SetSharedVals(season, gameId)
				if err := metl.ExtractData(); err != nil {
					ch.exerrs <- fmt.Sprintf("failed to extract data from %s | %v\n", metl.Request.URL, err)
				}

				metl.Dataset.CleanTempFields()
				tableSets := metl.Dataset.SliceInsertRows()[0]
				for i, pgt := range metl.PGTargets {
					gwg.Add(1)
					go func(i int, pgt PGTarget) {
						defer gwg.Done()
						rows := tableSets[i].([][]any)
						if err := metl.BuildAndInsert(db, lg, &pgt, rows); err != nil {
							ch.dberrs <- fmt.Sprintf("failed to do insert into %s | %v\n", pgt.PGTable, err)
						}
					}(i, pgt)
				}
				gwg.Wait()
			}()
			cwg.Wait()
		}
		mwg.Wait()
		fmt.Printf("chunk %d/%d complete | pausing for %d seconds\n",
			ic, len(b.ChunkedGameIDs), delay)
		time.Sleep(delay * time.Second)
	}

	// LOGGER GOROUTINE (AT BOTTOM JUST BEFORE CLOSING CHANS)
	mwg.Add(1)
	go func() {
		defer mwg.Done()
		for {
			select {
			case msg := <-ch.gnerrs:
				lg.CCLog("GENERAL/SETUP ERROR:"+msg, true, b.RowCount, nil)
			case msg := <-ch.dberrs:
				lg.CCLog("DB ERROR:"+msg, true, b.RowCount, nil)
			case msg := <-ch.exerrs:
				lg.CCLog("EXTRACT ERROR:"+msg, true, b.RowCount, nil)
			case msg := <-ch.success:
				lg.CCLog("SUCCESS:"+msg, true, b.RowCount, nil)
			case <-ch.done:
				lg.CCLog("COMPLETE", true, b.RowCount, nil)
				return
			}
		}
	}()
	// wait for global wait group and close all channel
	mwg.Wait()
	close(ch.success)
	close(ch.gnerrs)
	close(ch.dberrs)
	close(ch.exerrs)

	// done signal stops logger
	<-ch.done
}
