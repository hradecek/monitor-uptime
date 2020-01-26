package main

import (
	"net/http"
	"net/http/httptrace"
	"time"
)

type Uptime struct {
	StatusCode int
	TTFB time.Duration
}
// GetUptime returns HTTP status code and time to first byte (TTFB) for provided host
func GetUptime(host string, timeout int) (*Uptime, error) {
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
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}
	start = time.Now()
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return &Uptime{StatusCode: res.StatusCode, TTFB: first_byte / time.Millisecond}, nil
}
