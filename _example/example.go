// _example/example.go
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/kenshaw/httplog"
)

func main() {
	cl := &http.Client{
		Transport: httplog.NewPrefixedRoundTripLogger(nil, os.Stdout),
		// without request or response body
		// Transport: httplog.NewPrefixedRoundTripLogger(nil, os.Stdout, httplog.WithReqResBody(false, false)),
	}
	req, err := http.NewRequest("GET", "https://google.com", nil)
	if err != nil {
		log.Fatal(err)
	}
	res, err := cl.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
}
