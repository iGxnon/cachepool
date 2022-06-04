package helper

import (
	"context"
	"database/sql"
	"github.com/igxnon/cachepool"
)

type ExecResult struct {
	Result sql.Result
	Err    error
}

// Query 尝试在缓存中搜索，如果没有就去数据库里查询(db 字段不为 nil 情况下)，查询的结果会插入缓存中
// 目前只支持 T 是 map 的情况
// TODO T 支持 struct 等类型
func Query[T any](c *cachepool.CachePool, query string, args ...any) (rows []T, err error) {
	return QueryWithContext[T](context.Background(), c, query, args...)
}

// QueryWithContext 目前只支持 T 是 map 的情况
// TODO T 支持 struct 等类型
func QueryWithContext[T any](
	ctx context.Context,
	c *cachepool.CachePool,
	query string, args ...any,
) (rows []T, err error) {
	return handleRows[[]T](ctx, c, query, args...)
}

// QueryRow 目前只支持 T 是 map 的情况
// TODO T 支持 struct 等类型
func QueryRow[T any](c *cachepool.CachePool, query string, args ...any) (rows T, err error) {
	return QueryRowWithContext[T](context.Background(), c, query, args...)
}

// QueryRowWithContext 目前只支持 T 是 map 的情况
// TODO T 支持 struct 等类型
func QueryRowWithContext[T any](
	ctx context.Context,
	c *cachepool.CachePool,
	query string, args ...any,
) (row T, err error) {
	return handleRow[T](ctx, c, query, args...)
}

// CacheExec 将 Exec 写入请求存入缓存，适当时间再进行写入
func CacheExec(c *cachepool.CachePool, query string, args ...any) (future chan ExecResult) {
	return CacheExecWithContext(context.Background(), c, query, args...)
}

func CacheExecWithContext(
	ctx context.Context,
	c *cachepool.CachePool,
	query string, args ...any,
) (future chan ExecResult) {
	panic("unimplemented")
	return
}
