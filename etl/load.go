// insert extracted/transformed data into postgres db
package etl

import (
	"database/sql"
	"fmt"

	"github.com/jdetok/golib/pgresd"
)

// insert games in RespSchedule into post gres using pgresd package
func (rs *RespSchedule) InsertGames(db *sql.DB) error {
	tbl := "game_from_schedule"
	schema := "intake"
	pkey := "id"
	cols, err := pgresd.ColumnsInTable(db, tbl)
	fmt.Println(cols)
	if err != nil {
		return err
	}
	var rows [][]any = rs.SliceInsertRows()
	ins := pgresd.MakeInsert(schema, tbl, pkey, cols, rows)
	var rc int64 = 0
	return ins.InsertFast(db, &rc)
}
