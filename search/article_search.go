package search

import (
	"article-service/model"
	"context"

	"github.com/google/uuid"
)

//go:generate mockgen -source=article_search.go -destination=./mock_search/article_search_mock.go
type IArticleSearch interface {
	Index(ctx context.Context, article model.Article) error
	Search(ctx context.Context, query string) ([]uuid.UUID, error)
}
