package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

// Handler for "POST /v1/reviews" endpoint
func (app *application) createMovieReviewHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Create Movie Review")
}

// Handler for "GET /v1/reviews/:id" endpoint
func (app *application) showMovieReviewHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Show the details for the Movie Review: %d", id)
}
