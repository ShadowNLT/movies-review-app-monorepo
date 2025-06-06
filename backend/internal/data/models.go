package data

import (
	"cinepulse.nlt.net/internal/data/movie_reviews"
	"database/sql"
)

type Models struct {
	MovieReviews movie_reviews.MovieReviewModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		MovieReviews: movie_reviews.MovieReviewModel{DB: db},
	}
}
