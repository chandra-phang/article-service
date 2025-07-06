package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"article-service/apperror"
	"article-service/db/db_client"
	"article-service/db/transaction"
	"article-service/infrastructure/log"
	"article-service/model"
	"article-service/utils"
)

type ArticleRepo struct {
}

func GetArticleRepository() IArticleRepository {
	return ArticleRepo{}
}

func (r ArticleRepo) Create(ctx context.Context, article *model.Article) error {
	conn := transaction.GetClientOrTxn(ctx, db_client.GetDB)

	query := `
		INSERT INTO articles
			(id, title, body, author_id, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	res, err := conn.Exec(
		ctx,
		query,
		&article.ID,
		&article.Title,
		&article.Body,
		&article.Author.ID,
		time.Now(),
	)
	if err != nil {
		log.Errorf(ctx, err, "[ArticleRepo][Create] Exec failed")
		return apperror.ErrCreateRecordFailed
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		log.Errorf(ctx, err, "[ArticleRepo][Create] No affected rows")
		return apperror.ErrNoAffectedRows
	}

	return nil
}

func (r ArticleRepo) List(ctx context.Context, filter ArticleFilter) ([]*model.Article, error) {
	conn := transaction.GetClientOrTxn(ctx, db_client.GetDB)
	query := `
		SELECT
			articles.id,
			articles.title,
			articles.body,
			articles.created_at,
			authors.id AS author_id,
			authors.name AS author_name
		FROM articles
		JOIN authors ON articles.author_id = authors.id
		{{whereFilters}}
		{{orderBy}}
		{{limitAndOffset}}
	`

	var params []interface{}
	var whereFilters []string

	if len(filter.Ids) > 0 {
		idsQuery := "articles.id IN ("
		ids := []string{}
		for _, id := range filter.Ids {
			params = append(params, id)
			ids = append(ids, fmt.Sprintf("$%d", len(params)))
		}
		idsQuery += strings.Join(ids, ", ") + ")"
		whereFilters = append(whereFilters, idsQuery)
	}

	if filter.AuthorName != "" {
		params = append(params, "%"+filter.AuthorName+"%")
		whereFilters = append(whereFilters, fmt.Sprintf("LOWER(authors.name) LIKE LOWER($%d)", len(params)))
	}

	if len(whereFilters) != 0 {
		filters := "WHERE " + strings.Join(whereFilters, " AND ")
		query = strings.ReplaceAll(query, "{{whereFilters}}", filters)
	} else {
		query = strings.ReplaceAll(query, "{{whereFilters}}", "")
	}

	sortBy := "articles.created_at"
	sortDirection := "DESC"
	if filter.SortBy != "" && filter.SortDirection != "" {
		sortBy = filter.SortBy
		sortDirection = filter.SortDirection
	}

	sort := sortBy + " " + sortDirection
	orderByQuery := fmt.Sprintf("ORDER BY %s", sort)
	query = strings.ReplaceAll(query, "{{orderBy}}", orderByQuery)

	limit := utils.SetLimit(filter.Limit)
	params = append(params, limit, filter.Offset)
	limitOffsetQuery := fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(params)-1, len(params))
	query = strings.ReplaceAll(query, "{{limitAndOffset}}", limitOffsetQuery)

	rows, err := conn.Query(ctx, query, params...)
	if err != nil {
		log.Errorf(ctx, err, "[ArticleRepo][List] Query failed")
		return nil, apperror.ErrGetRecordFailed
	}

	var articles = []*model.Article{}
	for rows.Next() {
		var article model.Article
		err = rows.Scan(
			&article.ID,
			&article.Title,
			&article.Body,
			&article.CreatedAt,
			&article.Author.ID,
			&article.Author.Name,
		)
		if err != nil {
			log.Errorf(ctx, err, "[ArticleRepo][List] Scan failed")
			return nil, apperror.ErrScanRecordFailed
		}

		articles = append(articles, &article)
	}

	return articles, nil
}

func (r ArticleRepo) GetRecordsCount(ctx context.Context, filter ArticleFilter) (int64, error) {
	conn := transaction.GetClientOrTxn(ctx, db_client.GetDB)
	query := `
		SELECT COUNT(*)
		FROM articles
		JOIN authors ON articles.author_id = authors.id
		{{whereFilters}}
	`

	var params []interface{}
	var whereFilters []string

	if len(filter.Ids) > 0 {
		idsQuery := "articles.id IN ("
		ids := []string{}
		for _, id := range filter.Ids {
			params = append(params, id)
			ids = append(ids, fmt.Sprintf("$%d", len(params)))
		}
		idsQuery += strings.Join(ids, ", ") + ")"
		whereFilters = append(whereFilters, idsQuery)
	}

	if filter.AuthorName != "" {
		params = append(params, "%"+filter.AuthorName+"%")
		whereFilters = append(whereFilters, fmt.Sprintf("LOWER(authors.name) LIKE LOWER($%d)", len(params)))
	}

	if len(whereFilters) != 0 {
		filters := "WHERE " + strings.Join(whereFilters, " AND ")
		query = strings.ReplaceAll(query, "{{whereFilters}}", filters)
	} else {
		query = strings.ReplaceAll(query, "{{whereFilters}}", "")
	}

	rows, err := conn.Query(ctx, query, params...)
	if err != nil {
		log.Errorf(ctx, err, "[ArticleRepo][GetRecordsCount] Query failed")
		return 0, apperror.ErrGetRecordFailed
	}

	var rowsCount int64
	for rows.Next() {
		err = rows.Scan(&rowsCount)
		if err != nil {
			log.Errorf(ctx, err, "[ArticleRepo][GetRecordsCount] Scan failed")
			return 0, apperror.ErrScanRecordFailed
		}
	}

	return rowsCount, nil
}
