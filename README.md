# httplog

Package `httplog` provides a standard http.RoundTripper transport that can be
used with standard HTTP clients to log the raw (outgoing) HTTP request and
response.

## Example

```go
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
		// Transport: httplog.NewPrefixedRoundTripLogger(nil, os.Stdout, httplog.WithResReqBody(false, false)),
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
```
