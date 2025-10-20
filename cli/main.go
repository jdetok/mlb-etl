// entry point
package main

import (
	"fmt"
	"log"

	"github.com/jdetok/golib/pgresd"
	"github.com/jdetok/mlb-etl/etl"
	"github.com/jdetok/mlb-etl/logd"
)

// super quick error handling for testing, replace later
func ErrHndl(err error) {
	fmt.Println("an error occured: killing program")
	log.Fatal(err)
}

func main() {
	lg := logd.Logder{Prj: "mlb-etl"}
	lg.Log("starting", nil, nil)

	// database connection
	db, dbErr := pgresd.ConnectDB()
	if dbErr != nil {
		lg.Log("couldn't connect to db", dbErr, nil)
	}
	db.SetMaxOpenConns(200)
	db.SetMaxIdleConns(200)
	lg.DB = db

	start := 2000
	end := 2025

	// max number of goroutines
	maxcon := 10

	// total row count
	total_rows := int64(0)

	betl := etl.BatchETL{
		StartSzn:  start,
		EndSzn:    end,
		MaxGoRtns: maxcon,
		RowCount:  0,
		TotalRC:   &total_rows,
	}

	err := betl.RunManySznETL(db, &lg)
	if err != nil {
		lg.Log("error(s) occured running many seasons", err, betl.TotalRC)
	}
}
