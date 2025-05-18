package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	// healthcheck
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	// movie Reviews
	router.HandlerFunc(http.MethodPost, "/v1/reviews", app.createMovieReviewHandler)
	router.HandlerFunc(http.MethodGet, "/v1/reviews/:id", app.showMovieReviewHandler)

	return router
}
