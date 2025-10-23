package etl

import (
	"github.com/jdetok/mlb-etl/logd"
	"github.com/jdetok/mlb-etl/pgresd"
)

// MLB stats api base url
const BASE string = "https://statsapi.mlb.com/api"

// idea is to implement in all the structs returned JSON is unmarshalled into
// then can call that and do the individual logic for each data set to get into
// the intende d table. idea while trying to convert InsertGames to work with other structs
type ETLProcess interface {
	CleanTempFields() error   // run struct methods to convert temp fields
	SliceInsertRows() [][]any // run struct methods to get rows variable for Insert
}

type PGTarget struct {
	PGSchema string
	PGTable  string
	PGPKey   string
}

// struct to hold endpoint, database, etc information for a full etl process
type ETL struct {
	Dataset   ETLProcess         // struct that implements interface above
	Request   HTTPGet            // http get request
	PgSchema  string             // target db schema
	PgTable   string             // target db table
	PGTargets []PGTarget         // single endpoint inserting into many tables
	PgPKey    string             // primary key of db table
	InSt      pgresd.InsertStmnt // insert statement struct
	RowCount  int64              // row count for postgres insert
	Log       *logd.Logder
}

// parameters for RunManySznETL
type BatchETL struct {
	StartSzn  int
	EndSzn    int
	MaxGoRtns int
	RowCount  int64
	TotalRC   *int64 // global row count defined in main
}
