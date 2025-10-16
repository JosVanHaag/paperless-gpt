package main

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDocument containing extra parameters for testing
type TestDocument struct {
	ID         int
	Title      string
	Tags       []string
	FailUpdate bool // simulate update failure
}

// Use this for TestCases in your tests
type TestCase struct {
	name           string
	documents      []TestDocument
	expectedCount  int
	expectedError  string
	updateResponse int // HTTP status code for update response
}

// Test our HTTP-Client
func TestCreateCustomHTTPClient(t *testing.T) {
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Skipf("Skipping HTTP client test: %v", err)
	}

	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify custom header
		assert.Equal(t, "paperless-gpt", r.Header.Get("X-Title"), "Expected X-Title header")
		w.WriteHeader(http.StatusOK)
	}))
	server.Listener = listener
	server.Start()
	defer server.Close()

	// Get custom client
	client := createCustomHTTPClient()
	require.NotNil(t, client, "HTTP client should not be nil")

	// Make a request
	resp, err := client.Get(server.URL)
	require.NoError(t, err, "Request should not fail")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected 200 OK response")
}
