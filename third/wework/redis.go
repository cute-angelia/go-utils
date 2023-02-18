package wework

import (
	"context"
	"github.com/faabiosr/cachego"
	rd "github.com/go-redis/redis/v8"
	"time"
)

type (
	redisV struct {
		driver *rd.Client
	}
)

func NewRedis(driver *rd.Client) cachego.Cache {
	return &redisV{driver}
}

// Contains checks if cached key exists in Redis storage
func (r *redisV) Contains(key string) bool {
	ctx := context.Background()
	return r.driver.Exists(ctx, key).Val() == 1
}

// Delete the cached key from Redis storage
func (r *redisV) Delete(key string) error {
	ctx := context.Background()
	return r.driver.Del(ctx, key).Err()
}

// Fetch retrieves the cached value from key of the Redis storage
func (r *redisV) Fetch(key string) (string, error) {
	ctx := context.Background()
	return r.driver.Get(ctx, key).Result()
}

// FetchMulti retrieves multiple cached value from keys of the Redis storage
func (r *redisV) FetchMulti(keys []string) map[string]string {
	ctx := context.Background()
	result := make(map[string]string)

	items, err := r.driver.MGet(ctx, keys...).Result()
	if err != nil {
		return result
	}

	for i := 0; i < len(keys); i++ {
		if items[i] != nil {
			result[keys[i]] = items[i].(string)
		}
	}

	return result
}

// Flush removes all cached keys of the Redis storage
func (r *redisV) Flush() error {
	ctx := context.Background()
	return r.driver.FlushAll(ctx).Err()
}

// Save a value in Redis storage by key
func (r *redisV) Save(key string, value string, lifeTime time.Duration) error {
	ctx := context.Background()
	return r.driver.Set(ctx, key, value, lifeTime).Err()
}
