package repository

import (
	"article-service/model"
	"context"

	"github.com/google/uuid"
)

//go:generate mockgen -source=author_repo.go -destination=./mock_repository/author_repo_mock.go
type IAuthorRepository interface {
	Get(ctx context.Context, id uuid.UUID) (*model.Author, error)
}
