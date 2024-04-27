package redis

import (
	"context"
	"log"
	"rcsp/internal/model"
	"time"

	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	redis *cache.Cache
}

func NewRedisClient() *Redis {
	rd := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &Redis{cache.New(&cache.Options{
		Redis:      rd,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})}
}

var ctx = context.TODO()

func (r *Redis) AddOrder(order model.Order) error {
	if err := r.redis.Set(&cache.Item{
		Ctx:   ctx,
		Key:   order.OrderUid,
		Value: order,
		TTL:   time.Hour,
	}); err != nil {
		log.Printf("can't add order in cache, err: %e", err)
	}
	return nil
}

func (r *Redis) GetOrder(key string) (model.Order, error) {
	var order model.Order
	if err := r.redis.Get(ctx, key, &order); err != nil {
		log.Printf("can't get order from cache, err: %e", err)
		return model.Order{}, err
	}
	return order, nil
}

func (r *Redis) DeleteOrder(key string) error {
	err := r.redis.Delete(ctx, key)
	if err != nil {
		return nil
	}
	return err
}
