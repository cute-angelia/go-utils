package cache

import "time"

type (
	// Cache is the top-level cache interface
	Cache interface {

		// Get retrieve the cached key value
		Get(key string) (string, error)

		// GetMulti retrieve multiple cached keys value
		GetMulti(keys []string) map[string]string

		// Set cache a value by key
		Set(key string, value string, ttl time.Duration) error

		// Contains check if a cached key exists
		Contains(key string) bool

		// Delete remove the cached key
		Delete(key string) error

		// Flush remove all cached keys
		Flush() error

		// Scan scan prefix
		Scan(prefix string, f func(key string) error) (err error)

		// Fold all key
		Fold(f func(key string) error) (err error)
	}
)
