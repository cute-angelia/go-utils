package caches

import "time"

type (
	// Cache is the top-level cache interface
	Cache interface {

		// GenerateCacheKey 生产缓存 key ； bucket 用于内部业务隔离
		GenerateCacheKey(bucket string, key string) string

		// Get retrieve the cached key value
		Get(key string) (string, error)

		// GetMulti retrieve multiple cached keys value
		GetMulti(keys []string) map[string]string

		// Set cache a value by key
		Set(key string, value string, ttl time.Duration) error

		// SetWithBucket 在一个 bucket 保存数据，方便 scan 查找数据，并清楚一个 bucket 的数据
		SetWithBucket(bucket string, key string, value string, ttl time.Duration) error

		// Contains check if a cached key exists
		Contains(key string) bool

		// Delete remove the cached key
		Delete(key string) error

		// Flush remove all cached keys
		Flush() error

		// Scan 查询 bucket 里面所有数据
		Scan(bucket string, f func(key string) error) (err error)

		// Fold all key
		Fold(f func(key string) error) (err error)
	}
)
