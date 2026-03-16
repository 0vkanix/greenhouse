package main

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/0vkanix/greenlight/internal/assert"
)

func TestCreateMovieHandler(t *testing.T) {
	app := newTestApplication(t)
	server := newTestServer(t, app.routes())
	defer server.Close()

	code, _, body := server.post(t, "/v1/movies", nil)
	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, body, "create a new movie")
}

func TestShowMovieHandler(t *testing.T) {
	app := newTestApplication(t)
	server := newTestServer(t, app.routes())
	defer server.Close()

	tests := []struct {
		name       string
		movieID    string
		wantCode   int
		wantInBody []string
	}{
		{
			name:       "Valid request",
			movieID:    "1",
			wantCode:   http.StatusOK,
			wantInBody: []string{`"id":1`, `"title":"Casablanca"`, `"runtime":102`},
		},
		{
			name:     "Invalid ID parameter",
			movieID:  "-1",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Invalid ID (non-numeric)",
			movieID:  "abc",
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := server.get(t, fmt.Sprintf("/v1/movies/%s", tt.movieID))

			assert.Equal(t, code, tt.wantCode)
			for _, want := range tt.wantInBody {
				assert.StringContains(t, body, want)
			}
		})
	}
}
