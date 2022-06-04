package cachepool

import (
	"github.com/igxnon/cachepool/cache"
	"github.com/streadway/amqp"
	"testing"
	"time"
)

// Test sync with MQ
func TestMQInOnePool(t *testing.T) {
	pool := NewDefault(nil)
	// use message queue, sync some cache
	conn, _ := amqp.Dial("amqp://guest:guest@localhost:5672/")
	ch, _ := conn.Channel()
	pool.UseMQ(ch, "cache")

	// publish an importance message into cache (
	_ = cache.Publish(ch, "下北沢一番臭の伝説", 114514, time.Minute*5)

	// sleep for a second
	time.Sleep(time.Second)

	got, exp, ok := pool.Cache.GetWithExpiration("下北沢一番臭の伝説")
	if ok {
		t.Log(got, exp)
	}

	time.Sleep(time.Second)

	// stop using message queue
	pool.StopMQ()
}

func TestMQInManyPool(t *testing.T) {
	var (
		p1      = NewDefault(nil)
		p2      = NewDefault(nil)
		p3      = NewDefault(nil)
		conn, _ = amqp.Dial("amqp://guest:guest@localhost:5672/")
		ch1, _  = conn.Channel()
		ch2, _  = conn.Channel()
		ch3, _  = conn.Channel()
	)

	p1.UseMQ(ch1, "cache1")
	p2.UseMQ(ch2, "cache2")
	p3.UseMQ(ch3, "cache3")

	// publish an importance message into cache (
	_ = cache.Publish(ch1, "下北沢一番臭の伝説", struct {
		Age    int
		Prefix string
		Movie  string
	}{
		24,
		"野獣せんべい",
		"真夏の夜の银夢",
	}, time.Minute*5)

	// sleep for a second
	time.Sleep(time.Second)

	got, exp, ok := p1.Cache.GetWithExpiration("下北沢一番臭の伝説")
	if ok {
		t.Log(got, exp)
	}

	got, exp, ok = p2.Cache.GetWithExpiration("下北沢一番臭の伝説")
	if ok {
		t.Log(got, exp)
	}

	got, exp, ok = p3.Cache.GetWithExpiration("下北沢一番臭の伝説")
	if ok {
		t.Log(got, exp)
	}

	time.Sleep(time.Second)

	// stop using message queue
	p1.StopMQ()
	p2.StopMQ()
	p3.StopMQ()
}
