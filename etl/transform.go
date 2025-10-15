// primarily functions attached to response structs - clean data before inserting into db
package etl

import (
	"encoding/json"
	"fmt"
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
	dt, err := time.Parse(layout, *source)
	if err != nil {
		return err
	}
	*target = dt
	return nil
}
