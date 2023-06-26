package cache

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
)

type cacheImpl struct {
	cache *cache.Cache
}

type Config struct {
	URL string `envconfig:"REDIS_URL" required:"true"`
}

func NewCache(config Config) Cache {
	opt, err := redis.ParseURL(config.URL)
	if err != nil {
		log.Fatal("new cache error", err)
	}
	rdb := redis.NewClient(&redis.Options{
		Addr: opt.Addr,
		// Username:     opt.Username,
		Password:     opt.Password,
		DB:           opt.DB,
		DialTimeout:  opt.DialTimeout,
		ReadTimeout:  opt.ReadTimeout,
		WriteTimeout: opt.WriteTimeout,
		PoolSize:     opt.PoolSize,
	})

	cache := cache.New(&cache.Options{
		Redis:      rdb,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	c := cacheImpl{cache}
	return &c
}

func (c *cacheImpl) Get(ctx context.Context, key string, value any) error {
	return c.cache.Get(context.Background(), key, &value)
}

func (c *cacheImpl) Set(ctx context.Context, key string, value any) error {
	err := c.cache.Set(&cache.Item{
		Key:   key,
		Value: value,
		TTL:   0,
	})
	return err
}

func (c *cacheImpl) SetWithTTL(ctx context.Context, key string, value any, ttl time.Duration) error {
	err := c.cache.Set(&cache.Item{
		Key:   key,
		Value: value,
		TTL:   ttl,
	})
	return err
}

func (c *cacheImpl) Del(ctx context.Context, key string) error {
	err := c.cache.Delete(ctx, key)
	return err
}
