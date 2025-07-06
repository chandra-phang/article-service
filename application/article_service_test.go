package application

import (
	"article-service/apperror"
	"article-service/db/db_client"
	"article-service/db/repository"
	"article-service/db/repository/mock_repository"
	v1req "article-service/dto/request/v1_req"
	"article-service/factory"
	"article-service/infrastructure/elasticsearch"
	"article-service/model"
	"article-service/search/mock_search"
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_GetArticleService(t *testing.T) {
	svc := GetArticleService()
	assert.Nil(t, svc)

	elasticsearch.InitElasticSearchMock()
	InitArticleService()

	svc = GetArticleService()
	assert.NotNil(t, svc)
}

func Test_CreateArticle_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	db_client.InitDatabaseMockExpectingTxn(ctrl, true, true)

	author := factory.SampleAuthorChandra
	dto := v1req.CreateArticleDTO{
		Title:    "New Title",
		Body:     "New Body",
		AuthorId: author.ID.String(),
	}

	articleRepo := mock_repository.NewMockIArticleRepository(ctrl)
	authorRepo := mock_repository.NewMockIAuthorRepository(ctrl)
	articleSearch := mock_search.NewMockIArticleSearch(ctrl)

	authorRepo.EXPECT().Get(gomock.Any(), author.ID).Return(&author, nil)
	articleRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
	articleSearch.EXPECT().Index(gomock.Any(), gomock.Any()).Return(nil)

	svc := ArticleSvc{
		articleRepo:   articleRepo,
		authorRepo:    authorRepo,
		articleSearch: articleSearch,
	}
	id, err := svc.CreateArticle(context.Background(), dto)
	assert.NotEqual(t, uuid.Nil, id)
	assert.Nil(t, err)
}

func Test_CreateArticle_ReturnErr_WhenStartTransactionFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	db_client.InitDatabaseMockExpectingTxn(ctrl, false, true)

	author := factory.SampleAuthorChandra
	dto := v1req.CreateArticleDTO{
		Title:    "New Title",
		Body:     "New Body",
		AuthorId: author.ID.String(),
	}

	articleRepo := mock_repository.NewMockIArticleRepository(ctrl)
	authorRepo := mock_repository.NewMockIAuthorRepository(ctrl)
	articleSearch := mock_search.NewMockIArticleSearch(ctrl)

	svc := ArticleSvc{
		articleRepo:   articleRepo,
		authorRepo:    authorRepo,
		articleSearch: articleSearch,
	}
	id, err := svc.CreateArticle(context.Background(), dto)
	assert.Equal(t, uuid.Nil, id)
	assert.Equal(t, apperror.ErrStartTransactionFailed, err)
}

func Test_CreateArticle_ReturnErr_WhenAuthorNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	db_client.InitDatabaseMockExpectingTxn(ctrl, true, nil)

	author := factory.SampleAuthorChandra
	dto := v1req.CreateArticleDTO{
		Title:    "New Title",
		Body:     "New Body",
		AuthorId: author.ID.String(),
	}

	articleRepo := mock_repository.NewMockIArticleRepository(ctrl)
	authorRepo := mock_repository.NewMockIAuthorRepository(ctrl)
	articleSearch := mock_search.NewMockIArticleSearch(ctrl)

	authorRepo.EXPECT().Get(gomock.Any(), author.ID).Return(nil, apperror.ErrObjectNotExists)

	svc := ArticleSvc{
		articleRepo:   articleRepo,
		authorRepo:    authorRepo,
		articleSearch: articleSearch,
	}
	id, err := svc.CreateArticle(context.Background(), dto)
	assert.Equal(t, uuid.Nil, id)
	assert.Equal(t, apperror.ErrAuthorNotFound, err)
}

func Test_CreateArticle_ReturnErr_WhenGetAuthorFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	db_client.InitDatabaseMockExpectingTxn(ctrl, true, nil)

	author := factory.SampleAuthorChandra
	dto := v1req.CreateArticleDTO{
		Title:    "New Title",
		Body:     "New Body",
		AuthorId: author.ID.String(),
	}

	articleRepo := mock_repository.NewMockIArticleRepository(ctrl)
	authorRepo := mock_repository.NewMockIAuthorRepository(ctrl)
	articleSearch := mock_search.NewMockIArticleSearch(ctrl)

	authorRepo.EXPECT().Get(gomock.Any(), author.ID).Return(nil, apperror.ErrGetRecordFailed)

	svc := ArticleSvc{
		articleRepo:   articleRepo,
		authorRepo:    authorRepo,
		articleSearch: articleSearch,
	}
	id, err := svc.CreateArticle(context.Background(), dto)
	assert.Equal(t, uuid.Nil, id)
	assert.Equal(t, apperror.ErrGetRecordFailed, err)
}

func Test_CreateArticle_ReturnErr_WhenCreateArticleFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	db_client.InitDatabaseMockExpectingTxn(ctrl, true, nil)

	author := factory.SampleAuthorChandra
	dto := v1req.CreateArticleDTO{
		Title:    "New Title",
		Body:     "New Body",
		AuthorId: author.ID.String(),
	}

	articleRepo := mock_repository.NewMockIArticleRepository(ctrl)
	authorRepo := mock_repository.NewMockIAuthorRepository(ctrl)
	articleSearch := mock_search.NewMockIArticleSearch(ctrl)

	authorRepo.EXPECT().Get(gomock.Any(), author.ID).Return(&author, nil)
	articleRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(apperror.ErrCreateRecordFailed)

	svc := ArticleSvc{
		articleRepo:   articleRepo,
		authorRepo:    authorRepo,
		articleSearch: articleSearch,
	}
	id, err := svc.CreateArticle(context.Background(), dto)
	assert.Equal(t, uuid.Nil, id)
	assert.Equal(t, apperror.ErrCreateRecordFailed, err)
}

func Test_CreateArticle_ReturnErr_WhenIndexArticleFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	db_client.InitDatabaseMockExpectingTxn(ctrl, true, nil)

	author := factory.SampleAuthorChandra
	dto := v1req.CreateArticleDTO{
		Title:    "New Title",
		Body:     "New Body",
		AuthorId: author.ID.String(),
	}

	articleRepo := mock_repository.NewMockIArticleRepository(ctrl)
	authorRepo := mock_repository.NewMockIAuthorRepository(ctrl)
	articleSearch := mock_search.NewMockIArticleSearch(ctrl)

	authorRepo.EXPECT().Get(gomock.Any(), author.ID).Return(&author, nil)
	articleRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
	articleSearch.EXPECT().Index(gomock.Any(), gomock.Any()).Return(apperror.ErrIndexElasticFailed)

	svc := ArticleSvc{
		articleRepo:   articleRepo,
		authorRepo:    authorRepo,
		articleSearch: articleSearch,
	}
	id, err := svc.CreateArticle(context.Background(), dto)
	assert.Equal(t, uuid.Nil, id)
	assert.Equal(t, apperror.ErrIndexElasticFailed, err)
}

func Test_CreateArticle_ReturnErr_WhenCommitTransactionFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	db_client.InitDatabaseMockExpectingTxn(ctrl, true, false)

	author := factory.SampleAuthorChandra
	dto := v1req.CreateArticleDTO{
		Title:    "New Title",
		Body:     "New Body",
		AuthorId: author.ID.String(),
	}

	articleRepo := mock_repository.NewMockIArticleRepository(ctrl)
	authorRepo := mock_repository.NewMockIAuthorRepository(ctrl)
	articleSearch := mock_search.NewMockIArticleSearch(ctrl)

	authorRepo.EXPECT().Get(gomock.Any(), author.ID).Return(&author, nil)
	articleRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
	articleSearch.EXPECT().Index(gomock.Any(), gomock.Any()).Return(nil)

	svc := ArticleSvc{
		articleRepo:   articleRepo,
		authorRepo:    authorRepo,
		articleSearch: articleSearch,
	}
	id, err := svc.CreateArticle(context.Background(), dto)
	assert.Equal(t, uuid.Nil, id)
	assert.Equal(t, apperror.ErrCommitTransactionFailed, err)
}

func Test_ListArticle_Success_WithEmptyQuery(t *testing.T) {
	ctrl := gomock.NewController(t)
	db_client.InitDatabaseMock()

	articleRepo := mock_repository.NewMockIArticleRepository(ctrl)
	authorRepo := mock_repository.NewMockIAuthorRepository(ctrl)
	articleSearch := mock_search.NewMockIArticleSearch(ctrl)

	dto := v1req.ListArticlesDTO{
		AuthorName:    "Chandra",
		SortBy:        "created_at",
		SortDirection: "desc",
		Limit:         10,
		Page:          1,
	}
	mockArticles := []*model.Article{&factory.SampleArticle1}
	mockRecordsCount := int64(1)
	expectedFilter := repository.ArticleFilter{
		AuthorName:    dto.AuthorName,
		SortBy:        dto.SortBy,
		SortDirection: dto.SortDirection,
		Limit:         dto.Limit,
		Offset:        0,
	}

	articleRepo.EXPECT().List(gomock.Any(), expectedFilter).Return(mockArticles, nil)
	articleRepo.EXPECT().GetRecordsCount(gomock.Any(), expectedFilter).Return(mockRecordsCount, nil)

	svc := ArticleSvc{
		articleRepo:   articleRepo,
		authorRepo:    authorRepo,
		articleSearch: articleSearch,
	}

	articles, recordsCount, err := svc.ListArticles(context.Background(), dto)
	assert.Equal(t, mockArticles, articles)
	assert.Equal(t, mockRecordsCount, recordsCount)
	assert.Nil(t, err)
}

func Test_ListArticle_Success_WithQuery(t *testing.T) {
	ctrl := gomock.NewController(t)
	db_client.InitDatabaseMock()

	articleRepo := mock_repository.NewMockIArticleRepository(ctrl)
	authorRepo := mock_repository.NewMockIAuthorRepository(ctrl)
	articleSearch := mock_search.NewMockIArticleSearch(ctrl)

	dto := v1req.ListArticlesDTO{
		Query:         "Article",
		AuthorName:    "Chandra",
		SortBy:        "created_at",
		SortDirection: "desc",
		Limit:         10,
		Page:          1,
	}
	mockArticles := []*model.Article{&factory.SampleArticle1}
	mockRecordsCount := int64(1)
	mockArticleIds := []uuid.UUID{factory.SampleArticle1.ID}

	expectedFilter := repository.ArticleFilter{
		Ids:           mockArticleIds,
		AuthorName:    dto.AuthorName,
		SortBy:        dto.SortBy,
		SortDirection: dto.SortDirection,
		Limit:         10,
		Offset:        0,
	}

	articleSearch.EXPECT().Search(gomock.Any(), dto.Query).Return(mockArticleIds, nil)
	articleRepo.EXPECT().List(gomock.Any(), expectedFilter).Return(mockArticles, nil)
	articleRepo.EXPECT().GetRecordsCount(gomock.Any(), expectedFilter).Return(mockRecordsCount, nil)

	svc := ArticleSvc{
		articleRepo:   articleRepo,
		authorRepo:    authorRepo,
		articleSearch: articleSearch,
	}

	articles, recordsCount, err := svc.ListArticles(context.Background(), dto)
	assert.Equal(t, mockArticles, articles)
	assert.Equal(t, mockRecordsCount, recordsCount)
	assert.Nil(t, err)
}

func Test_ListArticle_ReturnErr_WhenSearchFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	db_client.InitDatabaseMock()

	articleRepo := mock_repository.NewMockIArticleRepository(ctrl)
	authorRepo := mock_repository.NewMockIAuthorRepository(ctrl)
	articleSearch := mock_search.NewMockIArticleSearch(ctrl)

	dto := v1req.ListArticlesDTO{Query: "Article"}
	articleSearch.EXPECT().Search(gomock.Any(), dto.Query).Return(nil, apperror.ErrSearchElasticFailed)

	svc := ArticleSvc{
		articleRepo:   articleRepo,
		authorRepo:    authorRepo,
		articleSearch: articleSearch,
	}

	articles, recordsCount, err := svc.ListArticles(context.Background(), dto)
	assert.Nil(t, articles)
	assert.Equal(t, int64(0), recordsCount)
	assert.Equal(t, apperror.ErrSearchElasticFailed, err)
}

func Test_ListArticle_ReturnErr_WhenListArticlesFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	db_client.InitDatabaseMock()

	articleRepo := mock_repository.NewMockIArticleRepository(ctrl)
	authorRepo := mock_repository.NewMockIAuthorRepository(ctrl)
	articleSearch := mock_search.NewMockIArticleSearch(ctrl)

	dto := v1req.ListArticlesDTO{Query: "Article"}
	mockArticleIds := []uuid.UUID{factory.SampleArticle1.ID}
	expectedFilter := repository.ArticleFilter{
		Ids:    mockArticleIds,
		Limit:  20,
		Offset: 0,
	}
	articleSearch.EXPECT().Search(gomock.Any(), dto.Query).Return(mockArticleIds, nil)
	articleRepo.EXPECT().List(gomock.Any(), expectedFilter).Return(nil, apperror.ErrGetRecordFailed)

	svc := ArticleSvc{
		articleRepo:   articleRepo,
		authorRepo:    authorRepo,
		articleSearch: articleSearch,
	}

	articles, recordsCount, err := svc.ListArticles(context.Background(), dto)
	assert.Nil(t, articles)
	assert.Equal(t, int64(0), recordsCount)
	assert.Equal(t, apperror.ErrGetRecordFailed, err)
}

func Test_ListArticle_ReturnErr_WhenGetRecordsCountFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	db_client.InitDatabaseMock()

	articleRepo := mock_repository.NewMockIArticleRepository(ctrl)
	authorRepo := mock_repository.NewMockIAuthorRepository(ctrl)
	articleSearch := mock_search.NewMockIArticleSearch(ctrl)

	dto := v1req.ListArticlesDTO{Query: "Article"}
	mockArticleIds := []uuid.UUID{factory.SampleArticle1.ID}
	mockArticles := []*model.Article{&factory.SampleArticle1}
	expectedFilter := repository.ArticleFilter{
		Ids:    mockArticleIds,
		Limit:  20,
		Offset: 0,
	}

	articleSearch.EXPECT().Search(gomock.Any(), dto.Query).Return(mockArticleIds, nil)
	articleRepo.EXPECT().List(gomock.Any(), expectedFilter).Return(mockArticles, nil)
	articleRepo.EXPECT().GetRecordsCount(gomock.Any(), expectedFilter).Return(int64(0), apperror.ErrGetRecordFailed)

	svc := ArticleSvc{
		articleRepo:   articleRepo,
		authorRepo:    authorRepo,
		articleSearch: articleSearch,
	}

	articles, recordsCount, err := svc.ListArticles(context.Background(), dto)
	assert.Nil(t, articles)
	assert.Equal(t, int64(0), recordsCount)
	assert.Equal(t, apperror.ErrGetRecordFailed, err)
}
