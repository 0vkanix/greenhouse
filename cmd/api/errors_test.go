package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/0vkanix/greenlight/internal/assert"
)

func TestMethodNotAllowedResponse(t *testing.T) {
	app := newTestApplication(t)
	server := newTestServer(t, app.routes())
	defer server.Close()

	code, _, _ := server.post(t, "/v1/healthcheck", nil)

	assert.Equal(t, code, http.StatusMethodNotAllowed)
}

func TestNotFoundResponse(t *testing.T) {
	app := newTestApplication(t)
	server := newTestServer(t, app.routes())
	defer server.Close()

	code, _, _ := server.get(t, "/v1/movie/-1")

	assert.Equal(t, code, http.StatusNotFound)
}

func TestBadRequestResponse(t *testing.T) {
	app := newTestApplication(t)
	rr := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/", nil)

	err := errors.New("bad request error")
	app.badRequestResponse(rr, r, err)

	assert.Equal(t, rr.Code, http.StatusBadRequest)
	assert.StringContains(t, rr.Body.String(), "bad request error")
}

func TestServerErrorResponse(t *testing.T) {
	app := newTestApplication(t)
	rr := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	app.serverErrorResponse(rr, r, errors.New("test error"))

	assert.Equal(t, rr.Code, http.StatusInternalServerError)
	assert.StringContains(t, rr.Body.String(), ErrInternalServerError)
}

func TestErrorResponse(t *testing.T) {
	app := newTestApplication(t)

	rr := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	app.errorResponse(rr, r, http.StatusInternalServerError, make(chan int))

	assert.Equal(t, rr.Code, http.StatusInternalServerError)
}

func TestFailedValidationResponse(t *testing.T) {
	app := newTestApplication(t)
	rr := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	err := make(map[string]string)

	err["error"] = "failed validation error"

	app.failedValidationResponse(rr, r, err)

	assert.Equal(t, rr.Code, http.StatusUnprocessableEntity)
}
