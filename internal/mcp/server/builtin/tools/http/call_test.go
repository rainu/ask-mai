package http

import (
	"context"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCallDescriptor_Run(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "test-value", r.Header.Get("X-Test-Header"))

		// Request-Body lesen
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		defer r.Body.Close()

		assert.Equal(t, `{"test":"data"}`, string(body))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"Success"}`))
	}))
	defer server.Close()

	callDescriptor := CallDescriptor{
		Method: http.MethodPost,
		Url:    server.URL,
		Header: map[string]string{
			"X-Test-Header": "test-value",
			"Content-Type":  "application/json",
		},
		Body: strings.NewReader(`{"test":"data"}`),
	}

	result, err := callDescriptor.Run(context.Background(), http.DefaultClient)
	if err != nil {
		t.Fatalf("Fehler bei der Ausführung des Requests: %v", err)
	}
	delete(result.Header, "Date")

	assert.Equal(t, &CallResult{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Header: map[string][]string{
			"Content-Length": {"21"},
			"Content-Type":   {"application/json"},
		},
		Body: `{"message":"Success"}`,
	}, result)
}
