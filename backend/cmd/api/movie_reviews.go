package main

import (
	"cinepulse.nlt.net/internal/constants"
	"cinepulse.nlt.net/internal/data"
	"fmt"
	"net/http"
	"time"
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

	// Create a placeholder instance of the MovieReview struct
	movieReview := data.MovieReview{
		ID:        id,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ImdbID:    "tt24069848",
		Rating:    4,
		Statement: data.MovieReviewStatement{
			Comment:   "Lorem ipsum",
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		},
		Reactions: data.MovieReviewReactionMap{
			data.Agree: []int64{1, 3, 4},
			data.Funny: []int64{5, 6},
		},
		Version: 1,
	}

	err = app.writeJSON(w, http.StatusOK, movieReview, nil)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, constants.ErrorMessages[http.StatusInternalServerError], http.StatusInternalServerError)
	}
}
