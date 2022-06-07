package test

import (
	"github.com/gomodule/redigo/redis"
	"github.com/igxnon/cachepool"
	"github.com/igxnon/cachepool/pkg/go-cache"
	"testing"
	"time"
)

type Bar struct {
	foo string
}

func (b *Bar) Marshal() []byte {
	return []byte(b.foo)
}

func (b *Bar) Unmarshal(bs []byte) bool {
	b.foo = string(bs)
	return true
}

func TestDoubleCachePool(t *testing.T) {
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		t.Error("redis does not connect")
	}

	pool := cachepool.NewDouble(
		cachepool.WithBuidinGlobalCache(time.Minute*30, conn),
		cachepool.WithCache(cache.NewCache(time.Minute*5, time.Minute*10)))

	pool.SetDefault("foo", &Bar{foo: "yee"})
	bytes, ok := pool.Get("foo")
	if !ok {
		t.Error("not ok")
	}

	var bar = &Bar{}
	bar.Unmarshal(bytes.([]byte))
	if bar.foo != "yee" {
		t.Error("not yee")
	}
}

func BenchmarkDoubleCachePoolGet(b *testing.B) {
	benchmarkDoubleCachePoolGet(b, cache.NewCache(time.Minute*5, time.Minute*10))
}

func BenchmarkDoubleSyncMapCachePoolGet(b *testing.B) {
	benchmarkDoubleCachePoolGet(b, cache.NewSyncMapCache(time.Minute*5, time.Minute*10))
}

func benchmarkDoubleCachePoolGet(b *testing.B, c cache.ICache) {
	b.StopTimer()

	conn, _ := redis.Dial("tcp", "127.0.0.1:6379")

	pool := cachepool.NewDouble(
		cachepool.WithBuidinGlobalCache(time.Minute*30, conn),
		cachepool.WithCache(c))

	pool.SetDefault("foo", &Bar{foo: "yee"})
	bytes, ok := pool.Get("foo")
	if !ok {
		b.Error("not ok")
	}

	var bar = &Bar{}
	bar.Unmarshal(bytes.([]byte))
	if bar.foo != "yee" {
		b.Error("not yee")
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		pool.Get("foo")
	}
}
