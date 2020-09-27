package httplog_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"

	"github.com/kenshaw/httplog"
)

func ExampleRoundTripLoggerWriter() {
	ts := httptest.NewServer(writeHTML(`<body>hello</body>`))
	defer ts.Close()

	w := NewMyWriter(os.Stdout)
	// do http request (w is a io.Writer)
	transport := httplog.NewPrefixedRoundTripLogger(nil, w)
	cl := &http.Client{
		Transport: transport,
	}
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		panic(err)
	}
	res, err := cl.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	// Output:
	// -> GET / HTTP/1.1
	// -> User-Agent: Go-http-client/1.1
	// -> Accept-Encoding: gzip
	// ->
	// ->
	// <- HTTP/1.1 200 OK
	// <- Content-Length: 18
	// <- Content-Type: text/html
	// <-
	// <- <body>hello</body>
}

func ExampleRoundTripLoggerPrintf() {
	ts := httptest.NewServer(writeHTML(`<body>hello</body>`))
	defer ts.Close()

	// do http request (printf has same signature as fmt.Printf)
	transport := httplog.NewPrefixedRoundTripLogger(nil, printf)
	cl := &http.Client{
		Transport: transport,
	}
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		panic(err)
	}
	res, err := cl.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	// Output:
	// -> GET / HTTP/1.1
	// -> User-Agent: Go-http-client/1.1
	// -> Accept-Encoding: gzip
	// ->
	// ->
	// <- HTTP/1.1 200 OK
	// <- Content-Length: 18
	// <- Content-Type: text/html
	// <-
	// <- <body>hello</body>
}

func ExampleRoundTripLoggerLogf() {
	ts := httptest.NewServer(writeHTML(`<body>hello</body>`))
	defer ts.Close()

	// do http request (logf has same signature as log.Printf)
	transport := httplog.NewPrefixedRoundTripLogger(nil, logf)
	cl := &http.Client{
		Transport: transport,
	}
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		panic(err)
	}
	res, err := cl.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	// Output:
	// -> GET / HTTP/1.1
	// -> User-Agent: Go-http-client/1.1
	// -> Accept-Encoding: gzip
	// ->
	// ->
	// <- HTTP/1.1 200 OK
	// <- Content-Length: 18
	// <- Content-Type: text/html
	// <-
	// <- <body>hello</body>
}

var cleanRE = regexp.MustCompile(`\n(->|<-) (Host|Date):.*`)
var spaceRE = regexp.MustCompile(`(?m)\s+$`)

func NewMyWriter(w io.Writer) *MyWriter {
	return &MyWriter{w: w}
}

type MyWriter struct {
	w io.Writer
}

func (w *MyWriter) Write(buf []byte) (int, error) {
	clean := cleanRE.ReplaceAll(buf, nil)
	clean = spaceRE.ReplaceAll(clean, nil)
	return os.Stdout.Write(append(clean, '\n'))
}

func printf(s string, v ...interface{}) (int, error) {
	clean := cleanRE.ReplaceAllString(fmt.Sprintf(s, v...), "")
	clean = spaceRE.ReplaceAllString(clean, "")
	return fmt.Println(clean)
}

func logf(s string, v ...interface{}) {
	clean := cleanRE.ReplaceAllString(fmt.Sprintf(s, v...), "")
	clean = spaceRE.ReplaceAllString(clean, "")
	fmt.Println(clean)
}

func writeHTML(content string) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "text/html")
		io.WriteString(res, strings.TrimSpace(content))
	})
}
