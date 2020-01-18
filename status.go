package main

import (
	"net/http"
	"net/http/httptrace"
	"time"
)

type Status struct {
	StatusCode int
	TTFB time.Duration
}
// GetStatus returns HTTP status code and time to first byte (TTFB) for provided host
func GetStatus(host string) (*Status, error) {
	var start time.Time
	var first_byte time.Duration
	req, err := http.NewRequest("GET", host, nil)
	if err != nil {
		return nil, err
	}
	trace := &httptrace.ClientTrace{
		GotFirstResponseByte: func() {
            first_byte = time.Since(start)
        },
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	start = time.Now()
	res, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	return &Status{StatusCode: res.StatusCode, TTFB: first_byte / time.Millisecond}, nil
}