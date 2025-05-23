package data

import (
	"cinepulse.nlt.net/internal/validator"
	"strings"
	"time"
	"unicode/utf8"
)

type MovieReviewReaction = string

const (
	Agree            MovieReviewReaction = "Agree"            // 👍
	Insightful       MovieReviewReaction = "Insightful"       //🤯
	Funny            MovieReviewReaction = "Funny"            //😂
	ThoughtProvoking MovieReviewReaction = "ThoughtProvoking" //🤔
	Disagree         MovieReviewReaction = "Disagree"         //👎
	WellSaid         MovieReviewReaction = "WellSaid"         //🙌
)

type MovieReviewReactionMap = map[MovieReviewReaction][]int64

type MovieReviewStatement struct {
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MovieReview struct {
	ID        int64                  `json:"id"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
	Reactions MovieReviewReactionMap `json:"reactions"`
	ImdbID    string                 `json:"imdb_id"`
	Rating    int8                   `json:"rating"`
	Statement MovieReviewStatement   `json:"statement"`
	Version   int64                  `json:"version"` // This will be incremented every time the user edits any of the editable information about the review
}

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
