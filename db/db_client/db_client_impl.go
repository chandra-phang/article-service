package db_client

// impl.go contains the implementation of package psqlclient interface.

import (
	"article-service/apperror"
	"article-service/db/transaction"
	"article-service/infrastructure/log"
	"article-service/lib"
	"context"
	"database/sql"

	"github.com/golang/mock/gomock"
)

// dbClient implements IConnection
type dbClient struct {
	pool *sql.DB
}

func (c *dbClient) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := c.pool.QueryContext(ctx, query, args...)
	if err != nil {
		log.Errorf(ctx, err, "[DB_Client][Query] Select() errored (%v) ---- %s", err, normalizeWhitespace(query))
	}
	return rows, err
}

func (c *dbClient) Exec(ctx context.Context, query string, args ...interface{}) (IResult, error) {
	result, err := c.pool.ExecContext(ctx, query, args...)
	if err != nil {
		log.Errorf(ctx, err, "[DB_Client][Exec] Exec errored (%v) ---- %s", err, normalizeWhitespace(query))
	}
	return result, err
}

func (c *dbClient) TransactionContextKey() transaction.ContextKey {
	return "psql transaction"
}

func (c *dbClient) StartTransaction(ctx context.Context) (transaction.ITransaction, error) {
	caller := lib.WhoCalledMe()

	tx, err := c.pool.BeginTx(ctx, nil)
	txID := generateTransactionID()

	if err != nil {
		log.Errorf(ctx, err, "[DB_Client][StartTransaction] [%s] Error acquiring transaction by caller: %s", txID, caller)
	}

	_, err = tx.Exec("SET CONSTRAINTS ALL DEFERRED")
	if err != nil {
		log.Errorf(ctx, err, "[DB_Client][StartTransaction] [%s] Failed %s", txID, caller)
	}

	return &txnImpl{
		txID,
		tx,
		false,
		nil,
		false,
		false,
	}, err
}

// txnImpl implements ITransaction
type txnImpl struct {
	id       string
	tx       *sql.Tx
	commitOK bool

	// only for testing purpose
	ctrl                *gomock.Controller
	expectStartSuccess  bool
	expectCommitSuccess interface{}
}

func (w *txnImpl) Commit(ctx context.Context) error {
	caller := lib.WhoCalledMe()

	err := w.tx.Commit()
	if err == nil {
		w.commitOK = true
	} else {
		log.Errorf(ctx, err, "[DB_Trx][Commit] [%s] Error during transaction commit by caller: %s", w.id, caller)
	}

	return err
}

func (w *txnImpl) Rollback(ctx context.Context) error {
	caller := lib.WhoCalledMe()

	if w.commitOK {
		return nil
	}

	err := w.tx.Rollback()
	if err != nil {
		log.Errorf(ctx, err, "[DB_Trx][Rollback] [%s] Error during transaction rollback, caller: %s", w.id, caller)
	}
	return err
}

func (w *txnImpl) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := w.tx.QueryContext(ctx, query, args...)
	if err != nil {
		log.Errorf(ctx, err, "[DB_Trx][Query] [%s] Error during txn.Select, sql: %s", w.id, normalizeWhitespace(query))
	}
	return rows, err
}

func (w *txnImpl) Exec(ctx context.Context, query string, args ...interface{}) (IResult, error) {
	result, err := w.tx.ExecContext(ctx, query, args...)
	if err != nil {
		log.Errorf(ctx, err, "[DB_Trx][Exec] [%s] Error during txn.Exec, sql: %s", w.id, normalizeWhitespace(query))
	}
	return result, err
}

func (w *txnImpl) StartTransaction(ctx context.Context) (transaction.ITransaction, error) {
	// normal flow
	if w.ctrl == nil {
		return w, nil
	}

	// mock transaction only for testing purpose
	mockTx := transaction.NewMockITransaction(w.ctrl)
	if !w.expectStartSuccess {
		return mockTx, apperror.ErrStartTransactionFailed
	}

	if w.expectCommitSuccess != nil {
		if w.expectCommitSuccess == true {
			mockTx.EXPECT().Commit(gomock.Any()).Return(nil)
		} else {
			mockTx.EXPECT().Commit(gomock.Any()).Return(apperror.ErrCommitTransactionFailed)
		}
	}

	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil)

	return mockTx, nil
}

func (w *txnImpl) TransactionContextKey() transaction.ContextKey {
	return "psql transaction"
}
