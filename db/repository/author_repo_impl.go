package repository

import (
	"article-service/apperror"
	"article-service/db/db_client"
	"article-service/db/transaction"
	"article-service/infrastructure/log"
	"article-service/model"
	"context"

	"github.com/google/uuid"
)

type AuthorRepo struct {
}

func GetAuthorRepository() IAuthorRepository {
	return AuthorRepo{}
}

func (r AuthorRepo) Get(ctx context.Context, id uuid.UUID) (*model.Author, error) {
	conn := transaction.GetClientOrTxn(ctx, db_client.GetDB)
	query := `
		SELECT
			authors.id,
			authors.name
		FROM authors
		WHERE id = $1
	`

	rows, err := conn.Query(ctx, query, id)
	if err != nil {
		log.Errorf(ctx, err, "[AuthorRepo][Get] Query failed")
		return nil, apperror.ErrGetRecordFailed
	}

	var author = model.Author{}
	for rows.Next() {
		err = rows.Scan(
			&author.ID,
			&author.Name,
		)
		if err != nil {
			log.Errorf(ctx, err, "[AuthorRepo][Get] Scan failed")
			return nil, apperror.ErrScanRecordFailed
		}
	}

	if author.ID == uuid.Nil {
		return nil, apperror.ErrObjectNotExists
	}

	return &author, nil
}
