package main

import (
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := envelope{
		"status":      "available",
		"environment": app.config.env, "version": version,
	}
	_ = app.writeJSON(w, r, http.StatusOK, data, nil)
}
