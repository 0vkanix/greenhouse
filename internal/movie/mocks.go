package movie

import (
	"context"

	"github.com/google/uuid"
)

type StubMovieRepository struct {
	Movies map[uuid.UUID]*Movie
	Error  error
}

func (s *StubMovieRepository) Insert(ctx context.Context, movie *Movie) error {
	if s.Error != nil {
		return s.Error
	}

	movie.ID = uuid.New()
	movie.Version = 1
	return nil
}

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
