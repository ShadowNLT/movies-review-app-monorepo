package data

import (
	"database/sql"
	"errors"
)

// ErrRecordNotFound is for when we are looking up something that doesn't exist in our database
var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	MovieReviews MovieReviewModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		MovieReviews: MovieReviewModel{DB: db},
	}
}
