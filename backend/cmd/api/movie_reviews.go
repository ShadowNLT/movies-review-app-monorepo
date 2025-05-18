package main

import (
	"fmt"
	"net/http"
)

// Handler for "POST /v1/reviews" endpoint
func (app *application) createMovieReviewHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Create Movie Review")
}

// Handler for "GET /v1/reviews/:id" endpoint
func (app *application) showMovieReviewHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Show the details for the Movie Review: %d", id)
}
