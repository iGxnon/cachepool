package cachepool

import (
	"context"
	"database/sql"
	"github.com/igxnon/cachepool/cache"
	"github.com/streadway/amqp"
	"time"
	_ "unsafe"
)

type CachePool struct {
	Cache    *cache.Cache
	Db       *sql.DB
	cancelMQ context.CancelFunc
}

// UseMQ 使用 rabbitmq 同步一些缓存
func (c *CachePool) UseMQ(ch *amqp.Channel, name string) {
	ctx, cancel := context.WithCancel(context.Background())
	c.cancelMQ = cancel
	go run_mq(ctx, c.Cache, ch, name)
}

// StopMQ 停用 rabbitmq
func (c *CachePool) StopMQ() {
	if c.cancelMQ != nil {
		c.cancelMQ()
	}
}

//go:linkname run_mq douyin-common/cachepool/cache.runSyncFromMQ
//noinspection ALL
func run_mq(context.Context, *cache.Cache, *amqp.Channel, string)

func NewDefault(db *sql.DB) *CachePool {
	return &CachePool{
		Cache: cache.New(time.Minute*5, time.Hour/2, time.Minute*10),
		Db:    db,
	}
}

func New(db *sql.DB, cache *cache.Cache) *CachePool {
	return &CachePool{
		Cache: cache,
		Db:    db,
	}
}
