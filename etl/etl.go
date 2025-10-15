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

	"github.com/jdetok/golib/pgresd"
)

// get ETL struct
// ETLProcess interface implementation // db target info, endpoint info passed
func MakeETL(ds ETLProcess, sch, tbl, pkey, endpt string, params []Param) *ETL {
	gr := HTTPGet{
		Base:     BASE,
		Endpoint: endpt,
		Params:   params,
	}

	e := ETL{
		Dataset:  ds,
		Request:  gr,
		PgSchema: sch,
		PgTable:  tbl,
		PgPKey:   pkey,
		RowCount: 0,
	}
	return &e
}

// run entire ETL processes
func (e *ETL) RunFullETL(db *sql.DB) error {
	if err := e.ExtractData(); err != nil {
		return fmt.Errorf("** error extracting data from %s\n%w", e.Request.URL, err)
	}

	// call the appropriate struct method from the interface
	if err := e.Dataset.CleanTempFields(); err != nil {
		return fmt.Errorf(
			"** error cleaning fields for ETLProcess implementation struct\n%w", err)
	}

	cols, err := pgresd.ColumnsInTable(db, e.PgTable)
	if err != nil {
		return fmt.Errorf(
			"** error making slice of columns\n%w", err)
	}

	rows := e.Dataset.SliceInsertRows()

	e.InSt = pgresd.MakeInsert(e.PgSchema, e.PgTable, e.PgPKey, cols, rows)

	return e.InSt.InsertFast(db, &e.RowCount)
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
	// fmt.Println(e.Request.URL)
	// fmt.Println(string(js))

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

// try to rewrite InsertGames (load.go) as a generic using the ETL struct
// interface?? so i can do each structs own methods?
func (e *ETL) InsertIntoDB() error {

	return nil
}
