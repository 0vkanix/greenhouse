package movie

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrRecordNotFound is returned when a requested movie record does not exist in the database.
var ErrRecordNotFound = errors.New("record not found")

// RepositoryInterface defines the contract for movie data access, allowing for 
// both concrete database implementations and mock implementations for testing.
type RepositoryInterface interface {
	Insert(ctx context.Context, movie *Movie) error
	Get(ctx context.Context, id uuid.UUID) (*Movie, error)
	Update(ctx context.Context, movie *Movie) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// Repository is a concrete implementation of RepositoryInterface using a 
// PostgreSQL connection pool and generated sqlc queries.
type Repository struct {
	pool    *pgxpool.Pool
	queries *Queries
}

// NewRepository creates and returns a new movie repository instance.
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool:    pool,
		queries: New(pool),
	}
}

// Delete removes a movie record from the database by its UUID. 
// It returns ErrRecordNotFound if the ID does not exist.
func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rowsAffected, err := r.queries.Delete(ctx, id)
	if err != nil {
		return err
	} else if rowsAffected.RowsAffected() == 0 {
		return ErrRecordNotFound
	}

	return nil
}

// Update modifies an existing movie record. It hydrates the movie pointer 
// with the new system-generated version number.
func (r *Repository) Update(ctx context.Context, movie *Movie) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	arg := UpdateParams{
		ID:      movie.ID,
		Title:   movie.Title,
		Year:    movie.Year,
		Runtime: movie.Runtime,
		Genres:  movie.Genres,
	}

	version, err := r.queries.Update(ctx, arg)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return ErrRecordNotFound
		default:
			return err
		}
	}

	movie.Version = version

	return nil
}

// Insert adds a new movie record to the database. It hydrates the movie pointer 
// with system-generated fields (ID, CreatedAt, Version).
func (r *Repository) Insert(ctx context.Context, movie *Movie) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	arg := InsertParams{
		Title:   movie.Title,
		Year:    movie.Year,
		Runtime: movie.Runtime,
		Genres:  movie.Genres,
	}

	result, err := r.queries.Insert(ctx, arg)
	if err != nil {
		return err
	}

	movie.ID = result.ID
	movie.CreatedAt = result.CreatedAt
	movie.Version = result.Version

	return nil
}

// Get retrieves a single movie record by its UUID. It returns ErrRecordNotFound 
// if no matching record is found.
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
