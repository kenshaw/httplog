// _example/main.go
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
