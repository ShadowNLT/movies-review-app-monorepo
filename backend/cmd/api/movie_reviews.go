package main

import (
	"cinepulse.nlt.net/internal/data/movie_reviews"
	"cinepulse.nlt.net/internal/data/movie_reviews/inputs"
	"cinepulse.nlt.net/internal/data/shared"
	"cinepulse.nlt.net/internal/validator"
	"errors"
	"fmt"
	"net/http"
)

// Handler for "POST /v1/reviews" endpoint
func (app *application) createMovieReviewHandler(w http.ResponseWriter, r *http.Request) {
	var input inputs.CreateMovieReviewInput

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

	v := validator.New()
	if inputs.ValidateCreateMovieReviewInput(v, &input); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	result, err := app.models.MovieReviews.Insert(&input)
	if err != nil {
		switch {
		case errors.Is(err, movie_reviews.ErrDuplicateImdbID):
			app.conflictResponse(w, r, errors.New("a review for the same imdb ID already exists"))
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// We include a Location header to let the client know which URL they can find
	// The newly-created resource at.
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/reviews/%d", result.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"review": result}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// Handler for "GET /v1/reviews" endpoint
func (app *application) listMovieReviewsHandler(w http.ResponseWriter, r *http.Request) {
	var input inputs.ListMovieReviewsQueryInput

	v := validator.New()
	qs := r.URL.Query()

	input.Page = app.readInt(qs, "page", 1, v)
	input.PageSize = app.readInt(qs, "page_size", 20, v)

	if inputs.ValidateListMovieReviewsQueryInput(v, &input); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	reviews, metadata, err := app.models.MovieReviews.GetAll(&input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie_reviews": reviews, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// Handler for "GET /v1/reviews/:id" endpoint
func (app *application) showMovieReviewHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	movieReview, err := app.models.MovieReviews.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, shared.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movieReview": movieReview}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// Handler for "DELETE /v1/reviews/:id" endpoint
func (app *application) deleteMovieReviewHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.models.MovieReviews.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, shared.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "movie review successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// Handler for "PATCH /v1/reviews/:id" endpoint
func (app *application) updateMovieReviewHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// Fetch the version of the movieReview with given ID
	movieReviewVersion, err := app.models.MovieReviews.GetVersionFor(id)
	if err != nil {
		switch {
		case errors.Is(err, shared.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	var input inputs.UpdateMovieReviewInput
	err = app.readJSON(w, r, &input, 2048)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Rating == nil && input.StatementComment == nil {
		app.badRequestResponse(w, r, errors.New("at least one of rating or statement_comment must be specified"))
		return
	}

	v := validator.New()
	if inputs.ValidateUpdateMovieReviewInput(v, &input); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	result, err := app.models.MovieReviews.Update(&input, id, movieReviewVersion)
	if err != nil {
		switch {
		case errors.Is(err, shared.ErrEditConflict):
			app.conflictResponse(w, r, errors.New("unable to update the record due to an edit conflict, please try again"))
			return
		default:
			app.serverErrorResponse(w, r, err)
		}
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"movieReview": result}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
