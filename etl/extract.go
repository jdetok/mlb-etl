package etl

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
	gr.BuildURL()
	client := &http.Client{}
	req, err := http.NewRequest("GET", gr.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new HTTP request\n%w", err)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP client failed to send HTTP request\n%w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read HTTP response body\n%w", err)
	}
	return body, nil
}

/*
concat HTTPGet items together to build query string
if only a single value is passed (e.g. []Param{{Key: "158"}}), only that value
will be appended to the url (preceded by a /) -> this looks like /v1/teams/158
*/
func (gr *HTTPGet) BuildURL() {
	var url string = fmt.Sprintf("%s/%s", gr.Base, gr.Endpoint)
	var params []Param // edit params to iterate through when building q string
	// build query string if there are parameters
	lenP := len(gr.Params)
	if lenP > 0 {
		// handle url edge cases (anything but /api/items?key=val&key1=val1)
		if gr.Params[0].Key != "" && gr.Params[0].Val == "" { // first value empty
			if lenP == 1 {
				url += fmt.Sprintf("/%s", gr.Params[0].Key)
				gr.URL = url
				return // return early if only one
			} // for teams/158/roster
			if gr.Params[1].Key != "" && gr.Params[1].Val == "" {
				url += fmt.Sprintf("/%s/%s", gr.Params[0].Key, gr.Params[1].Key)
				if lenP == 2 { // assign & return if only 2
					gr.URL = url
					return
				} else { // make params only third item and on
					params = gr.Params[2:]
				}
			}
		} else { // first key val good
			params = gr.Params
		}
		// build query string
		url += "?"
		for i, p := range params {
			url += fmt.Sprintf("%s=%s", p.Key, p.Val)
			if i < (len(params) - 1) { // concat & if not last param
				url += "&"
			}
		}
	}
	gr.URL = url
}
