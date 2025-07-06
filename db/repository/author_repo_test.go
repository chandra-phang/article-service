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
	"github.com/stretchr/testify/assert"
)

func Test_Author_Get_Success(t *testing.T) {
	mock := db_client.InitDatabaseMock()

	author := factory.SampleAuthorChandra
	query := regexp.QuoteMeta(`
		SELECT
			authors.id,
			authors.name
		FROM authors
		WHERE id = $1
	`)

	columns := []string{"id", "name"}
	mock.ExpectQuery(query).WithArgs(author.ID).
		WillReturnRows(
			sqlmock.NewRows(columns).AddRow(
				author.ID,
				author.Name,
			),
		)

	repo := GetAuthorRepository()
	result, err := repo.Get(context.Background(), author.ID)

	assert.Equal(t, &author, result)
	assert.Nil(t, err)
}

func Test_Author_Get_ReturnErr_WhenQueryFailed(t *testing.T) {
	mock := db_client.InitDatabaseMock()

	author := factory.SampleAuthorChandra
	query := regexp.QuoteMeta(`
		SELECT
			authors.id,
			authors.name
		FROM authors
		WHERE id = $1
	`)

	stubErr := errors.New("db error")
	mock.ExpectQuery(query).WithArgs(author.ID).
		WillReturnError(stubErr)

	repo := GetAuthorRepository()
	result, err := repo.Get(context.Background(), author.ID)

	assert.Empty(t, result)
	assert.Equal(t, apperror.ErrGetRecordFailed, err)
}

func Test_Author_Get_ReturnErr_WhenScanFailed(t *testing.T) {
	mock := db_client.InitDatabaseMock()

	author := factory.SampleAuthorChandra
	query := regexp.QuoteMeta(`
		SELECT
			authors.id,
			authors.name
		FROM authors
		WHERE id = $1
	`)

	columns := []string{"id", "name"}
	mock.ExpectQuery(query).WithArgs(author.ID).
		WillReturnRows(
			sqlmock.NewRows(columns).AddRow(
				"invalid-uuid",
				author.Name,
			),
		)

	repo := GetAuthorRepository()
	result, err := repo.Get(context.Background(), author.ID)

	assert.Empty(t, result)
	assert.Equal(t, apperror.ErrScanRecordFailed, err)
}

func Test_Author_Get_ReturnErr_WhenRecordNotFound(t *testing.T) {
	mock := db_client.InitDatabaseMock()

	author := factory.SampleAuthorChandra
	query := regexp.QuoteMeta(`
		SELECT
			authors.id,
			authors.name
		FROM authors
		WHERE id = $1
	`)

	columns := []string{"id", "name"}
	mock.ExpectQuery(query).WithArgs(author.ID).
		WillReturnRows(sqlmock.NewRows(columns))

	repo := GetAuthorRepository()
	result, err := repo.Get(context.Background(), author.ID)

	assert.Empty(t, result)
	assert.Equal(t, apperror.ErrObjectNotExists, err)
}
