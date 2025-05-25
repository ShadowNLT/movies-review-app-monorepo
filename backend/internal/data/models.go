package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	MovieReviews MovieReviewModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		MovieReviews: MovieReviewModel{DB: db},
	}
}
