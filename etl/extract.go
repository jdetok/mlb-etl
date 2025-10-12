package main

/* extract.go
this file should contain functions & types to extract data from the MLB API
i.e. func to build & send a get request & return the JSON as a slice of bytes
*/

import (
	"fmt"
	"io"
	"net/http"
)

// get request struct
type HTTPGet struct {
	Base     string
	Endpoint string
	Params   []Param
	URL      string // call BuildURL
}

// parameters for query string
type Param struct {
	Key string
	Val string
}

/*
build URL form HTTPGet struct, create new client, create new get request,
send request with client, read the body of the response
*/
func (gr *HTTPGet) SendGetRequest() ([]byte, error) {
	url := gr.BuildURL()
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

/*
concat HTTPGet items together to build query string
if only a single value is passed (e.g. []Param{{Key: "158"}}), only that value
will be appended to the url (preceded by a /) -> this looks like /v1/teams/158
*/
func (gr *HTTPGet) BuildURL() string {
	var url string = fmt.Sprintf("%s/%s", gr.Base, gr.Endpoint)

	// build query string if there are parameters
	lenP := len(gr.Params)
	if lenP > 0 {
		// HANDLE base/endpoint/value e.g. v1/teams/158
		if lenP == 1 && gr.Params[0].Val == "" {
			url += fmt.Sprintf("/%s", gr.Params[0].Key)
			return url
		}
		url += "?" // start query string
		for i, p := range gr.Params {
			url += fmt.Sprintf("%s=%s", p.Key, p.Val)
			if i < (len(gr.Params) - 1) { // concat & if not last param
				url += "&"
			}
		}
	}

	return url
}
