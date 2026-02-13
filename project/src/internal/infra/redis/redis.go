// Package redis
package redis

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
	Cache  *cache.Cache
}

func NewRedisClient(addr string) *RedisClient {
	c := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
		PoolSize: 10,
	})

	myCache := cache.New(&cache.Options{
		Redis:      c,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	return &RedisClient{
		Client: c,
		Cache:  myCache,
	}
}

func (r *RedisClient) Delete(ctx context.Context, key string) bool {
	return r.Client.Del(ctx, key).Err() == nil
}

func (r *RedisClient) Add(ctx context.Context, key string, value string) bool {
	return r.Client.Set(ctx, key, value, 0).Err() != nil
}

func (r *RedisClient) Get(ctx context.Context, key string) string {
	return r.Client.Get(ctx, key).Val()
}

func (r *RedisClient) ClearGroup(ctx context.Context, group string) error {
	var groupRoutes map[string]string
	err := r.Cache.Get(ctx, group, &groupRoutes)
	if err != nil {
		return err
	}
	if len(groupRoutes) == 0 {
		return nil
	}
	for x := range groupRoutes {
		err = r.Cache.Delete(ctx, x)
		if err != nil {
			return err
		}
	}
	return r.Cache.Delete(ctx, group)
}

func (r *RedisClient) SetCachedRoute(ctx context.Context, group string, route string, ttl time.Duration) error {
	var groupRoutes map[string]string
	err := r.Cache.Get(ctx, group, &groupRoutes)
	if err != nil {
		return err
	}
	if len(groupRoutes) == 0 {
		groupRoutes = make(map[string]string)
		groupRoutes[route] = route
		err := r.Cache.Set(&cache.Item{
			Ctx:   ctx,
			Key:   group,
			Value: &groupRoutes,
			TTL:   ttl,
		})
		if err != nil {
			log.Printf("error setting cache %v", err)
			return err
		}
	} else {
		// se a rota já existe não é setada novamente
		_, exists := groupRoutes[route]
		if exists {
			return nil
		}
		groupRoutes[route] = route
		err := r.Cache.Set(&cache.Item{
			Ctx:   ctx,
			Key:   group,
			Value: &groupRoutes,
			TTL:   ttl,
		})
		if err != nil {
			log.Printf("error setting cache %v", err)
			return err
		}
	}
	return nil
}

func (r *RedisClient) DeleteCache(ctx context.Context, key string) error {
	return r.Cache.Delete(ctx, key)
}

func (r *RedisClient) GetAppDataKey() string {
	return "app:data"
}
