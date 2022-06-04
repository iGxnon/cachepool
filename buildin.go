package cachepool

import (
	"database/sql"
	"github.com/igxnon/cachepool/cache"
	"time"
)

var buildInCachePool = &CachePool{
	// default expired in 5 min, deleted in 30min, synced in 10min
	Cache: cache.New(time.Minute*5, time.Hour/2, time.Minute*10),
}

func SetDatabase(db *sql.DB) {
	buildInCachePool.Db = db
}

func SetDefaultCache(c *cache.Cache) {
	buildInCachePool.Cache = c
}

// Get make sure what type you want to get
func Get[T any](key string) (T, bool) {
	got, ok := buildInCachePool.Cache.Get(key)
	return got.(T), ok
}

func GetWithExpiration[T any](key string) (T, time.Time, bool) {
	got, exp, ok := buildInCachePool.Cache.GetWithExpiration(key)
	return got.(T), exp, ok
}

func Set(key string, value any, exp time.Duration) {
	buildInCachePool.Cache.Set(key, value, exp)
}

func Add(key string, value any, exp time.Duration) error {
	return buildInCachePool.Cache.Add(key, value, exp)
}

func Replace(key string, new any, exp time.Duration) error {
	return buildInCachePool.Cache.Replace(key, new, exp)
}
