package uptime

import (
	"crypto/tls"
	"net/http"
	"net/http/httptrace"
	"time"
)

// Represents result of single uptime monitor run
type Result struct {
	StatusCode   int
	TTFB         time.Duration // Measured Time To First Byte
	DNSLookup    time.Duration // Measured duration of DNS lookup
	TLSHandshake time.Duration // Measured duration of TLS handshake
}

// Creates single HTTP request and collects uptime monitor's metrics that are returned as result
// In case of failure error is returned instead
func GetUptime(host string, timeout int) (*Result, error) {
	var connStartTime, dnsStartTime, tlsStartTime time.Time
	var firstByteDuration, dnsDuration, tlsDuration time.Duration

	req, err := http.NewRequest("GET", host, nil)
	if err != nil {
		return nil, err
	}

	trace := &httptrace.ClientTrace{
		GetConn: func(_ string) {
			connStartTime = time.Now()
		},
		GotFirstResponseByte: func() {
			firstByteDuration = time.Since(connStartTime)
		},
		DNSStart: func(_ httptrace.DNSStartInfo) {
			dnsStartTime = time.Now()
		},
		DNSDone: func(_ httptrace.DNSDoneInfo) {
			dnsDuration = time.Since(dnsStartTime)
		},
		TLSHandshakeStart: func() {
			tlsStartTime = time.Now()
		},
		TLSHandshakeDone: func(_ tls.ConnectionState, _ error) {
			tlsDuration = time.Since(tlsStartTime)
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
		StatusCode:   res.StatusCode,
		TTFB:         firstByteDuration.Round(time.Millisecond),
		DNSLookup:    dnsDuration.Round(time.Millisecond),
		TLSHandshake: tlsDuration.Round(time.Millisecond),
	}, nil
}
