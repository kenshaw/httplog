package main

import (
	"log"
	"net/http"

	"github.com/kenshaw/httplog"
)

func main() {
	transport := httplog.NewPrefixedRoundTripLogger(nil, log.Printf)
	cl := &http.Client{
		Transport: transport,
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
