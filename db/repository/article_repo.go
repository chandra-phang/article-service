package repository

import (
	"context"

	"article-service/model"

	"github.com/google/uuid"
)

//go:generate mockgen -source=article_repo.go -destination=./mock_repository/article_repo_mock.go
type IArticleRepository interface {
	Create(ctx context.Context, article *model.Article) error
	List(ctx context.Context, filter ArticleFilter) ([]*model.Article, error)
	GetRecordsCount(ctx context.Context, filter ArticleFilter) (int64, error)
}

type ArticleFilter struct {
	Ids           []uuid.UUID
	AuthorName    string
	SortBy        string
	SortDirection string
	Limit         int
	Offset        int
}
