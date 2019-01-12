package cache

import (
	"time"

	cache "github.com/patrickmn/go-cache"
)

// Create a cache
var c = cache.New(5*time.Minute, 10*time.Minute)

// GetCache ...
func GetCache() *cache.Cache {
	return c
}

// GetDefaultExpiration ...
func GetDefaultExpiration() time.Duration {
	return cache.DefaultExpiration
}
