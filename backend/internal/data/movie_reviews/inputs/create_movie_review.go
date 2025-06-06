package inputs

import (
	"cinepulse.nlt.net/internal/validator"
	"strings"
	"unicode/utf8"
)

type CreateMovieReviewInput struct {
	ImdbID           string `json:"imdb_id"`
	Rating           int8   `json:"rating"`
	StatementComment string `json:"statement_comment"`
}

func ValidateCreateMovieReviewInput(v *validator.Validator, input *CreateMovieReviewInput) {
	v.AddErrorIfNot(strings.TrimSpace(input.ImdbID) != "", "imdb_id", "must be provided")
	v.AddErrorIfNot(input.Rating >= 1, "rating", "must be greater than zero")
	v.AddErrorIfNot(input.Rating <= 5, "rating", "must be at most equal to 5")
	v.AddErrorIfNot(input.StatementComment != "", "statement_comment", "must be provided")
	v.AddErrorIfNot(utf8.RuneCountInString(input.StatementComment) <= 280, "statement_comment", "must not have more than 280 characters")
}
