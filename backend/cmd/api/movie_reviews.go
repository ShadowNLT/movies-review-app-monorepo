package main

import (
	"cinepulse.nlt.net/internal/data"
	"fmt"
	"net/http"
	"time"
)

// Handler for "POST /v1/reviews" endpoint
func (app *application) createMovieReviewHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ImdbID           string `json:"imdb_id"`
		Rating           int8   `json:"rating"`
		StatementComment string `json:"statement_comment"`
	}

	// To calculate the maxBytes that this payload can ever be
	// We look at each property.
	// In the following, the overhead mentioned are the extra bytes needed to encode
	// the JSON property name like "imdb_id" or "rating"
	// An imdb_id can take up to 9 bytes (
	// a rating is an integer value in [1:5] which only needs a byte
	// StatementComment has at most 280 Characters which are UTF-8 characters for which a rune(32 bits) are needed
	// which means we need 4 * 280 bytes ~= 1120 bytes
	// The overhead to encode for each field are 11 + 10 + 23 = 64 bytes
	// Which gives a total of 1194 bytes, and we will round it to the next power of two: 2048 bytes
	maxBytes := int64(2048)

	err := app.readJSON(w, r, &input, maxBytes)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	fmt.Fprintf(w, "%+v\n", input)
}

// Handler for "GET /v1/reviews/:id" endpoint
func (app *application) showMovieReviewHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
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

	err = app.writeJSON(w, http.StatusOK, envelope{"review": movieReview}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
