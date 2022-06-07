package redigo_cache

import (
	"github.com/gomodule/redigo/redis"
	"github.com/igxnon/cachepool/pkg/go-cache"
	"strconv"
	"sync"
	"testing"
	"time"
)

type Bar struct {
	foo string
}

func (b *Bar) Marshal() []byte {
	return []byte(b.foo)
}

func (b *Bar) Unmarshal(bytes []byte) bool {
	b.foo = string(bytes)
	return true
}

var conn, _ = redis.Dial("tcp", "127.0.0.1:6379")

func TestGlobalCache(t *testing.T) {
	tc := NewGlobalCache(time.Second*10, conn)
	tc.Set("foobarba", &Bar{foo: "hello"}, cache.DefaultExpiration)
	var got = &Bar{}
	ok := tc.GetUnmarshal("foobarba", got)
	if !ok || got.foo != "hello" {
		t.Error("error")
	}
}

func BenchmarkGlobalCacheGetExpiring(b *testing.B) {
	benchmarkGlobalCacheGet(b, 5*time.Minute)
}

func BenchmarkGlobalCacheGetNotExpiring(b *testing.B) {
	benchmarkGlobalCacheGet(b, cache.NoExpiration)
}

func benchmarkGlobalCacheGet(b *testing.B, exp time.Duration) {
	b.StopTimer()
	tc := NewGlobalCache(exp, conn)
	tc.Set("foobarba", &Bar{foo: "hello"}, cache.DefaultExpiration)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tc.Get("foobarba")
	}
}

func BenchmarkGlobalCacheGetManyConcurrentExpiring(b *testing.B) {
	benchmarkGlobalCacheGetManyConcurrent(b, 5*time.Minute)
}

func BenchmarkGlobalCacheGetManyConcurrentNotExpiring(b *testing.B) {
	benchmarkGlobalCacheGetManyConcurrent(b, cache.NoExpiration)
}

func benchmarkGlobalCacheGetManyConcurrent(b *testing.B, exp time.Duration) {
	b.StopTimer()
	n := 10000
	tsc := NewGlobalCache(exp, conn)
	keys := make([]string, n)
	for i := 0; i < n; i++ {
		k := "foo" + strconv.Itoa(i)
		keys[i] = k
		tsc.Set(k, &Bar{foo: "hello"}, cache.DefaultExpiration)
	}
	each := b.N / n
	wg := new(sync.WaitGroup)
	wg.Add(n)
	for _, v := range keys {
		go func(k string) {
			for j := 0; j < each; j++ {
				tsc.Get(k)
			}
			wg.Done()
		}(v)
	}
	b.StartTimer()
	wg.Wait()
}
