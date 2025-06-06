package inputs

import (
	"cinepulse.nlt.net/internal/validator"
	"unicode/utf8"
)

type UpdateMovieReviewInput struct {
	Rating           *int8   `json:"rating"`
	StatementComment *string `json:"statement_comment"`
}

func ValidateUpdateMovieReviewInput(v *validator.Validator, input *UpdateMovieReviewInput) {
	if input.Rating != nil {
		v.AddErrorIfNot(*input.Rating >= 1, "rating", "must be greater than zero")
		v.AddErrorIfNot(*input.Rating <= 5, "rating", "must be at most equal to 5")
	}

	if input.StatementComment != nil {
		v.AddErrorIfNot(*input.StatementComment != "", "statement_comment", "must be provided")
		v.AddErrorIfNot(utf8.RuneCountInString(*input.StatementComment) <= 280, "statement_comment", "must not have more than 280 characters")
	}
}
