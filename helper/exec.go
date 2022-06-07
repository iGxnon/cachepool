package helper

import (
	"context"
	"github.com/igxnon/cachepool"
)

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
