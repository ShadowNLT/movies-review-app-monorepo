package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	// We convert our helpers to http.Handler using the http.HandlerFunc
	// and making them the custom error handler for 404 and 405 responses
	// that could originate from the router itself
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// healthcheck
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	// movie Reviews
	router.HandlerFunc(http.MethodPost, "/v1/reviews", app.createMovieReviewHandler)
	router.HandlerFunc(http.MethodGet, "/v1/reviews/:id", app.showMovieReviewHandler)

	return app.recoverPanic(router)
}
