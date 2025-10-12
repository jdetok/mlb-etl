package main

/* etl.go
this file should contain functions that call different ETL functions
i.e a function to send a get request (extract), unmarshal the response into a
defined data structure (transform), & insert the data into a database (load)
*/

/*
create HTTP GET request from the passed endpoint and parameters, send the
request with an HTTP client and get the JSON response, unmarshal the JSON into
the struct passed as [T]
*/
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
