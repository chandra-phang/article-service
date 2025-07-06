package v1req

import (
	"context"

	"article-service/apperror"
	"article-service/infrastructure/log"

	"github.com/go-playground/validator/v10"
)

type ListArticlesDTO struct {
	Query         string
	AuthorName    string
	SortBy        string `validate:"omitempty,oneof=created_at title author_name"`
	SortDirection string `validate:"omitempty,oneof=asc desc"`
	Limit         int
	Page          int
}

type CreateArticleDTO struct {
	Title    string `json:"title" validate:"required"`
	Body     string `json:"body" validate:"required"`
	AuthorId string `json:"authorId" validate:"required,uuid"`
}

func (dto ListArticlesDTO) Validate(ctx context.Context) error {
	validate := validator.New()
	if err := validate.Struct(dto); err != nil {
		err = apperror.TryTranslateValidationErrors(err)
		log.Errorf(ctx, err, "[V1][ListArticlesDTO] Validation failed. dto: %v", dto)
		return err
	}

	return nil
}

func (dto CreateArticleDTO) Validate(ctx context.Context) error {
	validate := validator.New()
	if err := validate.Struct(dto); err != nil {
		err = apperror.TryTranslateValidationErrors(err)
		log.Errorf(ctx, err, "[V1][CreateArticleDTO] Validation failed. Dto: %v", dto)
		return err
	}

	return nil
}
