package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

var ErrInvalidIdParam = errors.New("invalid id parameter")

func (app *application) readIDParam(r *http.Request) (int64, error) {
	idParam := chi.URLParamFromCtx(r.Context(), "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || id < 1 {
		return 0, ErrInvalidIdParam
	}

	return id, nil
}
