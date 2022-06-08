package test

import (
	"database/sql"
	"github.com/igxnon/cachepool"
	"github.com/igxnon/cachepool/helper"
	"github.com/igxnon/cachepool/pkg/go-cache"
	redigocache "github.com/igxnon/cachepool/pkg/redigo-cache"
	"testing"
	"time"
)

func TestGlobalCacheSugar(t *testing.T) {
	cache := redigocache.NewGlobalCacheSugar(time.Minute*30, conn)
	cache.SetDefault("foo", "bar")
	var bar string
	cache.GetUnmarshal("foo", &bar)
	if bar != "bar" {
		t.Error("error")
	}
}

func TestGlobalCacheSugarHelper(t *testing.T) {
	dsn := "root:12345678@tcp(127.0.0.1:3306)/awesome?charset=utf8mb4&parseTime=True"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return
	}
	cache := redigocache.NewGlobalCacheSugar(time.Minute*30, conn)
	pool := cachepool.New(cachepool.WithCache(cache), cachepool.WithDatabase(db))
	got, err := helper.QueryRow[Bar](pool, "SELECT * FROM t LIMIT 1 OFFSET 1")
	if err != nil {
		t.Error(err)
		return
	}
	got, err = helper.QueryRow[Bar](pool, "SELECT * FROM t LIMIT 1 OFFSET 1")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(got)
}

func BenchmarkGlobalCacheSugarGetExpiring(b *testing.B) {
	benchmarkGlobalCacheSugarGet(b, time.Minute*10)
}

func BenchmarkGlobalCacheSugarGetNoExpiring(b *testing.B) {
	benchmarkGlobalCacheSugarGet(b, cache.NoExpiration)
}

func benchmarkGlobalCacheSugarGet(b *testing.B, exp time.Duration) {
	b.StopTimer()
	tc := redigocache.NewGlobalCacheSugar(exp, conn)
	tc.SetDefault("foobarba", Bar{Yee: "hello"})
	var bar Bar
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tc.GetUnmarshal("foobarba", &bar)
		if bar.Yee != "hello" {
			b.Error("not hello")
		}
	}
}
