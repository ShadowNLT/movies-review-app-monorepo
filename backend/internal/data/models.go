package data

import (
	"cinepulse.nlt.net/internal/data/movie_reviews"
	"cinepulse.nlt.net/internal/data/users"
	"database/sql"
)

type Models struct {
	MovieReviews movie_reviews.MovieReviewModel
	Users        users.UserModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		MovieReviews: movie_reviews.MovieReviewModel{DB: db},
		Users:        users.UserModel{DB: db},
	}
}
