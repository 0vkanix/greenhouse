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

	code, _, body := server.get(t, "/v1/healthcheck")

	assert.Equal(t, code, http.StatusOK)
	assert.StringContains(t, body, "status: available")
	assert.StringContains(t, body, "environment: test")
	assert.StringContains(t, body, "version: 1.0.0")
}
