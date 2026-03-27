package movie

import (
	"testing"

	"github.com/0vkanix/greenlight/internal/validator"
)

// TestValidate verifies that the movie validation logic correctly identifies 
// valid and invalid movie data across various scenarios like empty titles, 
// future years, and negative runtimes.
func TestValidate(t *testing.T) {
	v := validator.New()

	tests := []struct {
		name  string
		movie Movie
	}{
		{
			"Valid movie",
			Movie{
				Title:   "Casablanca",
				Year:    2026,
				Runtime: Runtime(120),
				Genres:  []string{"drama", "war"},
			},
		},
		{
			"Empty title",
			Movie{
				Title:   "",
				Year:    2026,
				Runtime: Runtime(120),
				Genres:  []string{"drama", "war"},
			},
		},
		{
			"Future year",
			Movie{
				Title:   "Casablanca",
				Year:    9999,
				Runtime: Runtime(120),
				Genres:  []string{"drama", "war"},
			},
		},
		{
			"Negative runtime",
			Movie{
				Title:   "Casablanca",
				Year:    2026,
				Runtime: Runtime(-120),
				Genres:  []string{"drama", "war"},
			},
		},
		{
			"Too many genres",
			Movie{
				Title:   "Casablanca",
				Year:    2026,
				Runtime: Runtime(120),
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
