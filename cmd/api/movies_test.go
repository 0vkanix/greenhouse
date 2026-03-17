package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/0vkanix/greenlight/internal/assert"
	"github.com/0vkanix/greenlight/internal/data"
)

func TestCreateMovieHandler(t *testing.T) {
	app := newTestApplication(t)
	server := newTestServer(t, app.routes())
	defer server.Close()

	tests := []struct {
		name      string
		input     string
		wantCode  int
		wantTitle string
	}{
		{
			name:      "Valid request",
			input:     `{"title":"Casablanca","year":1942,"runtime":102,"genres":["drama","war"],"version":1}`,
			wantCode:  http.StatusOK,
			wantTitle: "Casablanca",
		},
		{
			name:     "Malformed JSON (syntax error)",
			input:    `{"title":"Casablanca","year":1942,`,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "Invalid field type",
			input:    `{"title":"Casablanca","runtime":"102 mins"}`,
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := server.post(t, "/v1/movies", []byte(tt.input))

			assert.Equal(t, code, tt.wantCode)

			if tt.wantTitle != "" {
				assert.StringContains(t, body, tt.wantTitle)
			}
		})
	}
}

func TestShowMovieHandler(t *testing.T) {
	app := newTestApplication(t)
	server := newTestServer(t, app.routes())
	defer server.Close()

	tests := []struct {
		name     string
		movieID  string
		wantCode int
		wantData *data.Movie
	}{
		{
			name:     "Valid request",
			movieID:  "1",
			wantCode: http.StatusOK,
			wantData: &data.Movie{
				ID:      1,
				Title:   "Casablanca",
				Runtime: 102,
			},
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

			if tt.wantData != nil {
				var got struct {
					Movie *data.Movie `json:"movie"`
				}

				err := json.NewDecoder(bytes.NewBufferString(body)).Decode(&got)
				if err != nil {
					t.Fatal(err)
				}

				assert.Equal(t, got.Movie.ID, tt.wantData.ID)
				assert.Equal(t, got.Movie.Title, tt.wantData.Title)
				assert.Equal(t, got.Movie.Runtime, tt.wantData.Runtime)
			}
		})
	}
}
