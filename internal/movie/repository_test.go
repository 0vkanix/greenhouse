package movie_test

import (
	"context"
	"os"
	"testing"

	"github.com/0vkanix/greenlight/internal/assert"
	"github.com/0vkanix/greenlight/internal/movie"
	"github.com/0vkanix/greenlight/internal/movie/types"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// newTestDB creates a connection to the test database, sets up the schema,
// and registers a teardown function.
func newTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	dsn := "postgres://db_test:test@localhost:5432/greenlight_test?sslmode=disable"
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		t.Fatal(err)
	}

	script, err := os.ReadFile("./testdata/setup.sql")
	if err != nil {
		pool.Close()
		t.Fatal(err)
	}

	_, err = pool.Exec(context.Background(), string(script))
	if err != nil {
		pool.Close()
		t.Fatal(err)
	}

	t.Cleanup(func() {
		defer pool.Close()

		script, err := os.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}

		_, err = pool.Exec(context.Background(), string(script))
		if err != nil {
			t.Fatal(err)
		}
	})

	return pool
}

func TestRepository_Get(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pool := newTestDB(t)
	repo := movie.NewRepository(pool)

	ctx := context.Background()

	// Seeded Movie: Moana
	id := uuid.MustParse("0ee58bdc-e2de-454f-ad22-e02caa53cc31")

	m, err := repo.Get(ctx, id)
	assert.NilError(t, err)

	assert.Equal(t, m.Title, "Moana")
	assert.Equal(t, int(m.Year), 2016)
}

func TestRepository_Insert(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pool := newTestDB(t)
	repo := movie.NewRepository(pool)

	ctx := context.Background()

	m := &movie.Movie{
		Title:   "Inception",
		Year:    2010,
		Runtime: types.Runtime(148),
		Genres:  []string{"action", "sci-fi"},
	}

	err := repo.Insert(ctx, m)
	assert.NilError(t, err)

	assert.NotEqual(t, m.ID, uuid.Nil)
	assert.NotEqual(t, m.CreatedAt.IsZero(), true)
	assert.Equal(t, int(m.Version), 1)

	// Verify it actually exists in the DB
	row, err := repo.Get(ctx, m.ID)
	assert.NilError(t, err)
	assert.Equal(t, row.Title, "Inception")
}

func TestRepository_Update(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pool := newTestDB(t)
	repo := movie.NewRepository(pool)

	ctx := context.Background()

	// Seeded Movie: Deadpool
	id := uuid.MustParse("72dda126-efbc-41f0-8c1d-14141484f540")

	m, err := repo.Get(ctx, id)
	assert.NilError(t, err)

	m.Title = "Deadpool (Updated)"
	err = repo.Update(ctx, m)
	assert.NilError(t, err)

	assert.Equal(t, int(m.Version), 2)

	row, err := repo.Get(ctx, m.ID)
	assert.NilError(t, err)
	assert.Equal(t, row.Title, "Deadpool (Updated)")
}

func TestRepository_Delete(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pool := newTestDB(t)
	repo := movie.NewRepository(pool)

	ctx := context.Background()

	// Seeded Movie: Black Panther
	id := uuid.MustParse("6350d919-27a6-4b7b-ae8d-4f744bdf0282")

	err := repo.Delete(ctx, id)
	assert.NilError(t, err)

	_, err = repo.Get(ctx, id)
	if err == nil {
		t.Error("expected error but got nil when fetching deleted movie")
	}
	assert.Equal(t, err, movie.ErrRecordNotFound)
}

func TestRepository_Get_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pool := newTestDB(t)
	repo := movie.NewRepository(pool)

	_, err := repo.Get(context.Background(), uuid.New())
	assert.Equal(t, err, movie.ErrRecordNotFound)
}
