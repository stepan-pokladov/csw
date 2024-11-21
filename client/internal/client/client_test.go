package client

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSend(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		responseBody   string
		expectedError  bool
		expectedOutput string
	}{
		{
			name:           "Successful Request",
			statusCode:     http.StatusOK,
			responseBody:   "Success",
			expectedError:  false,
			expectedOutput: "Report sent to",
		},
		{
			name:           "Server Error",
			statusCode:     http.StatusInternalServerError,
			responseBody:   "Internal Server Error",
			expectedError:  true,
			expectedOutput: "failed to send report",
		},
		{
			name:           "Bad Request",
			statusCode:     http.StatusBadRequest,
			responseBody:   "Bad Request",
			expectedError:  true,
			expectedOutput: "failed to send report",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Check headers
				if r.Header.Get("Content-Encoding") != "gzip" {
					t.Errorf("Expected Content-Encoding header to be 'gzip', got %s", r.Header.Get("Content-Encoding"))
				}
				if r.Header.Get("Content-Type") != "application/jsonl" {
					t.Errorf("Expected Content-Type header to be 'application/jsonl', got %s", r.Header.Get("Content-Type"))
				}
				// Write the response
				w.WriteHeader(tt.statusCode)
				fmt.Fprintln(w, tt.responseBody)
			}))
			defer server.Close()

			// Call the send function
			data := []byte(`{"test": "data"}`)
			err := send(server.URL, data)

			// Validate error state
			if (err != nil) != tt.expectedError {
				t.Fatalf("Expected error: %v, got: %v", tt.expectedError, err)
			}

			// Validate error message or output
			if err != nil && !bytes.Contains([]byte(err.Error()), []byte(tt.expectedOutput)) {
				t.Errorf("Expected error output to contain '%s', got '%s'", tt.expectedOutput, err.Error())
			}
		})
	}
}
