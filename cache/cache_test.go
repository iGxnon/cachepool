package cache

import (
	"fmt"
	gc "github.com/patrickmn/go-cache"
	"strconv"
	"sync"
	"testing"
	"time"
)

func BenchmarkCacheSetNoExpiration(b *testing.B) {
	b.StopTimer()
	tc := New(gc.NoExpiration, 0, time.Hour/2)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tc.Set("foo", "bar", gc.DefaultExpiration)
	}
}

func BenchmarkCacheSetRaw(b *testing.B) {
	b.StopTimer()
	tc := gc.New(gc.NoExpiration, 0)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tc.Set("foo", "bar", gc.DefaultExpiration)
	}
}

func BenchmarkCacheGetNoExpiration(b *testing.B) {

	tc := New(gc.NoExpiration, 0, time.Hour/2)

	tc.SetDefault("foo", "bar")

	tc.syncExpired()

	for i := 0; i < 100_000; i++ {
		tc.SetDefault(fmt.Sprintf("foo%d", i), "bar")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tc.Get("foo")
	}

}

func BenchmarkCacheGetRaw(b *testing.B) {

	tc := gc.New(gc.NoExpiration, 0)

	tc.SetDefault("foo", "bar")

	for i := 0; i < 100_000; i++ {
		tc.SetDefault(fmt.Sprintf("foo%d", i), "bar")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tc.Get("foo")
	}
}

func BenchmarkCacheGetManyConcurrentNoExpiration(b *testing.B) {
	b.StopTimer()
	n := 10000
	tc := New(gc.NoExpiration, 0, time.Hour/2)
	keys := make([]string, n)
	for i := 0; i < n; i++ {
		k := "foo" + strconv.Itoa(n)
		keys[i] = k
		tc.SetDefault(k, "bar")
	}
	// also try to comment the next line
	tc.syncExpired()

	each := b.N / n
	wg := new(sync.WaitGroup)
	wg.Add(n)
	for _, v := range keys {
		go func(a string) {
			for j := 0; j < each; j++ {
				tc.Get(a)
			}
			wg.Done()
		}(v)
	}
	b.StartTimer()
	wg.Wait()
}

func BenchmarkCacheGetManyConcurrentRaw(b *testing.B) {
	// This is the same as BenchmarkCacheGetConcurrent, but its result
	// can be compared against BenchmarkShardedCacheGetManyConcurrent
	// in sharded_test.go.
	b.StopTimer()
	n := 10000
	tc := gc.New(gc.NoExpiration, 0)
	keys := make([]string, n)
	for i := 0; i < n; i++ {
		k := "foo" + strconv.Itoa(n)
		keys[i] = k
		tc.SetDefault(k, "bar")
	}
	each := b.N / n
	wg := new(sync.WaitGroup)
	wg.Add(n)
	for _, v := range keys {
		go func(a string) {
			for j := 0; j < each; j++ {
				tc.Get(a)
			}
			wg.Done()
		}(v)
	}
	b.StartTimer()
	wg.Wait()
}
