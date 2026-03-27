package movie

import (
	"context"

	"github.com/google/uuid"
)

// StubMovieRepository is a mock implementation of movie.RepositoryInterface 
// designed for use in unit tests. It allows tests to control the data 
// and errors returned by the movie domain without a real database.
type StubMovieRepository struct {
	Movies map[uuid.UUID]*Movie
	Error  error
}

// Insert simulates a successful database insertion by generating a 
// new UUID and version for the provided movie.
func (s *StubMovieRepository) Insert(ctx context.Context, movie *Movie) error {
	if s.Error != nil {
		return s.Error
	}

	movie.ID = uuid.New()
	movie.Version = 1

	return nil
}

// Get simulates retrieving a movie record. It returns a movie from the 
// internal map or ErrRecordNotFound if the UUID is not present.
func (s *StubMovieRepository) Get(ctx context.Context, id uuid.UUID) (*Movie, error) {
	if s.Error != nil {
		return nil, s.Error
	}

	movie, ok := s.Movies[id]
	if !ok {
		return nil, ErrRecordNotFound
	}

	return movie, nil
}

// Update simulates a successful database update by updating the 
// internal map. It returns ErrRecordNotFound if the ID is missing.
func (s *StubMovieRepository) Update(ctx context.Context, movie *Movie) error {
	if s.Error != nil {
		return s.Error
	}

	_, ok := s.Movies[movie.ID]
	if !ok {
		return ErrRecordNotFound
	}

	s.Movies[movie.ID] = movie

	return nil
}

// Delete simulates a database deletion by removing the ID from the 
// internal map. It returns ErrRecordNotFound if the ID is missing.
func (s *StubMovieRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if s.Error != nil {
		return s.Error
	}

	if _, ok := s.Movies[id]; !ok {
		return ErrRecordNotFound
	}

	delete(s.Movies, id)

	return nil
}
