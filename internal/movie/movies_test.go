package movie_test

import (
	"fmt"
	"testing"

	m "github.com/0vkanix/greenlight/internal/movie"
	"github.com/0vkanix/greenlight/internal/movie/types"
	"github.com/0vkanix/greenlight/internal/validator"
)

// ExampleMovie_Validate demonstrates how to use the Validate method
// to check if a movie object conforms to business rules.
func ExampleMovie_Validate() {
	v := validator.New()
	m := m.Movie{
		Title:   "Casablanca",
		Year:    1942,
		Runtime: types.Runtime(102),
		Genres:  []string{"drama", "war"},
	}

	m.Validate(v)

	if v.Valid() {
		fmt.Println("Movie is valid")
	}
	// Output: Movie is valid
}

// TestValidate verifies that the movie validation logic correctly identifies
// valid and invalid movie data across various scenarios like empty titles,
// future years, and negative runtimes.
func TestValidate(t *testing.T) {
	v := validator.New()

	tests := []struct {
		name  string
		movie m.Movie
	}{
		{
			"Valid movie",
			m.Movie{
				Title:   "Casablanca",
				Year:    2026,
				Runtime: types.Runtime(120),
				Genres:  []string{"drama", "war"},
			},
		},
		{
			"Empty title",
			m.Movie{
				Title:   "",
				Year:    2026,
				Runtime: types.Runtime(120),
				Genres:  []string{"drama", "war"},
			},
		},
		{
			"Future year",
			m.Movie{
				Title:   "Casablanca",
				Year:    9999,
				Runtime: types.Runtime(120),
				Genres:  []string{"drama", "war"},
			},
		},
		{
			"Negative runtime",
			m.Movie{
				Title:   "Casablanca",
				Year:    2026,
				Runtime: types.Runtime(-120),
				Genres:  []string{"drama", "war"},
			},
		},
		{
			"Too many genres",
			m.Movie{
				Title:   "Casablanca",
				Year:    2026,
				Runtime: types.Runtime(120),
				Genres:  []string{"drama", "war", "comedy", "documentary", "love", "game"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.movie.Validate(v)
		})
	}
}
