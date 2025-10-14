# fetch official MLB data & insert into Postgres

## dev work/running container
use `docker compose up --build` to start database container & 
`docker compose down --rmi all` to stop & remove the container

## ADDING A NEW ETL PROCESS (fetch from endpoint, structure in go, insert into postgres)
- an ETL struct should be defined with MLB endpoint, database table/schema/primary key,
and a Dataset struct that implements the ETLProcess interface
### ETL struct
```go 
// struct to hold endpoint, database, etc information for a full etl process
type ETL struct {
	Dataset  ETLProcess         // struct that implements interface below
	Request  HTTPGet            // http get request
	PgSchema string             // target db schema
	PgTable  string             // target db table
	PgPKey   string             // primary key of db table
	InSt     pgresd.InsertStmnt // insert statement struct from postgres library
	RowCount int64              // row count for postgres insert
}
```
### ETLProcess interface
```go
type ETLProcess interface {
	CleanTempFields() error   // run struct methods to convert temp fields
	SliceInsertRows() [][]any // run struct methods to get rows variable for Insert
}
```
- the `MakeETL()` function exists to create this an ETL struct
    - the function creates an HTTPGet struct with the endpoint & query string parameters
```go
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
```

- the `RunFullETL()` function exists to call the full ETL process on a defined ETL struct 
    - `ExtractData()` uses the HTTPGet struct to send a get request to the MLB API
    - the interface is taken advantage of with `e.Dataset.CleanTeamFields()` and `e.Dataset.SliceInsertRows()`
        - this enables using different logic for different structs while still always using this ETL struct methods
```go
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
```
## examples of calling MakeETL
- ### multiple parameters
    - will produce `/api/v1/schedule?sportId=1&season=1995&gameType=R`
```go
e := etl.MakeETL(&etl.RespSchedule{},
		"intake", "game_from_schedule", "id", "v1/schedule",
		[]etl.Param{
			{Key: "sportId", Val: "1"},
			{Key: "season", Val: "1995"},
			{Key: "gameType", Val: "R"},
		},
	)
```
- ## single paramater, no ?
    - will produce `/api/v1/teams/138`
```go 
	te := etl.MakeETL(&etl.RespTeams{},
		"intake", "team_detail", "id", "v1/teams", []etl.Param{{Key: "138"}})
``` 

- ## no parameters
    - will produce `/api/v1/teams`
```go 
	te := etl.MakeETL(&etl.RespTeams{},
		"intake", "team_detail", "id", "v1/teams", []etl.Param{{}})
``` 

- ## single parameter, not in query string, endpoint at end
    - will produce `/api/v1/teams/138/roster`
```go 
	te := etl.MakeETL(&etl.RespTeams{},
		"intake", "person", "id", "v1/teams", []etl.Param{{Key:"138"}, {Key:"roster"}})
``` 