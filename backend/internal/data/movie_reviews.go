package data

import (
	"cinepulse.nlt.net/internal/validator"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	ErrDuplicateImdbID = errors.New("duplicate imdb_id")
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
	ID        int64                   `json:"id"`
	CreatedAt time.Time               `json:"created_at"`
	UpdatedAt time.Time               `json:"updated_at"`
	Reactions *MovieReviewReactionMap `json:"reactions"`
	ImdbID    string                  `json:"imdb_id"`
	Rating    int8                    `json:"rating"`
	Statement MovieReviewStatement    `json:"statement"`
	Version   int64                   `json:"version"` // This will be incremented every time the user edits any of the editable information about the review
}

type CreateMovieReviewInput struct {
	ImdbID           string `json:"imdb_id"`
	Rating           int8   `json:"rating"`
	StatementComment string `json:"statement_comment"`
}

type CreatedMovieReview struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Version   int64     `json:"version"`
}

type UpdateMovieReviewInput struct {
	Rating           *int8   `json:"rating"`
	StatementComment *string `json:"statement_comment"`
}

func ValidateCreateMovieReviewInput(v *validator.Validator, input *CreateMovieReviewInput) {
	v.AddErrorIfNot(strings.TrimSpace(input.ImdbID) != "", "imdb_id", "must be provided")
	v.AddErrorIfNot(input.Rating >= 1, "rating", "must be greater than zero")
	v.AddErrorIfNot(input.Rating <= 5, "rating", "must be at most equal to 5")
	v.AddErrorIfNot(input.StatementComment != "", "statement_comment", "must be provided")
	v.AddErrorIfNot(utf8.RuneCountInString(input.StatementComment) <= 280, "statement_comment", "must not have more than 280 characters")
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

type MovieReviewModel struct {
	DB *sql.DB
}

func (m MovieReviewModel) Insert(review *CreateMovieReviewInput) (*CreatedMovieReview, error) {
	query := `
         INSERT INTO movie_reviews (
                                    imdb_id,
                                    rating,
                                    statement_comment
         )
         VALUES ($1, $2, $3)
         RETURNING id, created_at, version;
   `
	var result CreatedMovieReview
	args := []any{review.ImdbID, review.Rating, review.StatementComment}
	err := m.DB.QueryRow(query, args...).Scan(&result.ID, &result.CreatedAt, &result.Version)
	if err != nil {
		var pqErr *pq.Error
		switch {
		case errors.As(err, &pqErr) && pqErr.Code.Name() == "unique_violation":
			return &CreatedMovieReview{}, ErrDuplicateImdbID
		default:
			return &CreatedMovieReview{}, err
		}
	}
	return &result, nil
}

func (m MovieReviewModel) Get(id int64) (*MovieReview, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		 SELECT id, imdb_id, rating, statement_comment, statement_created_at, statement_updated_at,
                created_at, updated_at, version
         FROM movie_reviews
		 WHERE id = $1;`

	var movieReview MovieReview
	movieReview.Reactions = nil
	err := m.DB.QueryRow(query, id).Scan(
		&movieReview.ID,
		&movieReview.ImdbID,
		&movieReview.Rating,
		&movieReview.Statement.Comment,
		&movieReview.Statement.CreatedAt,
		&movieReview.Statement.UpdatedAt,
		&movieReview.CreatedAt,
		&movieReview.UpdatedAt,
		&movieReview.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &movieReview, nil
}

func (m MovieReviewModel) Update(input *UpdateMovieReviewInput, id int64) (*MovieReview, error) {
	var (
		args       []any
		setClauses []string
	)
	argCount := 1
	if input.Rating != nil {
		setClauses = append(setClauses, fmt.Sprintf("rating = $%d", argCount))
		args = append(args, *input.Rating)
		argCount++
	}
	if input.StatementComment != nil {
		setClauses = append(setClauses, fmt.Sprintf("statement_comment = $%d", argCount))
		args = append(args, *input.StatementComment)
		argCount++

		setClauses = append(setClauses, "statement_updated_at = now()")
	}

	setClauses = append(setClauses, "updated_at = now()", "version = version + 1")

	args = append(args, id)
	query := fmt.Sprintf(`
		UPDATE movie_reviews
        SET %s
        WHERE id = $%d
        RETURNING id, imdb_id, rating, statement_comment, 
        statement_created_at, statement_updated_at, created_at, 
        updated_at, version`, strings.Join(setClauses, ", "), argCount)

	var movieReview MovieReview
	movieReview.Reactions = nil
	err := m.DB.QueryRow(query, args...).Scan(
		&movieReview.ID,
		&movieReview.ImdbID,
		&movieReview.Rating,
		&movieReview.Statement.Comment,
		&movieReview.Statement.CreatedAt,
		&movieReview.Statement.UpdatedAt,
		&movieReview.CreatedAt,
		&movieReview.UpdatedAt,
		&movieReview.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return &MovieReview{}, ErrRecordNotFound
		default:
			return &MovieReview{}, err
		}
	}

	return &movieReview, nil
}

func (m MovieReviewModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM movie_reviews
		WHERE id = $1;`

	result, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil

}
