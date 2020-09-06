// Package httplog provides a standard http.RoundTripper transport that can be
// used with standard HTTP clients to log the raw (outgoing) HTTP request and
// response.
package httplog

import (
	"bytes"
	"net/http"
	"net/http/httputil"
)

// DefaultTransport is the default transport used by the HTTP logger.
var DefaultTransport = http.DefaultTransport

// RoundTripLogger provides a standard http.RoundTripper transport that can be
// used with standard HTTP clients to log the raw (outgoing) HTTP request and
// response.
type RoundTripLogger struct {
	transport http.RoundTripper
	reqf      func([]byte)
	resf      func([]byte)
}

// NewRoundTripLogger creates a new HTTP transport that logs the raw (outgoing)
// HTTP request and response.
func NewRoundTripLogger(transport http.RoundTripper, reqf, resf func([]byte)) *RoundTripLogger {
	return &RoundTripLogger{
		transport: transport,
		reqf:      reqf,
		resf:      resf,
	}
}

// NewPrefixedRoundTripLogger creates a new HTTP transport that logs the raw
// (outgoing) HTTP request and response to the provided standard logger funcs.
// Prefixes requests and responses with "-> " and "<-", respectively. Adds an
// additional blank line ("\n\n") to the output of requests and responses.
func NewPrefixedRoundTripLogger(transport http.RoundTripper, logf func(string, ...interface{})) *RoundTripLogger {
	nl := []byte("\n")
	f := func(prefix []byte, buf []byte) {
		buf = append(prefix, bytes.ReplaceAll(buf, nl, append(nl, prefix...))...)
		logf("%s\n\n", string(buf))
	}
	return NewRoundTripLogger(
		transport,
		func(buf []byte) { f([]byte("-> "), buf) },
		func(buf []byte) { f([]byte("<- "), buf) },
	)
}

// RoundTrip satisfies the http.RoundTripper interface.
func (l *RoundTripLogger) RoundTrip(req *http.Request) (*http.Response, error) {
	trans := l.transport
	if trans == nil {
		trans = DefaultTransport
	}
	reqBody, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return nil, err
	}
	res, err := trans.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	resBody, err := httputil.DumpResponse(res, true)
	if err != nil {
		return nil, err
	}
	l.reqf(reqBody)
	l.resf(resBody)
	return res, err
}
