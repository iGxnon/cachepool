package cachepool

import (
	"database/sql"
	"github.com/igxnon/cachepool/pkg/go-cache"
	"time"
)

// DoubleCachePool act just like L1(localCache(readOnly map)) L2(localCache)
// L3(globalCache) cache, and SQL database is just like Memory
// if globalCache implemented fits all type of the value it stored and Get() could
// return the value directly, helper.Query could be used on this pool
type DoubleCachePool struct {
	localCache  cache.ICache
	globalCache cache.ICache
	db          *sql.DB
}

// Set cache sidecar
func (c *DoubleCachePool) Set(k string, x interface{}, d time.Duration) {
	c.globalCache.Set(k, x, d)
	c.localCache.Delete(k)
}

func (c *DoubleCachePool) SetDefault(k string, x interface{}) {
	c.Set(k, x, cache.DefaultExpiration)
}

func (c *DoubleCachePool) Add(k string, x interface{}, d time.Duration) error {
	err := c.globalCache.Add(k, x, d)
	if err != nil {
		return err
	}
	c.localCache.Delete(k)
	return nil
}

func (c *DoubleCachePool) Replace(k string, x interface{}, d time.Duration) error {
	err := c.globalCache.Replace(k, x, d)
	if err != nil {
		return err
	}
	c.localCache.Delete(k)
	return nil
}

func (c *DoubleCachePool) Get(k string) (interface{}, bool) {
	got, ok := c.localCache.Get(k)
	if ok {
		return got, ok
	}
	got, ok = c.globalCache.Get(k)
	if ok {
		c.localCache.SetDefault(k, got)
	}
	return got, ok
}

func (c *DoubleCachePool) GetWithExpiration(k string) (interface{}, time.Time, bool) {
	got, exp, ok := c.localCache.GetWithExpiration(k)
	if ok {
		return got, exp, ok
	}
	got, exp, ok = c.globalCache.GetWithExpiration(k)
	if ok {
		c.localCache.SetDefault(k, got)
	}
	return got, exp, ok
}

func (c *DoubleCachePool) Increment(k string, n int64) error {
	err := c.globalCache.Increment(k, n)
	if err != nil {
		return err
	}
	c.localCache.Delete(k)
	return nil
}

func (c *DoubleCachePool) Decrement(k string, n int64) error {
	err := c.globalCache.Decrement(k, n)
	if err != nil {
		return err
	}
	c.localCache.Delete(k)
	return nil
}

func (c *DoubleCachePool) Delete(k string) {
	c.globalCache.Delete(k)
	c.localCache.Delete(k)
}

func (c *DoubleCachePool) DeleteExpired() {
	c.globalCache.DeleteExpired()
	c.localCache.DeleteExpired()
}

func (c *DoubleCachePool) Items() map[string]cache.IItem {
	return c.globalCache.Items()
}

func (c *DoubleCachePool) ItemCount() int {
	return c.globalCache.ItemCount()
}

func (c *DoubleCachePool) Flush() {
	c.globalCache.Flush()
	c.localCache.Flush()
}

func (c *DoubleCachePool) GetDatabase() *sql.DB {
	return c.db
}

func NewDouble(opt ...Option) *DoubleCachePool {
	opts := loadOptions(opt...)
	if opts._globalCache == nil {
		panic("global cache should be declared")
	}
	return &DoubleCachePool{
		localCache:  opts.cache,
		globalCache: opts._globalCache,
		db:          opts.db,
	}
}
