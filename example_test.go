package cachepool

import (
	"github.com/igxnon/cachepool/cache"
	"github.com/streadway/amqp"
	"time"
)

func ExampleCachePool() {
	pool := NewDefault(nil)
	// Set
	pool.Cache.Set("foo", "bar", time.Minute*40)
	err := pool.Cache.Add("foo2", "bar2", cache.DefaultExpiration)
	if err != nil {
		// foo2 contains before
	}
	err = pool.Cache.Replace("foo", "barbar", cache.NoExpiration)
	if err != nil {
		// foo does not contain before
	}
	// Get
	_, _ = pool.Cache.Get("foo")
	_, _, _ = pool.Cache.GetWithExpiration("foo2")

	// increment and decrement
	pool.Cache.Set("foo3", 114514, cache.NoExpiration)

	_ = pool.Cache.Increment("foo3", 1919810) // then foo3 equals 2034324
	_ = pool.Cache.Decrement("foo3", 1919810)

	// use message queue, sync some cache
	conn, _ := amqp.Dial("amqp://guest:guest@localhost:5672/")
	ch, _ := conn.Channel()
	pool.UseMQ(ch, "cache1")

	// publish an importance message into cache (
	_ = cache.Publish(ch, "下北沢一番臭の伝説", struct {
		Age    int
		Prefix string
		Movie  string
	}{
		24,
		"野獣せんべい",
		"真夏の夜の银夢",
	}, time.Minute*5)

	time.Sleep(time.Second)

	// stop using message queue
	pool.StopMQ()
}
