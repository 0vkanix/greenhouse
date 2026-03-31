package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/0vkanix/greenlight/internal/assert"
)

func TestRecoverPanic(t *testing.T) {
	stubRepo := &StubMovieRepository{}
	app := newTestApplication(t, stubRepo)

	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("intentional panic")
	})

	h := app.recoverPanic(panicHandler)

	rr := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	h.ServeHTTP(rr, r)

	// Verify the response
	assert.Equal(t, rr.Code, http.StatusInternalServerError)
	assert.Equal(t, rr.Header().Get("Connection"), "close")
	assert.StringContains(t, rr.Body.String(), ErrInternalServerError)
}
