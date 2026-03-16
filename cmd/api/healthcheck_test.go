package main

import (
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
	assert.Equal(t, body, `{"status": "available", "environment": "test", "version": "1.0.0"}`)
}
