package data

import (
	"cinepulse.nlt.net/internal/validator"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	ErrDuplicateImdbID     = errors.New("duplicate imdb_id")
	RequestTimeOutDuration = 3 * time.Second
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

type ListMovieReviewsQueryInput struct {
	Page     int
	PageSize int
}

func (i ListMovieReviewsQueryInput) limit() int {
	return i.PageSize
}

func (i ListMovieReviewsQueryInput) offset() int {
	return (i.Page - 1) * i.PageSize
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

func ValidateListMovieReviewsQueryInput(v *validator.Validator, input *ListMovieReviewsQueryInput) {
	v.AddErrorIfNot(input.Page > 0, "page", "must be greater than zero")
	v.AddErrorIfNot(input.Page <= 10_000_000, "page", "must be a maximum of 10 million")
	v.AddErrorIfNot(input.PageSize > 0, "page_size", "must be greater than zero")
	v.AddErrorIfNot(input.PageSize <= 100, "page_size", "must be a maximum of 100")
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
	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeOutDuration)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&result.ID, &result.CreatedAt, &result.Version)
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

func (m MovieReviewModel) GetVersionFor(id int64) (int64, error) {
	if id < 1 {
		return 0, ErrRecordNotFound
	}

	query := `
         SELECT version
         FROM movie_reviews
         WHERE id = $1;`

	var version int64
	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeOutDuration)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(&version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return 0, ErrRecordNotFound
		default:
			return 0, err
		}
	}
	return version, nil
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

	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeOutDuration)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
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

func (m MovieReviewModel) Update(input *UpdateMovieReviewInput, id, version int64) (*MovieReview, error) {
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
	args = append(args, version)
	query := fmt.Sprintf(`
		UPDATE movie_reviews
        SET %s
        WHERE id = $%d AND version = $%d
        RETURNING id, imdb_id, rating, statement_comment, 
        statement_created_at, statement_updated_at, created_at, 
        updated_at, version`, strings.Join(setClauses, ", "), argCount, argCount+1)

	var movieReview MovieReview
	movieReview.Reactions = nil
	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeOutDuration)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
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
			return &MovieReview{}, ErrEditConflict
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

	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeOutDuration)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
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

func (m MovieReviewModel) GetAll(queryInput *ListMovieReviewsQueryInput) (reviews []*MovieReview, metadata Metadata, err error) {
	query := `
       	SELECT count(*) OVER(), id, imdb_id, rating, statement_comment, statement_created_at, statement_updated_at, created_at, updated_at, version
        FROM movie_reviews
       	ORDER BY updated_at DESC
       	LIMIT $1 OFFSET $2`

	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeOutDuration)
	defer cancel()

	// Get the total Count of Records
	totalRecords := 0
	err = m.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM movie_reviews").Scan(&totalRecords)

	args := []any{queryInput.limit(), queryInput.offset()}

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer func(rows *sql.Rows) {
		if cErr := rows.Close(); cErr != nil {
			if err == nil {
				err = cErr
			} else {
				err = errors.Join(
					err,
					cErr,
				)
			}
		}
	}(rows)

	reviews = []*MovieReview{}
	totalPaginatedRecords := 0

	for rows.Next() {
		var review MovieReview
		review.Reactions = nil

		err := rows.Scan(
			&totalPaginatedRecords,
			&review.ID,
			&review.ImdbID,
			&review.Rating,
			&review.Statement.Comment,
			&review.Statement.CreatedAt,
			&review.Statement.UpdatedAt,
			&review.CreatedAt,
			&review.UpdatedAt,
			&review.Version)
		if err != nil {
			return nil, Metadata{}, err
		}
		reviews = append(reviews, &review)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}
	metadata = calculateMetadata(totalPaginatedRecords, totalRecords, queryInput.Page, queryInput.PageSize)
	return reviews, metadata, nil
}
