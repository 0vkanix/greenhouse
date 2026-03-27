package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *Application) routes() http.Handler {
	router := chi.NewRouter()

	router.NotFound(http.HandlerFunc(app.notFoundResponse))
	router.MethodNotAllowed(http.HandlerFunc(app.methodNotAllowedResponse))
	router.Get("/v1/healthcheck", app.healthcheckHandler)
	router.Post("/v1/movies", app.createMovieHandler)
	router.Get("/v1/movies/{id}", app.showMovieHandler)
	router.Put("/v1/movies/{id}", app.updateMovieHandler)
	router.Delete("/v1/movies/{id}", app.deleteMovieHandler)

	return app.recoverPanic(router)
}
