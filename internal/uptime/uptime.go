package uptime

import (
	"net/http"
	"net/http/httptrace"
	"time"
)

// Represents result of single uptime monitor run
type Result struct {
	StatusCode int
	TTFB       time.Duration // Measured 'time to first byte'
}

// Creates single HTTP request and collects uptime monitor's metrics that are returned as result
// In case of failure error is returned instead
func GetUptime(host string, timeout int) (*Result, error) {
	var startTime time.Time
	var firstByteTime time.Duration

	req, err := http.NewRequest("GET", host, nil)
	if err != nil {
		return nil, err
	}

	trace := &httptrace.ClientTrace{
		GetConn: func(_ string) {
			startTime = time.Now()
		},
		GotFirstResponseByte: func() {
			firstByteTime = time.Since(startTime)
		},
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return &Result{
		StatusCode: res.StatusCode,
		TTFB:       firstByteTime / time.Millisecond,
	}, nil
}
