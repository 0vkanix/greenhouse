package main

import (
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
		name     string
		movieID  string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid request",
			movieID:  "1",
			wantCode: http.StatusOK,
			wantBody: "show the details of movie 1",
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
			code, _, body := server.get(t, "/v1/movies/"+tt.movieID)

			assert.Equal(t, code, tt.wantCode)
			if tt.wantBody != "" {
				assert.Equal(t, body, tt.wantBody)
			}
		})
	}
}
