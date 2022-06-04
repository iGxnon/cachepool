package cache

import (
	"context"
	gc "github.com/patrickmn/go-cache"
	"github.com/streadway/amqp"
	"go.uber.org/atomic"
	"golang.org/x/exp/maps"
	"runtime"
	"time"
)

const (
	// NoExpiration For use with functions that take an expiration time.
	NoExpiration time.Duration = -1
	// DefaultExpiration For use with functions that take an expiration time. Equivalent to
	// passing in the same expiration duration as was given to New() or
	// NewFrom() when the cache was created (e.g. 5 minutes.)
	DefaultExpiration time.Duration = 0
)

type Cache struct {
	*cache
	// play the same tricks as gc.Cache
}

type cache struct {
	*gc.Cache
	m map[string]gc.Item
	// protect syncExpired
	l            *atomic.Bool
	syncInterval time.Duration
	j            *janitor
}

func (c *cache) Get(k string) (interface{}, bool) {
	// try to get value in m first
	if value, ok := c.m[k]; ok && !c.l.Load() {
		if value.Expiration > 0 {
			if time.Now().UnixNano() > value.Expiration {
				return nil, false
			}
		}
		return value.Object, true
	}
	return c.Cache.Get(k)
}

func (c *cache) GetWithExpiration(k string) (interface{}, time.Time, bool) {
	if value, ok := c.m[k]; ok && !c.l.Load() {
		if value.Expiration > 0 {
			if time.Now().UnixNano() > value.Expiration {
				return nil, time.Time{}, false
			}

			// Return the item and the expiration time
			return value.Object, time.Unix(0, value.Expiration), true
		}

		// If expiration <= 0 (i.e. no expiration time set) then return the item
		// and a zeroed time.Time
		return value.Object, time.Time{}, true
	}
	return c.Cache.GetWithExpiration(k)
}

func (c *cache) Set(k string, x interface{}, d time.Duration) {
	if _, ok := c.m[k]; ok { // 命中 m 内的常驻缓存
		// 更新
		c.j.update <- struct {
			key  string
			item gc.Item
		}{key: k, item: gc.Item{
			Object:     x,
			Expiration: d.Nanoseconds(),
		}}
	}
	c.Cache.Set(k, x, d)
}

func (c *cache) SetDefault(k string, x interface{}) {
	c.Set(k, x, DefaultExpiration)
}

// TODO add more

// syncExpired sync gc.Cache into m
func (c *cache) syncExpired() {
	c.l.Store(true)
	defer c.l.Store(false)
	maps.Clear(c.m)
	for key, item := range c.Items() {
		// save NoExpiration and exist more than 30min into m
		if checkSticky(item) {
			c.m[key] = item
		}
	}
}

type janitor struct {
	Interval time.Duration
	stop     chan struct{}
	update   chan struct {
		key  string
		item gc.Item
	}
}

func (j *janitor) run(c *cache) {
	ticker := time.NewTicker(j.Interval)
	for {
		select {
		case <-ticker.C:
			c.syncExpired()
		case entry := <-j.update:
			c.l.Store(true)
			// save NoExpiration and exist more than 30min into m
			if checkSticky(entry.item) {
				c.m[entry.key] = entry.item
			} else {
				// 没有常驻缓存特征就删了
				delete(c.m, entry.key)
			}
			c.l.Store(false)
		case <-j.stop:
			ticker.Stop()
			return
		}
	}
}

func runJanitor(c *cache, ci time.Duration) {
	j := &janitor{
		Interval: ci,
		stop:     make(chan struct{}),
		update: make(chan struct {
			key  string
			item gc.Item
		}, 1),
	}
	c.j = j
	go j.run(c)
}

func stopJanitor(c *Cache) {
	c.j.stop <- struct{}{}
}

func checkSticky(item gc.Item) bool {
	switch item.Object.(type) {
	case int, int8, int16, int32, int64,
		uint, uintptr, uint8, uint16,
		uint32, uint64, float32, float64:
		return false // single number is not sticky(
	}
	now := time.Now().UnixNano()
	return item.Expiration <= 0 || now-item.Expiration > (time.Minute*30).Nanoseconds()
}

func New(defaultExpiration, cleanupInterval, syncInterval time.Duration) *Cache {
	if syncInterval > cleanupInterval && cleanupInterval > 0 {
		panic("syncInterval should not more than cleanupInterval")
	}
	ca := &cache{
		Cache:        gc.New(defaultExpiration, cleanupInterval),
		m:            make(map[string]gc.Item),
		l:            atomic.NewBool(false),
		syncInterval: syncInterval,
	}
	C := &Cache{ca}
	if syncInterval > 0 {
		runJanitor(ca, syncInterval)
		runtime.SetFinalizer(C, stopJanitor)
	}

	return C
}

func NewWithMQ(ctx context.Context, defaultExpiration, cleanupInterval, syncInterval time.Duration, ch *amqp.Channel, name string) *Cache {
	c := New(defaultExpiration, cleanupInterval, syncInterval)
	go runSyncFromMQ(ctx, c, ch, name)
	return c
}
