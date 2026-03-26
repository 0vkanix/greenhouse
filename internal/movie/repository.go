package movie

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrRecordNotFound = errors.New("record not found")

type RepositoryInterface interface {
	Insert(ctx context.Context, movie *Movie) error
	Get(ctx context.Context, id uuid.UUID) (*Movie, error)
}

type Repository struct {
	pool    *pgxpool.Pool
	queries *Queries
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool:    pool,
		queries: New(pool),
	}
}

func (r *Repository) Insert(ctx context.Context, movie *Movie) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	arg := InsertParams{
		Title:   movie.Title,
		Year:    movie.Year,
		Runtime: movie.Runtime,
		Genres:  movie.Genres,
	}

	row, err := r.queries.Insert(ctx, arg)
	if err != nil {
		return err
	}

	movie.ID = row.ID
	movie.CreatedAt = row.CreatedAt
	movie.Version = row.Version

	return nil
}

func (r *Repository) Get(ctx context.Context, id uuid.UUID) (*Movie, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	movie, err := r.queries.Get(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &movie, nil
}
