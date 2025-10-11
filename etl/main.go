package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

// MLB stats api base url
const BASE string = "https://statsapi.mlb.com/api/v1"

// super quick error handling for testing, replace later
func ErrHndl(err error) {
	fmt.Println("an error occured: killing program")
	log.Fatal(err)
}

// parameters for query string
type Param struct {
	Key string
	Val string
}

// get request struct
type HTTPGet struct {
	Base     string
	Endpoint string
	Params   []Param
	URL      string // call BuildURL
}

/*
build URL form HTTPGet struct, create new client, create new get request,
send request with client, read the body of the response
*/
func (gr *HTTPGet) SendGetRequest() {
	url := gr.BuildURL()
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		ErrHndl(err)
	}

	res, err := client.Do(req)
	if err != nil {
		ErrHndl(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		ErrHndl(err)
	}
	fmt.Println(string(body))
}

// concat HTTPGet items together to build query string
func (gr *HTTPGet) BuildURL() string {
	var url string = fmt.Sprintf("%s/%s", gr.Base, gr.Endpoint)

	// build query string if there are parameters
	if len(gr.Params) > 0 {
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

func main() {
	gr := HTTPGet{
		Base:     BASE,
		Endpoint: "schedule",
		Params: []Param{
			{Key: "sportId", Val: "1"},
			{Key: "season", Val: "2025"},
			{Key: "gameType", Val: "R"},
		},
	}
	gr.SendGetRequest()
}
