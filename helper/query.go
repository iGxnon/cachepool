package helper

import (
	"context"
	"database/sql"
	"github.com/igxnon/cachepool"
	"github.com/igxnon/cachepool/helper/internal"
)

type ExecResult struct {
	Result sql.Result
	Err    error
}

func Query[T any](c cachepool.ICachePool, query string, args ...any) (rows []T, err error) {
	return QueryWithContext[T](context.Background(), c, query, args...)
}

func QueryWithContext[T any](
	ctx context.Context,
	c cachepool.ICachePool,
	query string, args ...any,
) (rows []T, err error) {
	return internal.HandleRows[[]T](ctx, c, query, args...)
}

func QueryRow[T any](c cachepool.ICachePool, query string, args ...any) (rows T, err error) {
	return QueryRowWithContext[T](context.Background(), c, query, args...)
}

func QueryRowWithContext[T any](
	ctx context.Context,
	c cachepool.ICachePool,
	query string, args ...any,
) (row T, err error) {
	return internal.HandleRow[T](ctx, c, query, args...)
}
