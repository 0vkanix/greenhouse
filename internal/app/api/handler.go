package api

import (
	"errors"
	"fmt"
	"net/http"

	m "github.com/0vkanix/greenlight/internal/movie"
	"github.com/0vkanix/greenlight/internal/validator"
)

func (app *Application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string    `json:"title"`
		Year    int32     `json:"year"`
		Runtime m.Runtime `json:"runtime"`
		Genres  []string  `json:"genres"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	movie := &m.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	v := validator.New()

	if movie.Validate(v); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.Movies.Insert(r.Context(), movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%v", movie.ID))

	err = app.writeJSON(w, r, http.StatusOK, envelope{"movie": movie}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movie, err := app.Movies.Get(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, m.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, r, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
