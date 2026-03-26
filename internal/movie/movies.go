package movie

import (
	"time"

	"github.com/0vkanix/greenlight/internal/validator"
)

func (m *Movie) Validate(v *validator.Validator) {
	v.Check(m.Title != "", "title", "must be provided")
	v.Check(len(m.Title) < 500, "title", "must not be more than 500 bytes long")

	v.Check(m.Year != 0, "year", "must be provided")
	v.Check(m.Year >= 1888, "year", "must be greater than 1888")
	v.Check(m.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	v.Check(m.Runtime != 0, "runtime", "must be provided")
	v.Check(m.Runtime > 0, "runtime", "must be a positive integer")

	v.Check(m.Genres != nil, "genres", "must be provided")
	v.Check(len(m.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(m.Genres) <= 5, "genres", "must not contain more than 5 genres")

	v.Check(validator.Unique(m.Genres), "genres", "must not contain duplicate values")
}
