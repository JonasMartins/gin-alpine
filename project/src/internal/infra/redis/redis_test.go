package redis

import (
	"context"
	"log"
	"testing"
	"time"

	"gin-alpine/src/internal/configs"

	"github.com/stretchr/testify/assert"
)

func TestRedis(t *testing.T) {
	description := "Testing Redis cache functionality"
	defer func() {
		log.Printf("Test: %s\n", description)
		log.Println("Deferred tearing down.")
	}()

	cfg, _ := configs.LoadConfig()
	var rClient *RedisClient = nil

	t.Run("should validate config object", func(t *testing.T) {
		assert.NotNil(t, cfg)
	})

	t.Run("should validate redis client", func(t *testing.T) {
		rClient = NewRedisClient(cfg.RedisURL)
		assert.NotNil(t, rClient)
	})

	t.Run("should validate redis cache", func(t *testing.T) {
		if rClient == nil {
			rClient = NewRedisClient(cfg.RedisURL)
		}
		assert.NotNil(t, rClient.Cache)

		ctx := context.TODO()
		group := "users"
		routes := []string{
			"page=1&limit=2",
			"page=2&limit=2",
			"page=3&limit=2",
			"page=1&limit=5",
			"page=2&limit=10",
			"page=2&limit=10",
			"page=2&limit=10",
			"page=2&limit=10",
			"page=2&limit=10",
			"page=2&limit=10",
		}
		for _, x := range routes {
			err := rClient.SetCachedRoute(ctx, group, x, time.Minute)
			if err != nil {
				t.Fail()
			}
		}

		result := make(map[string]string)
		err := rClient.Cache.Get(ctx, group, &result)
		if err != nil {
			t.Fail()
		}
		err = rClient.Cache.Delete(ctx, group)
		if err != nil {
			t.Fail()
		}
	})

}
