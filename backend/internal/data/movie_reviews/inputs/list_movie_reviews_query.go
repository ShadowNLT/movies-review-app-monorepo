package inputs

import "cinepulse.nlt.net/internal/validator"

type ListMovieReviewsQueryInput struct {
	Page     int
	PageSize int
}

func (i ListMovieReviewsQueryInput) Limit() int {
	return i.PageSize
}

func (i ListMovieReviewsQueryInput) Offset() int {
	return (i.Page - 1) * i.PageSize
}

func ValidateListMovieReviewsQueryInput(v *validator.Validator, input *ListMovieReviewsQueryInput) {
	v.AddErrorIfNot(input.Page > 0, "page", "must be greater than zero")
	v.AddErrorIfNot(input.Page <= 10_000_000, "page", "must be a maximum of 10 million")
	v.AddErrorIfNot(input.PageSize > 0, "page_size", "must be greater than zero")
	v.AddErrorIfNot(input.PageSize <= 100, "page_size", "must be a maximum of 100")
}
