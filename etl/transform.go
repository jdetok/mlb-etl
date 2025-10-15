// primarily functions attached to response structs - clean data before inserting into db
package etl

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// common date string
const BASIC_DATE_STR string = "2006-01-02"

// GENERIC JSON TO GO STRUCT UNMARSHALER
// creates a variable of the desired type, attempts to unmarshal the passed js
// slice of bytes into that variable. returns pointer to the variable if successful
func MakeDS[T any](js []byte) (*T, error) {
	var v T
	if err := json.Unmarshal(js, &v); err != nil {
		fmt.Println(err)
		return &v, err
	}
	return &v, nil
}

// return string value if the string is greater than 1
// change to accept the number to check
func checkLen(str string) string {
	if len(str) > 1 {
		return str
	}
	return ""
}

// convert passed string to time.Time in target based on layout
func StrToDT(source *string, target *time.Time, layout string) error {
	if *source == "" {
		return fmt.Errorf("ERROR | cannot convert empty string")
	}
	dt, err := time.Parse(layout, *source)
	if err != nil {
		return err
	}
	*target = dt
	return nil
}

// concat season and person id to make player/season primary key field
func MakeSPrID(prId *uint64, season *string) (*uint64, error) {
	// first check that season is a valid integer
	if _, err := strconv.Atoi(*season); err != nil {
		return nil, fmt.Errorf(`SPrID generation failed | season <%s> must be an int | %w`,
			*season, err)
	}
	// player id as string
	prIdStr := strconv.FormatUint(*prId, 10)
	SPrIDStr := *season + prIdStr // concat ids as strings
	SPrIDUint, err := strconv.ParseUint(SPrIDStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf(`failed | generated <%s> could not be converted to uint64 | %w`,
			SPrIDStr, err)
	}
	// assign the uint64 to the original reference
	return &SPrIDUint, nil
}
