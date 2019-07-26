package config

import (
	"time"

	"webserver/cache"
)

// Create a cache
var c = cache.New(30*time.Minute, 60*time.Minute)

// GetCache ...
func GetCache() *cache.Cache {
	return c
}

// GetDefaultExpiration ...
func GetDefaultExpiration() time.Duration {
	return cache.DefaultExpiration
}
