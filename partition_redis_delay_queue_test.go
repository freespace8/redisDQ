package redisDQ

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

func TestNewPartitionRedisDelayQueue(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	q := NewPartitionRedisDelayQueue(client)
	topic := "test"
	count := 10
	var wg sync.WaitGroup
	wg.Add(count)
	go q.Consume(topic, 5, 0, func(msg *Msg) (func(), error) {
		fmt.Printf("consume partiton0: %+v\n", msg)
		wg.Done()
		return nil, nil
	})
	go q.Consume(topic, 5, 1, func(msg *Msg) (func(), error) {
		fmt.Printf("consume partiton1: %+v\n", msg)
		wg.Done()
		return nil, nil
	})
	for i := 0; i < count; i++ {
		q.Push(context.Background(), &Msg{
			Topic:     topic,
			Key:       topic + strconv.Itoa(i),
			Partition: i % 2,
			Body:      []byte(topic + strconv.Itoa(i)),
			Delay:     time.Second * time.Duration(count),
		})
		time.Sleep(time.Second)
	}
	wg.Wait()
}
