package uptime

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Given host is up,
// When uptime is retrieved,
// Then uptime result contains HTTP OK (200)
func TestGetUptimeWebsiteUp(t *testing.T) {
	// Given
	hostHTTP := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(1 * time.Second)
			w.WriteHeader(http.StatusOK)
	}))
	defer hostHTTP.Close()

	// When
	result, err := GetUptime(hostHTTP.URL, 10)

	// Then
	assert.Nil(t, err, "Unexpected error happened")
	assert.Equal(t, 200, result.StatusCode, "Unexpected status code")
	assert.GreaterOrEqual(t, int64(result.TTFB), int64(0), "Unexpected TTFB value")
	assert.GreaterOrEqual(t, int64(result.DNSLookup), int64(0), "Unexpected DNSLookup value")
	assert.GreaterOrEqual(t, int64(result.TLSHandshake), int64(0), "Unexpected TLSHandshake value")
}

// Given host is up,
// When uptime is retrieved
//      and host response by specific HTTP status code
// Then uptime result contains exactly that HTTP status code
func TestGetUptimeHTTPStatusCode(t *testing.T) {
	hostHTTP := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
	defer hostHTTP.Close()

	// When
	result, err := GetUptime(hostHTTP.URL, 10)

	// Then
	assert.Nil(t, err, "Unexpected error happened")
	assert.Equal(t, http.StatusNotFound, result.StatusCode, "Unexpected HTTP status code")
}

// Given host is up
//       and its response time is slow (5 seconds)
// When uptime is retrieved and timeout (4 seconds) is reached
// Then error is returned
func TestGetUptimeTimeout(t *testing.T) {
	// Given
	hostHTTP := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(5 * time.Second)
		}))
	defer hostHTTP.Close()

	// When
	_, err := GetUptime(hostHTTP.URL, 4)

	// Then
	assert.NotNil(t, err, "Error was expected")
}

// When uptime is retrieved
//      and provided host URL does not exists,
// Then error is returned
func TestGetUptimeNonExistingHostURL(t *testing.T) {
	// When
	_, err := GetUptime("non-existing-url", 10)

	// Then
	assert.NotNil(t, err, "Error was expected")
}

// When uptime is retrieved
//      and provided host URL is malformed,
// Then error is returned
func TestGetUptimeMalformedHostURL(t *testing.T) {
	// When
	_, err := GetUptime(string([]byte{00}), 10)

	// Then
	assert.NotNil(t, err, "Error was expected")
}
