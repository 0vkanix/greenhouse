package main

import (
	"net/http"

	"github.com/0vkanix/greenlight/internal/errors"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := envelope{
		"status":      "available",
		"environment": app.config.env, "version": version,
	}
	err := app.writeJSON(w, r, http.StatusOK, data, nil)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, errors.ErrInternalServerError.Error(), http.StatusInternalServerError)
	}
}
