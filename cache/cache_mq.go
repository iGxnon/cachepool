// cache_mq 用于在不同实例中同步共享的缓存，比如验证码等缓存

package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)

const exchangeName = "exchange.__cache_sync__"

func runSyncFromMQ(ctx context.Context, cache *Cache, ch *amqp.Channel, name string) {
	err := ch.ExchangeDeclare(
		exchangeName,
		"fanout",
		true,
		false,
		false,
		true,
		nil,
	)
	if err != nil {
		log.Panicf("declare exchange failed err %v\n", err)
	}
	_, err = ch.QueueDeclare(
		name,
		true,
		false,
		false,
		true,
		nil,
	)
	if err != nil {
		log.Panicf("declare queue failed err %v\n", err)
	}

	err = ch.QueueBind(
		name,
		fmt.Sprintf("%s-key", name),
		exchangeName,
		true, nil,
	)

	if err != nil {
		log.Panicf("bind queue failed err %v\n", err)
	}

	msg, err := ch.Consume(
		name,
		fmt.Sprintf("%s-consumer", name),
		true,
		true,
		false,
		true,
		nil,
	)
	if err != nil {
		log.Panicf("Err when create message queue, %v\n", err)
	}
	for {
		select {
		case <-ctx.Done():
			return
		case m, ok := <-msg:
			if !ok {
				log.Printf("cache %s mq channel closed\n", name)
				return
			}
			key, value, exp, err := decode(m.Body)
			if err != nil {
				log.Printf("Warning comsumer %s message %s err %v\n",
					m.ConsumerTag, m.MessageId, err)
				continue
			}
			cache.Set(key, value, exp)
			log.Printf("Save %s into cache %s from MQ\n", string(m.Body), name)
		}
	}
}

type data struct {
	Key   string        `json:"key"`
	Value any           `json:"value"`
	Exp   time.Duration `json:"Exp"`
}

func decode(b []byte) (key string, value any, exp time.Duration, err error) {
	d := data{}
	err = json.Unmarshal(b, &d)
	return d.Key, d.Value, d.Exp, err
}

// Publish 将缓存同步到所有实例里
func Publish(ch *amqp.Channel, key string, value any, d time.Duration) error {
	data := data{
		Key:   key,
		Value: value,
		Exp:   d,
	}
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return ch.Publish(exchangeName, "", false, false, amqp.Publishing{
		Timestamp:    time.Now(),
		MessageId:    key,
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         b,
	})
}
