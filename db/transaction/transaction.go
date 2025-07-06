package transaction

import (
	"context"
)

//go:generate mockgen -source=transaction.go -destination=./transaction_mock.go -package=transaction
type Manager interface {
	TransactionContextKey() ContextKey
	StartTransaction(context.Context) (ITransaction, error)
}

type ITransaction interface {
	Commit(context.Context) error
	Rollback(context.Context) error
}

type ContextKey string

func Start(ctx context.Context, mgr Manager) (context.Context, ITransaction, error) {
	txn, err := mgr.StartTransaction(ctx)
	if err != nil {
		return ctx, nil, err
	}
	ctxKey := mgr.TransactionContextKey()
	ctx = context.WithValue(ctx, ctxKey, txn)
	return ctx, txn, nil
}

func GetClientOrTxn[T Manager](ctx context.Context, fallbackGetter func() T) T {
	client := fallbackGetter()
	ctxKey := client.TransactionContextKey()
	if txn, ok := getTxnFromContext[T](ctx, ctxKey); ok {
		client = txn
	}
	return client
}

func getTxnFromContext[T Manager](ctx context.Context, key ContextKey) (T, bool) {
	txn, ok := ctx.Value(key).(T)
	return txn, ok
}
