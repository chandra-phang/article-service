package repository

import (
	"article-service/apperror"
	"article-service/db/db_client"
	"article-service/factory"
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_Article_Create_Success(t *testing.T) {
	mock := db_client.InitDatabaseMock()

	article := factory.SampleArticle1
	query := regexp.QuoteMeta(`
		INSERT INTO articles
			(id, title, body, author_id, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`)

	mock.ExpectExec(query).WithArgs(
		article.ID,
		article.Title,
		article.Body,
		article.Author.ID,
		sqlmock.AnyArg(),
	).WillReturnResult(
		sqlmock.NewResult(0, 1),
	)

	repo := GetArticleRepository()
	err := repo.Create(context.Background(), &article)

	assert.Nil(t, err)
}

func Test_Article_Create_ReturnErr_WhenExecFailed(t *testing.T) {
	mock := db_client.InitDatabaseMock()

	article := factory.SampleArticle1
	query := regexp.QuoteMeta(`
		INSERT INTO articles
			(id, title, body, author_id, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`)

	stubErr := errors.New("db error")
	mock.ExpectExec(query).WithArgs(
		article.ID,
		article.Title,
		article.Body,
		article.Author.ID,
		sqlmock.AnyArg(),
	).WillReturnError(
		stubErr,
	)

	repo := GetArticleRepository()
	err := repo.Create(context.Background(), &article)

	assert.Equal(t, apperror.ErrCreateRecordFailed, err)
}

func Test_Article_Create_ReturnErr_WhenNoRowsAffected(t *testing.T) {
	mock := db_client.InitDatabaseMock()

	article := factory.SampleArticle1
	query := regexp.QuoteMeta(`
		INSERT INTO articles
			(id, title, body, author_id, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`)

	mock.ExpectExec(query).WithArgs(
		article.ID,
		article.Title,
		article.Body,
		article.Author.ID,
		sqlmock.AnyArg(),
	).WillReturnResult(
		sqlmock.NewResult(0, 0),
	)

	repo := GetArticleRepository()
	err := repo.Create(context.Background(), &article)

	assert.Equal(t, apperror.ErrNoAffectedRows, err)
}

func Test_Article_List_Success_WithEmptyFilter(t *testing.T) {
	mock := db_client.InitDatabaseMock()

	article := factory.SampleArticle1
	query := regexp.QuoteMeta(`
		SELECT
			articles.id,
			articles.title,
			articles.body,
			articles.created_at,
			authors.id AS author_id,
			authors.name AS author_name
		FROM articles
		JOIN authors ON articles.author_id = authors.id
		ORDER BY articles.created_at DESC
		LIMIT $1 OFFSET $2
	`)

	columns := []string{
		"id",
		"title",
		"body",
		"created_at",
		"author_id",
		"author_name",
	}

	mock.ExpectQuery(query).WithArgs(20, 0).
		WillReturnRows(
			sqlmock.NewRows(columns).AddRow(
				article.ID,
				article.Title,
				article.Body,
				article.CreatedAt,
				article.Author.ID,
				article.Author.Name,
			),
		)

	repo := GetArticleRepository()
	filter := ArticleFilter{}
	articles, err := repo.List(context.Background(), filter)

	assert.Equal(t, &article, articles[0])
	assert.Nil(t, err)
}

func Test_Article_List_Success_WithFilledFilter(t *testing.T) {
	mock := db_client.InitDatabaseMock()

	article := factory.SampleArticle1
	query := regexp.QuoteMeta(`
		SELECT
			articles.id,
			articles.title,
			articles.body,
			articles.created_at,
			authors.id AS author_id,
			authors.name AS author_name
		FROM articles
		JOIN authors ON articles.author_id = authors.id
		WHERE articles.id IN ($1) AND LOWER(authors.name) LIKE LOWER($2)
		ORDER BY title asc
		LIMIT $3 OFFSET $4
	`)

	columns := []string{
		"id",
		"title",
		"body",
		"created_at",
		"author_id",
		"author_name",
	}

	mock.ExpectQuery(query).WithArgs(article.ID, "%"+article.Author.Name+"%", 10, 20).
		WillReturnRows(
			sqlmock.NewRows(columns).AddRow(
				article.ID,
				article.Title,
				article.Body,
				article.CreatedAt,
				article.Author.ID,
				article.Author.Name,
			),
		)

	repo := GetArticleRepository()
	filter := ArticleFilter{
		Ids:           []uuid.UUID{article.ID},
		AuthorName:    article.Author.Name,
		SortBy:        "title",
		SortDirection: "asc",
		Limit:         10,
		Offset:        20,
	}
	articles, err := repo.List(context.Background(), filter)

	assert.Equal(t, &article, articles[0])
	assert.Nil(t, err)
}

func Test_Article_List_ReturnErr_WhenQueryFailed(t *testing.T) {
	mock := db_client.InitDatabaseMock()

	query := regexp.QuoteMeta(`
		SELECT
			articles.id,
			articles.title,
			articles.body,
			articles.created_at,
			authors.id AS author_id,
			authors.name AS author_name
		FROM articles
		JOIN authors ON articles.author_id = authors.id
		ORDER BY articles.created_at DESC
		LIMIT $1 OFFSET $2
	`)

	stubErr := errors.New("db error")
	mock.ExpectQuery(query).
		WithArgs(20, 0).
		WillReturnError(
			stubErr,
		)

	repo := GetArticleRepository()
	filter := ArticleFilter{}
	articles, err := repo.List(context.Background(), filter)
	assert.Empty(t, articles)
	assert.Equal(t, apperror.ErrGetRecordFailed, err)
}

func Test_Article_List_ReturnErr_WhenScanFailed(t *testing.T) {
	mock := db_client.InitDatabaseMock()

	query := regexp.QuoteMeta(`
		SELECT
			articles.id,
			articles.title,
			articles.body,
			articles.created_at,
			authors.id AS author_id,
			authors.name AS author_name
		FROM articles
		JOIN authors ON articles.author_id = authors.id
		ORDER BY articles.created_at DESC
		LIMIT $1 OFFSET $2
	`)

	article := factory.SampleArticle1
	columns := []string{
		"id",
		"title",
		"body",
		"created_at",
		"author_id",
		"author_name",
	}
	mock.ExpectQuery(query).
		WithArgs(20, 0).
		WillReturnRows(sqlmock.NewRows(columns).AddRow(
			article.ID,
			article.Title,
			article.Body,
			"invalid-datetime",
			article.Author.ID,
			article.Author.Name,
		))

	repo := GetArticleRepository()
	filter := ArticleFilter{}
	articles, err := repo.List(context.Background(), filter)
	assert.Empty(t, articles)
	assert.Equal(t, apperror.ErrScanRecordFailed, err)
}

func Test_Article_List_Success_WhenResultIsEmpty(t *testing.T) {
	mock := db_client.InitDatabaseMock()

	query := regexp.QuoteMeta(`
		SELECT
			articles.id,
			articles.title,
			articles.body,
			articles.created_at,
			authors.id AS author_id,
			authors.name AS author_name
		FROM articles
		JOIN authors ON articles.author_id = authors.id
		ORDER BY articles.created_at DESC
		LIMIT $1 OFFSET $2
	`)

	columns := []string{
		"id",
		"title",
		"body",
		"created_at",
		"author_id",
		"author_name",
	}

	mock.ExpectQuery(query).
		WithArgs(20, 0).
		WillReturnRows(sqlmock.NewRows(columns))

	repo := GetArticleRepository()
	filter := ArticleFilter{}
	articles, err := repo.List(context.Background(), filter)
	assert.Empty(t, articles)
	assert.Nil(t, err)
}

func Test_Article_GetRecordsCount_Success_WithEmptyFilter(t *testing.T) {
	mock := db_client.InitDatabaseMock()

	query := regexp.QuoteMeta(`
		SELECT COUNT(*)
		FROM articles
		JOIN authors ON articles.author_id = authors.id
	`)

	columns := []string{"count"}
	mock.ExpectQuery(query).WillReturnRows(sqlmock.NewRows(columns).AddRow(1))

	repo := GetArticleRepository()
	filter := ArticleFilter{}
	recordsCount, err := repo.GetRecordsCount(context.Background(), filter)

	assert.Equal(t, int64(1), recordsCount)
	assert.Nil(t, err)
}

func Test_Article_GetRecordsCount_Success_WithFilledFilter(t *testing.T) {
	mock := db_client.InitDatabaseMock()

	article := factory.SampleArticle1
	query := regexp.QuoteMeta(`
		SELECT COUNT(*)
		FROM articles
		JOIN authors ON articles.author_id = authors.id
		WHERE articles.id IN ($1) AND LOWER(authors.name) LIKE LOWER($2)
	`)

	columns := []string{"count"}

	mock.ExpectQuery(query).WithArgs(article.ID, "%"+article.Author.Name+"%").
		WillReturnRows(sqlmock.NewRows(columns).AddRow(1))

	repo := GetArticleRepository()
	filter := ArticleFilter{
		Ids:        []uuid.UUID{article.ID},
		AuthorName: article.Author.Name,
	}
	recordsCount, err := repo.GetRecordsCount(context.Background(), filter)

	assert.Equal(t, int64(1), recordsCount)
	assert.Nil(t, err)
}

func Test_Article_GetRecordsCount_ReturnErr_WhenQueryFailed(t *testing.T) {
	mock := db_client.InitDatabaseMock()

	query := regexp.QuoteMeta(`
		SELECT COUNT(*)
		FROM articles
		JOIN authors ON articles.author_id = authors.id
	`)

	stubErr := errors.New("db error")
	mock.ExpectQuery(query).WillReturnError(stubErr)

	repo := GetArticleRepository()
	filter := ArticleFilter{}
	recordsCount, err := repo.GetRecordsCount(context.Background(), filter)
	assert.Equal(t, int64(0), recordsCount)
	assert.Equal(t, apperror.ErrGetRecordFailed, err)
}

func Test_Article_GetRecordsCount_ReturnErr_WhenScanFailed(t *testing.T) {
	mock := db_client.InitDatabaseMock()

	query := regexp.QuoteMeta(`
		SELECT COUNT(*)
		FROM articles
		JOIN authors ON articles.author_id = authors.id
	`)

	columns := []string{"count"}
	mock.ExpectQuery(query).WillReturnRows(sqlmock.NewRows(columns).AddRow("invalid-count"))

	repo := GetArticleRepository()
	filter := ArticleFilter{}
	recordsCount, err := repo.GetRecordsCount(context.Background(), filter)
	assert.Equal(t, int64(0), recordsCount)
	assert.Equal(t, apperror.ErrScanRecordFailed, err)
}

func Test_Article_GetRecordsCount_Success_WhenResultIsEmpty(t *testing.T) {
	mock := db_client.InitDatabaseMock()

	query := regexp.QuoteMeta(`
		SELECT COUNT(*)
		FROM articles
		JOIN authors ON articles.author_id = authors.id
	`)

	columns := []string{"count"}

	mock.ExpectQuery(query).WillReturnRows(sqlmock.NewRows(columns))

	repo := GetArticleRepository()
	filter := ArticleFilter{}
	recordsCount, err := repo.GetRecordsCount(context.Background(), filter)
	assert.Equal(t, int64(0), recordsCount)
	assert.Nil(t, err)
}
