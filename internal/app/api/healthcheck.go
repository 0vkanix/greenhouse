package api

import (
	"net/http"
)

func (app *Application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := envelope{
		"status":      "available",
		"environment": app.Config.Env, "version": version,
	}
	_ = app.writeJSON(w, r, http.StatusOK, data, nil)
}
