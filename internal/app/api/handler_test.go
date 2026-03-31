package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
			stubRepo := &StubMovieRepository{Error: tt.stubError}
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
		stubRepo  *StubMovieRepository
		wantCode  int
		wantData  *movie.Movie
		wantError string
	}{
		{
			name:    "Valid request",
			movieID: id.String(),
			stubRepo: &StubMovieRepository{
				Movies: map[uuid.UUID]*movie.Movie{id: wantMovie},
			},
			wantCode: http.StatusOK,
			wantData: wantMovie,
		},
		{
			name:     "Invalid ID parameter (non-UUID)",
			movieID:  "123",
			stubRepo: &StubMovieRepository{},
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Non-existent movie",
			movieID:  uuid.New().String(),
			stubRepo: &StubMovieRepository{Movies: map[uuid.UUID]*movie.Movie{}},
			wantCode: http.StatusNotFound,
		},
		{
			name:    "Database error",
			movieID: id.String(),
			stubRepo: &StubMovieRepository{
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

func TestUpdateMovieHandler(t *testing.T) {
	id := uuid.New()
	movieObj := &movie.Movie{
		ID:      id,
		Title:   "Casablanca",
		Year:    1942,
		Runtime: 102,
		Genres:  []string{"drama", "war"},
		Version: 1,
	}

	tests := []struct {
		name      string
		movieID   string
		input     string
		stubRepo  *StubMovieRepository
		wantCode  int
		wantTitle string
	}{
		{
			name:    "Valid request",
			movieID: id.String(),
			input:   `{"title":"Casablanca (Updated)","year":1942,"runtime":"102 mins","genres":["drama"]}`,
			stubRepo: &StubMovieRepository{
				Movies: map[uuid.UUID]*movie.Movie{id: movieObj},
			},
			wantCode:  http.StatusOK,
			wantTitle: "Casablanca (Updated)",
		},
		{
			name:     "Invalid ID parameter (non-UUID)",
			movieID:  "abc",
			input:    `{"title":"Updated"}`,
			stubRepo: &StubMovieRepository{},
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Malformed JSON",
			movieID:  id.String(),
			input:    `{"title":}`,
			stubRepo: &StubMovieRepository{},
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "Validation failure",
			movieID:  id.String(),
			input:    `{"title":""}`,
			stubRepo: &StubMovieRepository{},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:     "Non-existent movie",
			movieID:  uuid.New().String(),
			input:    `{"title":"Updated","year":1942,"runtime":"102 mins","genres":["drama"]}`,
			stubRepo: &StubMovieRepository{Movies: map[uuid.UUID]*movie.Movie{}},
			wantCode: http.StatusNotFound,
		},
		{
			name:    "Database error",
			movieID: id.String(),
			input:   `{"title":"Updated","year":1942,"runtime":"102 mins","genres":["drama"]}`,
			stubRepo: &StubMovieRepository{
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

			url := fmt.Sprintf("/v1/movies/%s", tt.movieID)
			code, _, body := server.put(t, url, []byte(tt.input))

			assert.Equal(t, code, tt.wantCode)

			if tt.wantTitle != "" {
				assert.StringContains(t, body, tt.wantTitle)
			}
		})
	}
}

func (ts *testServer) put(t *testing.T, urlPath string, data []byte) (int, http.Header, string) {
	r, err := http.NewRequest(http.MethodPut, ts.URL+urlPath, bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	r.Header.Set("Content-Type", "application/json")

	rs, err := ts.Client().Do(r)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	body = bytes.TrimSpace(body)
	return rs.StatusCode, rs.Header, string(body)
}
