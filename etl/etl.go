//	etl.go
//
// this file should contain functions that call different ETL functions
// i.e a function to send a get request (extract), unmarshal the response into a
// defined data structure (transform), & insert the data into a database (load)
package etl

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"sync"

	"github.com/jdetok/golib/pgresd"
	"github.com/jdetok/mlb-etl/logd"
)

// get ETL struct
// ETLProcess interface implementation // db target info, endpoint info passed
func MakeETL(ds ETLProcess, sch, tbl, pkey, endpt string, params []Param, lg *logd.Logder) *ETL {
	gr := HTTPGet{
		Base:     BASE,
		Endpoint: endpt,
		Params:   params,
		Log:      lg,
	}

	e := ETL{
		Dataset:  ds,
		Request:  gr,
		PgSchema: sch,
		PgTable:  tbl,
		PgPKey:   pkey,
		RowCount: 0,
		Log:      lg,
	}
	return &e
}

// run entire ETL processes
func (e *ETL) RunFullETL(db *sql.DB) error {
	if err := e.ExtractData(); err != nil {
		e.Log.Log(fmt.Sprintf("** error extracting data from %s", e.Request.URL), err, nil)
		return err
	}

	// call the appropriate struct method from the interface
	if err := e.Dataset.CleanTempFields(); err != nil {
		e.Log.Log("** error cleaning fields for ETLProcess implementation struct", err, nil)
		return err
	}

	cols, err := pgresd.ColumnsInTable(db, e.PgTable)
	if err != nil {
		e.Log.Log("** error making slice of columns", err, nil)
		return err
	}

	rows := e.Dataset.SliceInsertRows()

	e.InSt = pgresd.MakeInsert(e.PgSchema, e.PgTable, e.PgPKey, cols, rows)

	return e.InSt.InsertFast(db, &e.RowCount)
}

func IncrementRC(mu *sync.Mutex, total_rc, rc_to_add *int64) {
	mu.Lock()
	*total_rc += *rc_to_add
	mu.Unlock()
}

func CatchErr(mu *sync.Mutex, errs *[]error, err *error) {
	mu.Lock()
	*errs = append(*errs, *err)
	mu.Unlock()
}

// concurrent etl runs many seasons
func (b *BatchETL) RunManySznETL(db *sql.DB, lg *logd.Logder) error {

	// mutex and waitgroup for safe concurrency
	var mu sync.Mutex
	var wg sync.WaitGroup

	sem := make(chan struct{}, b.MaxGoRtns)

	// COLLECT ERRS
	var allErrs []error
	var errMu sync.Mutex

	// loop through each season, run ETLs
	for i := 0; i < b.EndSzn-b.StartSzn; i++ {
		sem <- struct{}{}
		wg.Add(1)
		go func(i int, rc *int64) {
			defer wg.Done()
			defer func() { <-sem }() // clear one spot in sem
			szn := strconv.Itoa(b.EndSzn - i)

			// SCHEDULE ENDPOINT TEST
			// schema | table | primary key | endpoint | endpoint parameters
			e := MakeETL(&RespSchedule{},
				"intake", "game_from_schedule", "id", "v1/schedule",
				[]Param{
					{Key: "sportId", Val: "1"},
					{Key: "season", Val: szn},
					{Key: "gameType", Val: "R"},
				}, lg,
			)
			if err := e.RunFullETL(db); err != nil {
				lg.Log("schedule endpoint failed", err, rc)
				CatchErr(&errMu, &allErrs, &err)
			}

			lg.Log(fmt.Sprintf("done with schedule etl for %s", szn), nil, &e.RowCount)
			IncrementRC(&mu, rc, &e.RowCount)

			var pl RespPlayers
			// sports/1/players?season=2025
			pl.Season = szn
			ple := MakeETL(&pl, "intake", "splayer", "sprid", "v1/sports",
				[]Param{
					{Key: "1"},
					{Key: "players"},
					{Key: "season", Val: pl.Season}}, lg)

			if err := ple.RunFullETL(db); err != nil {
				lg.Log("error with players endpoint", err, rc)
				CatchErr(&errMu, &allErrs, &err)
			}

			lg.Log(fmt.Sprintf("done with players etl for %s", pl.Season), nil, &ple.RowCount)

			IncrementRC(&mu, rc, &ple.RowCount)

			// TEAMS ETL
			te := MakeETL(&RespTeams{},
				"intake", "team_detail", "id", "v1/teams",
				[]Param{{Key: "season", Val: szn}}, lg)

			if err := te.RunFullETL(db); err != nil {
				lg.Log("error with teams endpoint", err, rc)
				CatchErr(&errMu, &allErrs, &err)
			}

			lg.Log(fmt.Sprintf("done with teams etl for %s", szn), nil, rc)

			IncrementRC(&mu, rc, &te.RowCount)
		}(i, &b.RowCount)
	}
	wg.Wait()

	numErrs := len(allErrs)
	if numErrs > 0 {
		err := fmt.Errorf("%d errors occured", numErrs)
		lg.Log("RunManySznETL error:", err, &b.RowCount)
		for i, e := range allErrs {
			lg.Log(fmt.Sprintf("error %d/%d", i, numErrs), e, &b.RowCount)
		}
		return err
	}
	lg.Log(fmt.Sprintf("FINAL ROW COUNT: %d",
		b.RowCount), nil, &b.RowCount)
	return nil
}

// create HTTP GET request from the passed endpoint and parameters, send the
// request with an HTTP client and get the JSON response, unmarshal the JSON into
// the struct passed as [T]
func (e *ETL) ExtractData() error {
	// send http request & get JSON response
	js, err := e.Request.SendGetRequest()
	if err != nil {
		return fmt.Errorf("** failed to send get request to %s\n%w", e.Request.URL, err)
	}

	// populate e.Dataset by unmarshalling json
	if err := e.ConvertJSONResp(js); err != nil {
		return fmt.Errorf("** failed to convert json to go struct\n%w", err)
	}
	return nil
}

// unmarshal json response to the implementing struct for ETL.Dataset
func (e *ETL) ConvertJSONResp(js []byte) error {
	if err := json.Unmarshal(js, e.Dataset); err != nil {
		return fmt.Errorf("failed to unmarshal json into %v\n%w",
			reflect.TypeOf(e.Dataset), err)
	}
	return nil
}

// create HTTP GET request from the passed endpoint and parameters, send the
// request with an HTTP client and get the JSON response, unmarshal the JSON into
// the struct passed as [T]
func GetAndMakeDS[T any](endpt string, params []Param) (*T, error) {
	// create get request
	gr := HTTPGet{
		Base:     BASE,
		Endpoint: endpt,
		Params:   params,
	}

	// get JSON response
	js, err := gr.SendGetRequest()
	if err != nil {
		return nil, err
	}

	// create & return the data structure passed at [T] from JSON response
	ds, err := MakeDS[T](js)
	if err != nil {
		return nil, err
	}
	return ds, nil
}
