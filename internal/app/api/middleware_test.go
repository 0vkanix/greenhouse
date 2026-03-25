package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/0vkanix/greenlight/internal/assert"
)

func TestRecoverPanic(t *testing.T) {
	app := newTestApplication(t)

	// Create a handler that panics
	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("intentional panic")
	})

	// Wrap it with our recoverPanic middleware
	h := app.recoverPanic(panicHandler)

	rr := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	h.ServeHTTP(rr, r)

	// Verify the response
	assert.Equal(t, rr.Code, http.StatusInternalServerError)
	assert.Equal(t, rr.Header().Get("Connection"), "close")
	assert.StringContains(t, rr.Body.String(), ErrInternalServerError)
}
