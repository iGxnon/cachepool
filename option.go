package cachepool

import (
	"database/sql"
	"github.com/gomodule/redigo/redis"
	"github.com/igxnon/cachepool/pkg/go-cache"
	redigocache "github.com/igxnon/cachepool/pkg/redigo-cache"
	"time"
)

type Option func(*Options)

type Options struct {
	db           *sql.DB
	cache        cache.ICache
	_globalCache cache.ICache
}

func loadOptions(options ...Option) *Options {
	opts := new(Options)
	for _, option := range options {
		option(opts)
	}
	if opts.cache == nil {
		opts.cache = cache.NewCache(time.Minute*5, time.Minute*30)
	}
	return opts
}

func WithCache(cache cache.ICache) Option {
	return func(opt *Options) {
		opt.cache = cache
	}
}

func WithGlobalCache(cache cache.ICache) Option {
	return func(opt *Options) {
		opt._globalCache = cache
	}
}

// WithBuidinGlobalCache use buildin redis global cache
// NOTE: after using buildin redis cache, all pointer of value set into cache should
// implement redigocache.Object
//
// type Object interface {
//		Marshal() []byte
//		Unmarshal([]byte) bool
// }
func WithBuidinGlobalCache(defaultExpiration time.Duration, conn redis.Conn) Option {
	return func(opt *Options) {
		opt._globalCache = redigocache.NewGlobalCache(defaultExpiration, conn)
	}
}

func WithDatabase(db *sql.DB) Option {
	return func(opt *Options) {
		opt.db = db
	}
}
