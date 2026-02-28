package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisCache(url, password string) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: password,
		DB:       0, // use default DB
	})

	return &RedisCache{
		client: rdb,
		ctx:    context.Background(),
	}
}

func (c *RedisCache) Set(key string, value interface{}, duration time.Duration) {
	data, err := json.Marshal(value)
	if err != nil {
		return
	}
	c.client.Set(c.ctx, key, data, duration)
}

func (c *RedisCache) Get(key string, dest interface{}) bool {
	val, err := c.client.Get(c.ctx, key).Result()
	if err != nil {
		return false
	}

	err = json.Unmarshal([]byte(val), dest)
	return err == nil
}
