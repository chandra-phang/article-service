package application

import (
	"context"

	"article-service/apperror"
	"article-service/db/db_client"
	"article-service/db/repository"
	v1req "article-service/dto/request/v1_req"
	"article-service/infrastructure/log"
	"article-service/model"
	"article-service/search"
	"article-service/utils"

	"github.com/google/uuid"
)

//go:generate mockgen -source=article_service.go -destination=./mock_application/article_service_mock.go
type IArticleService interface {
	CreateArticle(ctx context.Context, dto v1req.CreateArticleDTO) (uuid.UUID, error)
	ListArticles(ctx context.Context, dto v1req.ListArticlesDTO) ([]*model.Article, int64, error)
}

type ArticleSvc struct {
	articleRepo   repository.IArticleRepository
	authorRepo    repository.IAuthorRepository
	articleSearch search.IArticleSearch
}

var articleSvcSingleton IArticleService

func InitArticleService() {
	articleSvcSingleton = ArticleSvc{
		repository.GetArticleRepository(),
		repository.GetAuthorRepository(),
		search.GetArticleSearch(),
	}
}

func GetArticleService() IArticleService {
	return articleSvcSingleton
}

func (svc ArticleSvc) CreateArticle(ctx context.Context, dto v1req.CreateArticleDTO) (uuid.UUID, error) {
	ctx, txn, err := db_client.StartTransactionCtx(ctx)
	if err != nil {
		log.Errorf(ctx, err, "[ArticleSvc][CreateArticle] failed to start transaction")
		return uuid.Nil, apperror.ErrStartTransactionFailed
	}
	defer txn.Rollback(ctx)

	authorID, _ := uuid.Parse(dto.AuthorId)
	author, err := svc.authorRepo.Get(ctx, authorID)
	if err != nil {
		if err == apperror.ErrObjectNotExists {
			log.Errorf(ctx, err, "[ArticleSvc][CreateArticle] author not found, id: %s", dto.AuthorId)
			return uuid.Nil, apperror.ErrAuthorNotFound
		}
		log.Errorf(ctx, err, "[ArticleSvc][CreateArticle] authorRepo.Get is failed, id: %s", dto.AuthorId)
		return uuid.Nil, err
	}

	article := model.Article{
		ID:     utils.GenerateUUID(),
		Title:  dto.Title,
		Body:   dto.Body,
		Author: *author,
	}

	err = svc.articleRepo.Create(ctx, &article)
	if err != nil {
		log.Errorf(ctx, err, "[ArticleSvc][CreateArticle] articleRepo.Create is failed, article: %v", article)
		return uuid.Nil, err
	}

	err = svc.articleSearch.Index(ctx, article)
	if err != nil {
		log.Errorf(ctx, err, "[ArticleSvc][CreateArticle] articleSearch.Index is failed, article: %v", article)
		return uuid.Nil, err
	}

	if err = txn.Commit(ctx); err != nil {
		log.Errorf(ctx, err, "[ArticleSvc][CreateArticle] txn.Commit is failed!")
		return uuid.Nil, apperror.ErrCommitTransactionFailed
	}

	return article.ID, nil
}

func (svc ArticleSvc) ListArticles(ctx context.Context, dto v1req.ListArticlesDTO) ([]*model.Article, int64, error) {
	limit := utils.SetLimit(dto.Limit)

	filter := repository.ArticleFilter{
		AuthorName:    dto.AuthorName,
		SortBy:        dto.SortBy,
		SortDirection: dto.SortDirection,
		Limit:         limit,
		Offset:        utils.SetOffset(dto.Page, limit),
	}
	if dto.Query != "" {
		ids, err := svc.articleSearch.Search(ctx, dto.Query)
		if err != nil {
			log.Errorf(ctx, err, "[ArticleSvc][ListArticles] articleSearch.Search is failed, query: %s", dto.Query)
			return nil, 0, err
		}
		filter.Ids = ids
	}

	articles, err := svc.articleRepo.List(ctx, filter)
	if err != nil {
		log.Errorf(ctx, err, "[ArticleSvc][ListArticles] articleRepo.List is failed")
		return nil, 0, err
	}

	recordsCount, err := svc.articleRepo.GetRecordsCount(ctx, filter)
	if err != nil {
		log.Errorf(ctx, err, "[ArticleSvc][ListArticles] articleRepo.GetRecordsCount is failed")
		return nil, 0, err
	}

	return articles, recordsCount, nil
}
