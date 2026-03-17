package main

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/0vkanix/greenlight/internal/assert"
)

func TestHealthcheckHandler(t *testing.T) {
	app := newTestApplication(t)
	server := newTestServer(t, app.routes())
	defer server.Close()

	code, headers, body := server.get(t, "/v1/healthcheck")

	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, headers.Get("Content-Type"), "application/json")

	var got struct {
		Status      string `json:"status"`
		Environment string `json:"environment"`
		Version     string `json:"version"`
	}

	err := json.Unmarshal([]byte(body), &got)
	if err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	assert.Equal(t, got.Status, "available")
	assert.Equal(t, got.Environment, "test")
	assert.Equal(t, got.Version, "1.0.0")
}
