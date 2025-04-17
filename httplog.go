// Package httplog provides a standard [http.RoundTripper] transport that can
// be used with standard HTTP clients to log the raw (outgoing) HTTP request
// and response.
package httplog

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
)

// DefaultTransport is the default transport used by the HTTP logger.
var DefaultTransport = http.DefaultTransport

// RoundTripLogger provides a standard [http.RoundTripper] transport that can
// be used with standard HTTP clients to log the raw (outgoing) HTTP request
// and response.
type RoundTripLogger struct {
	transport http.RoundTripper
	reqf      func([]byte)
	resf      func([]byte)
	noReqBody bool
	noResBody bool
}

// NewRoundTripLogger creates a new HTTP transport that logs the raw (outgoing)
// HTTP request and response.
func NewRoundTripLogger(transport http.RoundTripper, reqf, resf func([]byte), opts ...Option) *RoundTripLogger {
	l := &RoundTripLogger{
		transport: transport,
		reqf:      reqf,
		resf:      resf,
	}
	for _, o := range opts {
		o(l)
	}
	return l
}

// NewPrefixedRoundTripLogger creates a new HTTP transport that logs the raw
// (outgoing) HTTP request and response to the provided logger.
//
// Prefixes requests and responses with "-> " and "<-", respectively. Adds an
// additional blank line ("\n\n") to the output of requests and responses.
//
// Valid types for logger:
//
//	io.Writer
//	func(string, ...interface{}) (int, error) // fmt.Printf
//	func(string, ...interface{}) // log.Printf
//
// Note: will panic() when an unknown logger type is passed.
func NewPrefixedRoundTripLogger(transport http.RoundTripper, logger any, opts ...Option) *RoundTripLogger {
	nl := []byte("\n")
	var f func([]byte, []byte)
	switch logf := logger.(type) {
	case io.Writer:
		f = func(prefix []byte, buf []byte) {
			buf = append(prefix, bytes.ReplaceAll(buf, nl, append(nl, prefix...))...)
			_, _ = logf.Write(append(buf, '\n', '\n'))
		}
	case func(string, ...any) (int, error): // fmt.Printf
		f = func(prefix []byte, buf []byte) {
			buf = append(prefix, bytes.ReplaceAll(buf, nl, append(nl, prefix...))...)
			_, _ = logf("%s\n\n", string(buf))
		}
	case func(string, ...any): // log.Printf
		f = func(prefix []byte, buf []byte) {
			buf = append(prefix, bytes.ReplaceAll(buf, nl, append(nl, prefix...))...)
			logf("%s\n\n", string(buf))
		}
	default:
		panic(fmt.Sprintf("unable to convert logf with type %T", logf))
	}
	return NewRoundTripLogger(
		transport,
		func(buf []byte) { f([]byte("-> "), buf) },
		func(buf []byte) { f([]byte("<- "), buf) },
		opts...,
	)
}

// RoundTrip satisfies the [http.RoundTripper] interface.
func (l *RoundTripLogger) RoundTrip(req *http.Request) (*http.Response, error) {
	transport := l.transport
	if transport == nil {
		transport = DefaultTransport
	}
	reqBody, err := httputil.DumpRequestOut(req, !l.noReqBody)
	if err != nil {
		return nil, err
	}
	res, err := transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	resBody, err := httputil.DumpResponse(res, !l.noResBody)
	if err != nil {
		return nil, err
	}
	l.reqf(reqBody)
	l.resf(resBody)
	return res, err
}

// Option is a roundtrip logger option.
type Option func(*RoundTripLogger)

// WithReqResBody is a roundtrip logger option to set whether or not to log the
// request and response body. Useful when body content is binary.
func WithReqResBody(req, res bool) Option {
	return func(l *RoundTripLogger) {
		l.noReqBody, l.noResBody = !req, !res
	}
}
