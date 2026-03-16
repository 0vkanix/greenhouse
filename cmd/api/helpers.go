package main

import (
	"encoding/json"
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

func (app *application) writeJSON(w http.ResponseWriter, r *http.Request, status int, data any, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	for key, values := range headers {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}
