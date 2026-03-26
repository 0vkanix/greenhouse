package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/0vkanix/greenlight/internal/assert"
	"github.com/0vkanix/greenlight/internal/movie"
	"github.com/google/uuid"
)

func TestCreateMovieHandler(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		stubError error
		wantCode  int
		wantTitle string
	}{
		{
			name:      "Valid request",
			input:     `{"title":"Casablanca","year":1942,"runtime":"102 mins","genres":["drama","war"]}`,
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
			input:    `{"title":"Casablanca","year":"Year 1942"}`,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "Validation failure (empty title)",
			input:    `{"title":"","year":1942,"runtime":"102 mins","genres":["drama"]}`,
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:      "Database error",
			input:     `{"title":"Casablanca","year":1942,"runtime":"102 mins","genres":["drama"]}`,
			stubError: errors.New("database error"),
			wantCode:  http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stubRepo := &movie.StubMovieRepository{Error: tt.stubError}
			app := newTestApplication(t, stubRepo)
			server := newTestServer(t, app.routes())
			defer server.Close()

			code, _, body := server.post(t, "/v1/movies", []byte(tt.input))

			assert.Equal(t, code, tt.wantCode)

			if tt.wantTitle != "" {
				assert.StringContains(t, body, tt.wantTitle)
			}
		})
	}
}

func TestShowMovieHandler(t *testing.T) {
	id := uuid.New()
	wantMovie := &movie.Movie{ID: id, Title: "Casablanca"}

	tests := []struct {
		name      string
		movieID   string
		stubRepo  *movie.StubMovieRepository
		wantCode  int
		wantData  *movie.Movie
		wantError string
	}{
		{
			name:    "Valid request",
			movieID: id.String(),
			stubRepo: &movie.StubMovieRepository{
				Movies: map[uuid.UUID]*movie.Movie{id: wantMovie},
			},
			wantCode: http.StatusOK,
			wantData: wantMovie,
		},
		{
			name:     "Invalid ID parameter (non-UUID)",
			movieID:  "123",
			stubRepo: &movie.StubMovieRepository{},
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Non-existent movie",
			movieID:  uuid.New().String(),
			stubRepo: &movie.StubMovieRepository{Movies: map[uuid.UUID]*movie.Movie{}},
			wantCode: http.StatusNotFound,
		},
		{
			name:    "Database error",
			movieID: id.String(),
			stubRepo: &movie.StubMovieRepository{
				Error: errors.New("database error"),
			},
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := newTestApplication(t, tt.stubRepo)
			server := newTestServer(t, app.routes())
			defer server.Close()

			code, _, body := server.get(t, fmt.Sprintf("/v1/movies/%s", tt.movieID))

			assert.Equal(t, code, tt.wantCode)

			if tt.wantData != nil {
				var got struct {
					Movie *movie.Movie `json:"movie"`
				}

				err := json.NewDecoder(bytes.NewBufferString(body)).Decode(&got)
				if err != nil {
					t.Fatal(err)
				}

				assert.Equal(t, got.Movie.ID, tt.wantData.ID)
				assert.Equal(t, got.Movie.Title, tt.wantData.Title)
			}
		})
	}
}
