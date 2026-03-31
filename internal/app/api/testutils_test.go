package api

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/0vkanix/greenlight/internal/movie"
	"github.com/google/uuid"
)

// StubMovieRepository is a mock implementation of movie.RepositoryInterface 
// for use in API handlers tests.
type StubMovieRepository struct {
	Movies map[uuid.UUID]*movie.Movie
	Error  error
}

func (s *StubMovieRepository) Insert(ctx context.Context, movieObj *movie.Movie) error {
	if s.Error != nil {
		return s.Error
	}
	movieObj.ID = uuid.New()
	movieObj.Version = 1
	return nil
}

func (s *StubMovieRepository) Get(ctx context.Context, id uuid.UUID) (*movie.Movie, error) {
	if s.Error != nil {
		return nil, s.Error
	}
	movieObj, ok := s.Movies[id]
	if !ok {
		return nil, movie.ErrRecordNotFound
	}
	return movieObj, nil
}

func (s *StubMovieRepository) Update(ctx context.Context, movieObj *movie.Movie) error {
	if s.Error != nil {
		return s.Error
	}
	_, ok := s.Movies[movieObj.ID]
	if !ok {
		return movie.ErrRecordNotFound
	}
	s.Movies[movieObj.ID] = movieObj
	return nil
}

func (s *StubMovieRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if s.Error != nil {
		return s.Error
	}
	if _, ok := s.Movies[id]; !ok {
		return movie.ErrRecordNotFound
	}
	delete(s.Movies, id)
	return nil
}

func newTestApplication(t *testing.T, repo movie.RepositoryInterface) *Application {
	return &Application{
		Config: Config{
			Port: 4000,
			Env:  "test",
		},
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		Movies: repo,
	}
}

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	server := httptest.NewServer(h)

	return &testServer{server}
}

func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, string) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
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

func (ts *testServer) post(t *testing.T, urlPath string, data []byte) (int, http.Header, string) {
	rs, err := ts.Client().Post(ts.URL+urlPath, "application/json", bytes.NewReader(data))
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
