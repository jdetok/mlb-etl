package main

import (
	"encoding/json"
	"fmt"
)

/*
GENERIC JSON TO GO STRUCT UNMARSHALER
creates a variable of the desired type, attempts to unmarshal the passed js
slice of bytes into that variable. returns pointer to the variable if successful
*/
func MakeDS[T any](js []byte) (*T, error) {
	var v T
	if err := json.Unmarshal(js, &v); err != nil {
		fmt.Println(err)
		return &v, err
	}
	return &v, nil
}
